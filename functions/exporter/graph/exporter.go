package graph

import (
	"bytes"
	"github.com/ahmetb/go-linq"
	"github.com/eriklupander/powertracker/functions/exporter/model"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

func ExportLinePlot(entries []model.Entry) ([]byte, error) {
	pts := make(plotter.XYs, len(entries))
	for i := range pts {
		pts[i].X = float64(entries[i].Created.Unix())
		pts[i].Y = entries[i].CurrentUsage
	}

	p := plot.New()
	p.Title.Text = "Energy usage"
	p.X.Label.Text = "Time"
	p.X.Min = float64(entries[0].Created.Unix())
	p.X.Tick.Marker = NewDateTimeTicks(6)
	p.Y.Label.Text = "Energy (Watts)"
	p.Y.Min = 0

	err := plotutil.AddLines(p, "Watts", pts)
	if err != nil {
		return nil, err
	}

	return toPNG(p)
}

func ExportHist(entries []model.Entry) ([]byte, error) {
	var pts plotter.XYs
	linq.From(entries).Select(func(i interface{}) interface{} {
		entry := i.(model.Entry)
		return plotter.XY{X: float64(entry.Created.Unix()), Y: entry.CurrentUsage}
	}).ToSlice(&pts)

	hist, err := plotter.NewHistogram(pts, len(pts))
	if err != nil {
		return nil, err
	}
	p := plot.New()
	p.Title.Text = "Energy usage"

	p.X.Label.Text = "Time (UTC)"
	p.X.Min = float64(entries[0].Created.Unix())
	p.X.Tick.Marker = NewUTCDateTimeTicks(6.0)

	p.Y.Label.Text = "Energy (Watts)"
	p.Y.Min = 0.0
	p.Y.Max = linq.From(entries).Select(func(i interface{}) interface{} {
		return i.(model.Entry).CurrentUsage
	}).Max().(float64)

	p.Add(hist)

	return toPNG(p)
}

func toPNG(p *plot.Plot) ([]byte, error) {
	buf := new(bytes.Buffer)
	writerTo, err := p.WriterTo(12*vg.Inch, 3*vg.Inch, "png")
	if err != nil {
		return nil, err
	}
	_, err = writerTo.WriteTo(buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}