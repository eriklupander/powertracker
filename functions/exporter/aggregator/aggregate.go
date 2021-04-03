package aggregator

import (
	"fmt"
	"github.com/eriklupander/powertracker/functions/exporter/model"
	"sort"
	"time"
)

func Aggregate(entries []model.Entry, aggregate string) ([]model.Entry, error) {

	layout := ""
	divider := 12.0
	switch aggregate {

	case "1h":
		layout = "2006-01-02 15"
		divider = 1.0
	case "1d":
		layout = "2006-01-02"
		divider = 1.0 / 24.0
	case "1M":
		layout = "2006-01"
		divider = 1.0 / 24.0 / 30.0
	case "5m":
		fallthrough
	default:
		// do nothing, this is the default
		layout = "2006-01-02 15:04"
	}

	if layout != "" {
		// transform entries to "per hour" instead
		// assume already sorted ASC
		m := make(map[string][]model.Entry, 0)
		for i := range entries {
			t := entries[i].Created.Format(layout)
			_, ok := m[t]
			if !ok {
				m[t] = make([]model.Entry, 0)
			}
			m[t] = append(m[t], entries[i])
		}

		// Now, sum each bucket and stuff into new list
		out := make([]model.Entry, 0)
		for k, v := range m {
			sum := 0.0
			for i := range v {
				sum += v[i].CurrentUsage
			}
			created, err := time.Parse(layout, k)
			if err != nil {
				fmt.Printf("error parsing created date after aggregation: %s error: %v\n", k, err)
				return nil, err
			}

			// sum needs to be adjusted to "effect".
			// If we've consumed 0.2kW in 5 minutes, that should translate to 60/5 x 0.2kW x 1000 to get effect in W
			sum = divider * sum * 1000
			out = append(out, model.Entry{CurrentUsage: sum, HomeId: v[0].HomeId, Created: created})
		}
		entries = out
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Created.Before(entries[j].Created)
	})

	return entries, nil
}
