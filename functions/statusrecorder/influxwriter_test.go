package main

import (
	"context"
	"strconv"
	"testing"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/sirupsen/logrus"
)

func Test_InfluxWrite(t *testing.T) {
	token := "not-here!"
	bucket := "chargerstatus"
	org := "erik.lupander@gmail.com"
	client := influxdb2.NewClient("https://eu-central-1-1.aws.cloud2.influxdata.com", token)
	defer client.Close()

	// user blocking write client for writes to desired bucket
	writeAPI := client.WriteAPIBlocking(org, bucket)

	now := time.Now()

	p := influxdb2.NewPointWithMeasurement("stat").
		AddTag("site", "Ionity Spekeröd").
		AddTag("weekday", now.Weekday().String()).
		AddTag("hour_of_day", strconv.Itoa(now.Hour())).
		AddField("available", 4).
		SetTime(time.Now())
	if err := writeAPI.WritePoint(context.Background(), p); err != nil {
		logrus.WithError(err).Info("error writing point")
	}
}
