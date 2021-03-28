package graph

import (
	"github.com/eriklupander/powertracker/functions/exporter/aggregator"
	"github.com/eriklupander/powertracker/functions/exporter/model"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
	"time"
)

func TestExport(t *testing.T) {
	data, err := Export(buildEntries())
	assert.NoError(t, err)
	assert.NotNil(t, data)
	ioutil.WriteFile("erik.png", data, os.FileMode(0755))
}

func TestExportWithAggregation(t *testing.T) {
	entries, err := aggregator.Aggregate(buildEntries(), "1h")
	data, err := Export(entries)
	assert.NoError(t, err)
	assert.NotNil(t, data)
	assert.NoError(t, ioutil.WriteFile("erik-agg.png", data, os.FileMode(0755)))
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
