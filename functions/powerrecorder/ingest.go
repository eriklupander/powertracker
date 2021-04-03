package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/timestreamwrite"
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
			// DualStack: true,
			Timeout: 10 * time.Second,
		}).DialContext,
		MaxIdleConns:          2,
		IdleConnTimeout:       5 * time.Second,
		TLSHandshakeTimeout:   5 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	// So client makes HTTP/2 requests
	if err := http2.ConfigureTransport(tr); err != nil {
		fmt.Println("Error configuring HTTP transport:")
		fmt.Println(err)
		return
	}

	//sess := session.Must(session.NewSessionWithOptions(session.Options{
	//	Config: aws.Config{
	//		HTTPClient: &http.Client{Transport: tr},
	//	},
	//	SharedConfigState: session.SharedConfigEnable,
	//}))

	sess, err := session.NewSession(&aws.Config{Region: aws.String("eu-west-1"), MaxRetries: aws.Int(3), HTTPClient: &http.Client{Transport: tr}})
	if err != nil {
		fmt.Println("Error creating session:")
		fmt.Println(err)
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
				MeasureValue:     aws.String(strconv.FormatFloat(float64(rec.AccumulatedConsumption), 'f', 6, 64)),
				MeasureValueType: aws.String("DOUBLE"),
				Time:             aws.String(strconv.FormatInt(time.Now().Unix(), 10)),
				TimeUnit:         aws.String("SECONDS"),
			},
		},
	}

	_, err := writeSvc.WriteRecords(writeRecordsInput)

	if err != nil {
		fmt.Println("Error:")
		fmt.Println(err)
	} else {
		fmt.Println("Write records is successful")
	}
}

func describeTable(writeSvc *timestreamwrite.TimestreamWrite, databaseName, tableName string) {
	// Describe table.
	describeTableInput := &timestreamwrite.DescribeTableInput{
		DatabaseName: aws.String(databaseName),
		TableName:    aws.String(tableName),
	}
	describeTableOutput, err := writeSvc.DescribeTable(describeTableInput)

	if err != nil {
		fmt.Println("Error:")
		fmt.Println(err)
	} else {
		fmt.Println("Describe table is successful, below is the output:")
		fmt.Println(describeTableOutput)
	}
}

func describeDb(writeSvc *timestreamwrite.TimestreamWrite, databaseName string) {
	// Describe database.
	describeDatabaseInput := &timestreamwrite.DescribeDatabaseInput{
		DatabaseName: aws.String(databaseName),
	}

	describeDatabaseOutput, err := writeSvc.DescribeDatabase(describeDatabaseInput)

	if err != nil {
		fmt.Println("Error:")
		fmt.Println(err)
	} else {
		fmt.Println("Describe database is successful, below is the output:")
		fmt.Println(describeDatabaseOutput)
	}
}

type record struct {
	HomeId                 string
	AccumulatedConsumption float64
}
