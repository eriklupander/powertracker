package graph

import (
	"github.com/eriklupander/powertracker/functions/exporter/aggregator"
	"github.com/eriklupander/powertracker/functions/exporter/model"
	"github.com/eriklupander/powertracker/functions/exporter/testutil"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
	"time"
)

func TestExportWith5MinuteAggregation(t *testing.T) {
	entries := testutil.LoadEntriesFromCSV(t, "../testdata/usage-20210327.csv")
	data, err := ExportLinePlot(entries)
	assert.NoError(t, err)
	assert.NotNil(t, data)
	assert.NoError(t, ioutil.WriteFile("usage-aggregated-5m.png", data, os.FileMode(0755)))
}
func TestExportWithHourAggregation(t *testing.T) {
	entries, err := aggregator.Aggregate1h(testutil.LoadEntriesFromCSV(t, "../testdata/usage-20210327.csv"))
	data, err := ExportLinePlot(entries)
	assert.NoError(t, err)
	assert.NotNil(t, data)
	assert.NoError(t, ioutil.WriteFile("usage-aggregated-1h.png", data, os.FileMode(0755)))
}
func TestExportWeekWith5MinuteAggregation(t *testing.T) {
	entries := testutil.LoadEntriesFromCSV(t, "../testdata/usage-210326-210403-5m.csv")
	data, err := ExportLinePlot(entries)
	assert.NoError(t, err)
	assert.NotNil(t, data)
	assert.NoError(t, ioutil.WriteFile("usage-week-5m.png", data, os.FileMode(0755)))
}

func TestExportWithHourAggregationMultiDay(t *testing.T) {
	entries, err := aggregator.Aggregate1h(testutil.LoadEntriesFromCSV(t, "../testdata/usage.csv"))
	data, err := ExportLinePlot(entries)
	assert.NoError(t, err)
	assert.NotNil(t, data)
	assert.NoError(t, ioutil.WriteFile("usage-aggregated-1h-multiday.png", data, os.FileMode(0755)))
}

func TestExportHistogram5m(t *testing.T) {
	entries := testutil.LoadEntriesFromCSV(t, "../testdata/usage-210326-210403-5m.csv")
	data, err := ExportHist(entries)
	assert.NoError(t, err)
	assert.NotNil(t, data)
	assert.NoError(t, ioutil.WriteFile("histogram-5m-multiday.png", data, os.FileMode(0755)))
}
func TestExportHistogram1h(t *testing.T) {
	entries := testutil.LoadEntriesFromCSV(t, "../testdata/usage-210326-210403-1h.csv")

	data, err := ExportHist(entries)
	assert.NoError(t, err)
	assert.NotNil(t, data)
	assert.NoError(t, ioutil.WriteFile("histogram-1h-multiday.png", data, os.FileMode(0755)))
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
