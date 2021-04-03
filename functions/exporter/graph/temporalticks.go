package graph

import (
	"gonum.org/v1/plot"
	"time"
)

// UTCDateTimeTicks
type UTCDateTimeTicks struct{

}

// UTCDateTimeTicks returns datetime ticks in UTC in the specified range.
func (UTCDateTimeTicks) Ticks(min, max float64) []plot.Tick {
	if max <= min {
		panic("illegal range")
	}
	start := time.Unix(int64(min), 0) //.Truncate(time.Hour)
	end := time.Unix(int64(max), 0)   //.Round(time.Hour)
	end = end.Add(time.Hour)          // make last our "inclusive" so a day's duration becomes 24 h
	duration := end.Sub(start)
	totalHours := duration.Hours()
	stepHours := totalHours / 6
	var ticks []plot.Tick

	d := time.Unix(start.UTC().Unix(), 0).Truncate(time.Minute)
	for d.Before(end) {
		ticks = append(ticks, plot.Tick{Value: float64(d.UTC().Unix()), Label: d.UTC().Format("2006-01-02 15:04")})
		d = d.Add(time.Duration(stepHours) * time.Hour)
	}
	return ticks
}


// DateTimeTicks
type DateTimeTicks struct{}

// Ticks returns Ticks in the specified range.
func (DateTimeTicks) Ticks(min, max float64) []plot.Tick {
	if max <= min {
		panic("illegal range")
	}
	start := time.Unix(int64(min), 0) //.Truncate(time.Hour)
	end := time.Unix(int64(max), 0)   //.Round(time.Hour)
	end = end.Add(time.Hour)          // make so our "inclusive" so a day's duration becomes 24 h
	duration := end.Sub(start)
	totalHours := duration.Hours()
	stepHours := totalHours / 6
	var ticks []plot.Tick
	d := time.Unix(start.Unix(), 0).Truncate(time.Minute)
	for d.Before(end) {
		ticks = append(ticks, plot.Tick{Value: float64(d.Unix()), Label: d.Format("2006-01-02 15:04")})
		d = d.Add(time.Duration(stepHours) * time.Hour)
	}
	return ticks
}
