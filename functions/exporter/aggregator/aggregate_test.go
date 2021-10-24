package aggregator

import (
	"github.com/ahmetb/go-linq"
	"github.com/eriklupander/powertracker/functions/exporter/model"
	"github.com/eriklupander/powertracker/functions/exporter/testutil"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAggregator1H(t *testing.T) {
	entries := testutil.LoadEntriesFromCSV(t, "../testdata/rawusage-2021-03-26_01.csv")
	agg, err := Aggregate1h(entries)
	assert.NoError(t, err)
	assert.Equal(t, 24, len(agg), "expect hourly 24 entries for a day")

	sum := linq.From(entries).Select(func(i interface{}) interface{} {
		return i.(model.Entry).CurrentUsage
	}).SumFloats()
	assert.Len(t, agg, 24)
	assert.InEpsilon(t, 19.895473, sum, 0.1, "expect sum of aggregated data to equal actual")
}

func TestAggregator1M(t *testing.T) {
	entries := testutil.LoadEntriesFromCSV(t, "../testdata/rawusage-2021-03-26_01.csv")
	agg, err := Aggregate1M(entries)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(agg), "expect hourly 24 entries for a day")

	sum := linq.From(entries).Select(func(i interface{}) interface{} {
		return i.(model.Entry).CurrentUsage
	}).SumFloats()
	assert.InEpsilon(t, 19.895473, sum, 0.1, "expect sum of aggregated data to equal actual")
}

func TestAggregator1D(t *testing.T) {
	entries := testutil.LoadEntriesFromCSV(t, "../testdata/rawusage-2021-03-26_01.csv")
	agg, err := Aggregate1D(entries)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(agg), "expect 1 entry for a day")

	sum := linq.From(entries).Select(func(i interface{}) interface{} {
		return i.(model.Entry).CurrentUsage
	}).SumFloats()
	assert.InEpsilon(t, 19.895473, sum, 0.1, "expect sum of aggregated data to equal actual")
}
