package main

import (
	"github.com/eriklupander/powertracker/functions/statusrecorder/model"
	"strconv"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

type RecordWriter interface {
	Write(r model.Record)
	Flush()
}

type InfluxWriter struct {
	writeApi api.WriteAPI
	bucket   string
	org      string
}

func (iw *InfluxWriter) Flush() {
	iw.writeApi.Flush()
}

func NewInfluxWriter(bucket, org string) *InfluxWriter {
	client := influxdb2.NewClient("https://eu-central-1-1.aws.cloud2.influxdata.com", influxTbToken)

	return &InfluxWriter{writeApi: client.WriteAPI(org, bucket)}
}

func (iw *InfluxWriter) Write(r model.Record) {
	now := time.Now()

	p := influxdb2.NewPointWithMeasurement("charger_availability").
		AddTag("site", r.SiteName).
		AddTag("weekday", now.Weekday().String()).
		AddTag("hour_of_day", strconv.Itoa(now.Hour())).
		AddField("available", r.Available).
		SetTime(time.Now())

	iw.writeApi.WritePoint(p)
}
