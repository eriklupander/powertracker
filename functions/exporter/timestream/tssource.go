package timestream

import (
	"fmt"
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
			Timeout: 10 * time.Second,
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

	query := "SELECT pr.homeId, pr.measure_value::double, pr.time FROM powertracker.power_record pr"

	// apply some semi-ugly date predicates if applicable
	if fromStr != "" || toStr != "" {
		from, fromErr := time.Parse("2006-01-02", fromStr)
		to, toErr := time.Parse("2006-01-02", toStr)
		if fromErr == nil && toErr == nil {
			query += " WHERE pr.time > '" + from.Format("2006-01-02") + "' AND pr.time < '" + to.Format("2006-01-02") + "'"
		} else if toErr == nil {
			query += " WHERE pr.time < '" + to.Format("2006-01-02") + "'"
		} else if fromErr == nil {
			query += " WHERE pr.time > '" + from.Format("2006-01-02") + "'"
		}
	}

	query += " ORDER BY pr.time"
	output, err := s.querySvc.Query(&timestreamquery.QueryInput{ClientToken: &idempotencyKey, QueryString: &query})
	if err != nil {
		return nil, err
	}

	entries := make([]model.Entry, 0)
	lastAcc := -1.0
	currentUsage := 0.0
	for _, row := range output.Rows {

		homeId := processScalarType(row.Data[0])
		measure, err := strconv.ParseFloat(processScalarType(row.Data[1]), 64)
		if err != nil {
			fmt.Printf("error parsing float: %s %v\n", processScalarType(row.Data[1]), err)
			return nil, err
		}

		if lastAcc == -1 {
			currentUsage = 0
		} else {
			// day switch
			if measure < lastAcc {
				lastAcc = 0
				currentUsage = 0
			} else {
				currentUsage = measure - lastAcc
			}
		}

		created, err := time.Parse("2006-01-02 15:04:05", processScalarType(row.Data[2]))
		if err != nil {
			return nil, err
		}

		entries = append(entries, model.Entry{
			HomeId:           homeId,
			CurrentUsage:     currentUsage,
			AccumulatedDaily: measure,
			Created:          created,
		})
		lastAcc = measure
	}

	return entries, nil
}
