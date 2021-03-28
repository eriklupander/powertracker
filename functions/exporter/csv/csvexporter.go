package csv

import (
	"bytes"
	"encoding/csv"
	"github.com/eriklupander/powertracker/functions/exporter/model"
	"strconv"
	"time"
)

func Export(entries []model.Entry) ([]byte, error) {
	records := [][]string{
		{"home_id", "current_usage", "acc_usage", "time"},
	}
	for _, e := range entries {
		records = append(records, []string{
			e.HomeId,
			strconv.FormatFloat(e.CurrentUsage, 'f', 6, 64),
			strconv.FormatFloat(e.AccumulatedDaily, 'f', 6, 64),
			e.Created.Format(time.RFC3339)},
		)
	}
	f := new(bytes.Buffer)

	w := csv.NewWriter(f)

	for _, record := range records {
		if err := w.Write(record); err != nil {
			return nil, err
		}
	}
	w.Flush()
	return f.Bytes(), nil
}
