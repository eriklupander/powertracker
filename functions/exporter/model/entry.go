package model

import "time"

type Entry struct {
	HomeId           string
	CurrentUsage     float64
	AccumulatedDaily float64
	Created          time.Time
}
