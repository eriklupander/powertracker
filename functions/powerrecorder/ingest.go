package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/timestreamwrite"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/http2"
	"net"
	"net/http"
	"strconv"
	"time"
)

func ingest(rec record) {
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
		logrus.WithError(err).Error("error configuring HTTP transport")
		return
	}

	sess, err := session.NewSession(&aws.Config{Region: aws.String("eu-west-1"), MaxRetries: aws.Int(3), HTTPClient: &http.Client{Transport: tr}})
	if err != nil {
		logrus.WithError(err).Error("error creating AWS session")
		return
	}
	writeSvc := timestreamwrite.New(sess)

	databaseName := "powertracker"
	tableName := "power_record"

	writeData(writeSvc, databaseName, tableName, rec)
}

func writeData(writeSvc *timestreamwrite.TimestreamWrite, databaseName string, tableName string, rec record) {
	writeRecordsInput := &timestreamwrite.WriteRecordsInput{
		DatabaseName: aws.String(databaseName),
		TableName:    aws.String(tableName),
		Records: []*timestreamwrite.Record{
			{
				Dimensions: []*timestreamwrite.Dimension{
					{
						Name:  aws.String("homeId"),
						Value: aws.String(rec.HomeId),
					},
				},
				MeasureName:      aws.String("energy_used"),
				MeasureValue:     aws.String(strconv.FormatFloat(rec.AccumulatedConsumption, 'f', 6, 64)),
				MeasureValueType: aws.String("DOUBLE"),
				Time:             aws.String(strconv.FormatInt(time.Now().Unix(), 10)),
				TimeUnit:         aws.String("SECONDS"),
			},
		},
	}

	if _, err := writeSvc.WriteRecords(writeRecordsInput); err != nil {
		logrus.WithError(err).Error("error writing power usage records")
	}
}

type record struct {
	HomeId                 string
	AccumulatedConsumption float64
}
