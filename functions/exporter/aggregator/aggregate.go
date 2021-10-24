package aggregator

import (
	"github.com/ahmetb/go-linq/v3"
	"github.com/eriklupander/powertracker/functions/exporter/model"
	"time"
)

func Aggregate(entries []model.Entry, aggregate string) ([]model.Entry, error) {

	switch aggregate {

	case "1h":
		return Aggregate1h(entries)
	case "1d":
		return Aggregate1D(entries)
	case "1M":
		return Aggregate1M(entries)
	case "5m":
		fallthrough
	default:
		return Aggregate5m(entries), nil
	}
}

// Aggregate5m assumes per 5m input and just multiplies the currentUsage by 12000 to get power in watts.
func Aggregate5m(entries []model.Entry) []model.Entry {
	for i := range entries {
		entries[i].CurrentUsage = entries[i].CurrentUsage * 12 * 1000 // convert to per-hour in watts
	}
	return entries
}

// Aggregate1h assumes 5m input data, groups entries by day+hour and outputs the average power
// for each hour of the time series.
func Aggregate1h(entries []model.Entry) ([]model.Entry, error) {
	output := make([]model.Entry, 0)

	// assume 5m input
	linq.From(entries).
		GroupByT(func(e model.Entry) string {
			return e.Created.Local().Format("2006-01-02 15") + ":00"
		}, func(e model.Entry) float64 {
			return e.CurrentUsage
		}).
		SortT(func(i, j linq.Group) bool {
			return i.Key.(string) < j.Key.(string)
		}).
		ForEachT(func(group linq.Group) {
			ts, err := time.Parse("2006-01-02 15:04", group.Key.(string))
			if err != nil {
				panic(err.Error())
			}
			output = append(output, model.Entry{
				CurrentUsage: linq.From(group.Group).SumFloats() * 1000,
				Created:      ts,
			})
		})
	return output, nil
}

// Aggregate1D assumes 5m input and groups entries per (local time) date, outputting the average
// power for each day. Note - not the cumulative energy used for the full day - the average power in watts
// during the date.
func Aggregate1D(entries []model.Entry) ([]model.Entry, error) {
	output := make([]model.Entry, 0)
	// assume 5m input
	linq.From(entries).
		GroupByT(func(e model.Entry) string {
			return e.Created.Local().Format("2006-01-02")
		}, func(e model.Entry) float64 {
			return e.CurrentUsage
		}).
		SortT(func(i, j linq.Group) bool {
			return i.Key.(string) < j.Key.(string)
		}).
		ForEachT(func(group linq.Group) {
			ts, err := time.Parse("2006-01-02", group.Key.(string))
			if err != nil {
				panic(err.Error())
			}
			output = append(output, model.Entry{
				CurrentUsage: linq.From(group.Group).SumFloats() / float64(len(group.Group)) * 12 * 1000,
				Created:      ts,
			})
		})
	return output, nil
}

func Aggregate1M(entries []model.Entry) ([]model.Entry, error) {
	output := make([]model.Entry, 0)
	// assume 5m input
	linq.From(entries).
		GroupByT(func(e model.Entry) string {
			return e.Created.Local().Format("2006-01")
		}, func(e model.Entry) float64 {
			return e.CurrentUsage
		}).
		SortT(func(i, j linq.Group) bool {
			return i.Key.(string) < j.Key.(string)
		}).
		ForEachT(func(group linq.Group) {
			ts, err := time.Parse("2006-01", group.Key.(string))
			if err != nil {
				panic(err.Error())
			}
			output = append(output, model.Entry{
				CurrentUsage: linq.From(group.Group).SumFloats(),
				Created:      ts,
			})
		})
	return output, nil
}

