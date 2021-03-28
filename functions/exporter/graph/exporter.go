package graph

import (
	"bytes"
	"github.com/eriklupander/powertracker/functions/exporter/model"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

func Export(entries []model.Entry) ([]byte, error) {
	xticks := plot.TimeTicks{Format: "2006-01-02\n15:04"}
	pts := make(plotter.XYs, len(entries))
	for i := range pts {
		pts[i].X = float64(entries[i].Created.Unix())
		pts[i].Y = entries[i].CurrentUsage
	}

	p := plot.New()

	p.Title.Text = "Power usage"
	p.X.Label.Text = "Time"
	p.X.Tick.Marker = xticks
	p.Y.Label.Text = "Power"

	err := plotutil.AddLines(p, "Watts", pts)
	if err != nil {
		return nil, err
	}

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
