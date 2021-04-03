package graph

import (
	"bytes"
	c "encoding/csv"
	"github.com/eriklupander/powertracker/functions/exporter/aggregator"
	"github.com/eriklupander/powertracker/functions/exporter/model"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestExportWith5MinuteAggregation(t *testing.T) {
	entries, err := aggregator.Aggregate(loadEntriesFromCSV(t, "../testdata/usage-20210327.csv"), "5m")
	data, err := ExportLinePlot(entries)
	assert.NoError(t, err)
	assert.NotNil(t, data)
	assert.NoError(t, ioutil.WriteFile("usage-aggregated-5m.png", data, os.FileMode(0755)))
}
func TestExportWithHourAggregation(t *testing.T) {
	entries, err := aggregator.Aggregate(loadEntriesFromCSV(t, "../testdata/usage-20210327.csv"), "1h")
	data, err := ExportLinePlot(entries)
	assert.NoError(t, err)
	assert.NotNil(t, data)
	assert.NoError(t, ioutil.WriteFile("usage-aggregated-1h.png", data, os.FileMode(0755)))
}
func TestExportWeekWith5MinuteAggregation(t *testing.T) {
	entries := loadEntriesFromCSV(t, "../testdata/usage-210326-210403-5m.csv")
	data, err := ExportLinePlot(entries)
	assert.NoError(t, err)
	assert.NotNil(t, data)
	assert.NoError(t, ioutil.WriteFile("usage-week-5m.png", data, os.FileMode(0755)))
}

func TestExportWithHourAggregationMultiDay(t *testing.T) {
	entries, err := aggregator.Aggregate(loadEntriesFromCSV(t, "../testdata/usage.csv"), "1h")
	data, err := ExportLinePlot(entries)
	assert.NoError(t, err)
	assert.NotNil(t, data)
	assert.NoError(t, ioutil.WriteFile("usage-aggregated-1h-multiday.png", data, os.FileMode(0755)))
}

func TestExportHistogram5m(t *testing.T) {
	entries := loadEntriesFromCSV(t, "../testdata/usage-210326-210403-5m.csv")
	data, err := ExportHist(entries)
	assert.NoError(t, err)
	assert.NotNil(t, data)
	assert.NoError(t, ioutil.WriteFile("histogram-5m-multiday.png", data, os.FileMode(0755)))
}
func TestExportHistogram1h(t *testing.T) {
	entries := loadEntriesFromCSV(t, "../testdata/usage-210326-210403-1h.csv")

	data, err := ExportHist(entries)
	assert.NoError(t, err)
	assert.NotNil(t, data)
	assert.NoError(t, ioutil.WriteFile("histogram-1h-multiday.png", data, os.FileMode(0755)))
}

func loadEntriesFromCSV(t *testing.T, filename string) []model.Entry {
	data, err := ioutil.ReadFile(filename)
	assert.NoError(t, err)

	reader := c.NewReader(bytes.NewBuffer(data))
	reader.FieldsPerRecord = 4
	records, err := reader.ReadAll()
	assert.NoError(t, err)

	out := make([]model.Entry, 0)
	for i := 1; i < len(records); i++ {
		usage, err := strconv.ParseFloat(records[i][1], 64)
		if err != nil {
			t.FailNow()
		}
		created, err := time.Parse("2006-01-02T15:04:05Z", records[i][3])
		if err != nil {
			t.FailNow()
		}
		out = append(out, model.Entry{
			CurrentUsage: usage,
			Created:      created,
		})
	}
	return out
}

func buildEntries() []model.Entry {
	then := time.Date(2021, 3, 1, 0, 0, 0, 0, time.Local)
	entries := make([]model.Entry, 1000)
	for i := 0; i < 1000; i++ {
		entries[i] = model.Entry{
			CurrentUsage: rand.Float64() * 1000,
			Created:      then,
		}
		then = then.Add(time.Minute * 5)
	}
	return entries
}
