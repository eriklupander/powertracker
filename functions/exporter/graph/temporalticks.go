package graph

import (
	"github.com/sirupsen/logrus"
	"gonum.org/v1/plot"
	"time"
)

// UTCDateTimeTicks
type UTCDateTimeTicks struct {
	Steps float64
}

func NewUTCDateTimeTicks(steps float64) UTCDateTimeTicks {
	if steps == 0.0 {
		steps = 6.0
	}
	return UTCDateTimeTicks{Steps: steps}
}

// UTCDateTimeTicks returns datetime ticks in UTC in the specified range.
func (u UTCDateTimeTicks) Ticks(min, max float64) []plot.Tick {
	if max <= min {
		logrus.Fatal("illegal range")
	}
	start := time.Unix(int64(min), 0)
	end := time.Unix(int64(max), 0)
	end = end.Add(time.Hour)               // make last our "inclusive" so a day's duration becomes 24 h

	stepHours :=  end.Sub(start).Hours() / u.Steps
	// failsafe
	if stepHours < 1.0 {
		stepHours = 1.0
	}
	logrus.Infof("computed %v stepHours", stepHours)
	var ticks []plot.Tick
	d := time.Unix(start.UTC().Unix(), 0).Truncate(time.Minute)
	for d.Before(end) {
		ticks = append(ticks, plot.Tick{Value: float64(d.UTC().Unix()), Label: d.UTC().Format("2006-01-02 15:04")})
		d = d.Add(time.Duration(stepHours) * time.Hour)
	}
	return ticks
}

// DateTimeTicks
type DateTimeTicks struct {
	Steps float64
}

func NewDateTimeTicks(steps float64) DateTimeTicks {
	if steps == 0.0 {
		steps = 6.0
	}
	return DateTimeTicks{Steps: steps}
}

// Ticks returns Ticks in the specified range where we try to "even out" things slightly.
func (dtt DateTimeTicks) Ticks(min, max float64) []plot.Tick {
	if max <= min {
		logrus.Fatal("illegal range")
	}
	start := time.Unix(int64(min), 0)
	end := time.Unix(int64(max), 0)
	end = end.Add(time.Hour) // make so our end "inclusive" so a day's duration becomes 24 h

	//
	stepHours := end.Sub(start).Hours() / dtt.Steps

	var ticks []plot.Tick
	d := time.Unix(start.Unix(), 0).Truncate(time.Minute)
	for d.Before(end) {
		ticks = append(ticks, plot.Tick{Value: float64(d.Unix()), Label: d.Format("2006-01-02 15:04")})
		d = d.Add(time.Duration(stepHours) * time.Hour)
	}
	return ticks
}
