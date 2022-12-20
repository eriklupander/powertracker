package main

import (
	"fmt"
	"github.com/eriklupander/powertracker/functions/statusrecorder/model"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func Test_ScrapeAll(t *testing.T) {
	mockWriter := NewMockWriter()
	assert.NoError(t, scrapeChargers(mockWriter))

	fmt.Printf("%v\n", mockWriter.entries)
}

type MockWriter struct {
	lock sync.Mutex
	entries []string
}

func NewMockWriter() *MockWriter {
	return &MockWriter{entries: make([]string, 0)}
}

func (m *MockWriter) Write(r model.Record) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.entries = append(m.entries, fmt.Sprintf("%v", r))
}

func (m *MockWriter) Flush() {

}
