package testutil

import (
	"bytes"
	c "encoding/csv"
	"github.com/eriklupander/powertracker/functions/exporter/model"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"strconv"
	"testing"
	"time"
)

func LoadEntriesFromCSV(t *testing.T, filename string) []model.Entry {
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
