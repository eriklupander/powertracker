package timestream

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/timestreamquery"
	"github.com/eriklupander/powertracker/functions/exporter/model"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/http2"
	"net"
	"net/http"
	"strconv"
	"time"
)

func NewDataSource() *Source {
	tr := &http.Transport{
		ResponseHeaderTimeout: 10 * time.Second,
		Proxy:                 http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			KeepAlive: 10 * time.Second,
			Timeout:   10 * time.Second,
		}).DialContext,
		MaxIdleConns:          2,
		IdleConnTimeout:       5 * time.Second,
		TLSHandshakeTimeout:   5 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	if err := http2.ConfigureTransport(tr); err != nil {
		logrus.Fatalf("error configuring HTTP transport: %v", err)
	}

	sess, err := session.NewSession(&aws.Config{Region: aws.String("eu-west-1"), MaxRetries: aws.Int(3), HTTPClient: &http.Client{Transport: tr}})
	if err != nil {
		logrus.Fatalf("error creating timestream session: %v", err)
	}
	querySvc := timestreamquery.New(sess)
	return &Source{
		querySvc: querySvc,
	}
}

type Source struct {
	querySvc *timestreamquery.TimestreamQuery
}

func (s *Source) GetAll(fromStr, toStr string) ([]model.Entry, error) {
	idempotencyKey := uuid.New().String()

	query := s.buildQuery(fromStr, toStr)
	logrus.Infof("Using query: %s", query)
	st := time.Now()
	output, err := s.querySvc.Query(&timestreamquery.QueryInput{ClientToken: &idempotencyKey, QueryString: &query})
	if err != nil {
		return nil, err
	}
	logrus.Infof("querying timestream DB took %v, producing %d results", time.Since(st), len(output.Rows))

	entries := make([]model.Entry, 0)
	lastAccumulativeValue := -1.0
	currentUsage := 0.0

	st = time.Now()

	for _, row := range output.Rows {

		homeId := *row.Data[0].ScalarValue
		measure, err := strconv.ParseFloat(*row.Data[1].ScalarValue, 64)
		if err != nil {
			return nil, err
		}
		created, err := time.Parse("2006-01-02 15:04:05", *row.Data[2].ScalarValue)
		if err != nil {
			return nil, err
		}

		if lastAccumulativeValue == -1.0 {
			currentUsage = 0.0
		} else {
			// day switch, then the measurement drops to 0 again.
			if measure < lastAccumulativeValue {
				lastAccumulativeValue = 0.0
				currentUsage = 0.0
			} else {
				currentUsage = measure - lastAccumulativeValue
			}
		}

		entries = append(entries, model.Entry{
			HomeId:           homeId,
			CurrentUsage:     currentUsage,
			AccumulatedDaily: measure,
			Created:          created,
		})
		lastAccumulativeValue = measure
	}
	logrus.Infof("processing timestream query results took %v", time.Since(st))

	return entries, nil
}

func (s *Source) buildQuery(fromStr string, toStr string) string {
	query := "SELECT pr.homeId, pr.measure_value::double, pr.time FROM powertracker.power_record pr"

	// apply some semi-ugly date predicates if applicable
	if fromStr != "" || toStr != "" {
		from, fromErr := time.ParseInLocation("2006-01-02", fromStr, time.Local)
		to, toErr := time.ParseInLocation("2006-01-02", toStr, time.Local)
		if fromErr == nil && toErr == nil {
			logrus.Infof("Local TZ is: %v", from.Location())
			query += " WHERE pr.time > '" + from.Format("2006-01-02 15:04:05") + "' AND pr.time < '" + to.Format("2006-01-02 15:04:05") + "'"
		} else if toErr == nil {
			logrus.Infof("Local TZ is: %v", to.Location())

			query += " WHERE pr.time < '" + to.Format("2006-01-02 15:04:05") + "'"
		} else if fromErr == nil {
			logrus.Infof("Local TZ is: %v", from.Location())

			query += " WHERE pr.time > '" + from.Format("2006-01-02 15:04:05") + "'"
		}
	}

	query += " ORDER BY pr.time"
	return query
}
