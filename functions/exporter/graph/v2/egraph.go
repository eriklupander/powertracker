package v2

import (
	"bytes"
	"github.com/ahmetb/go-linq/v3"
	"github.com/eriklupander/powertracker/functions/exporter/model"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
	"github.com/pkg/errors"
	"math"
)

func BarChart(entries []model.Entry) ([]byte, error) {
	aggrTitle, timeFormat := resolveAggregation(entries)
	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{PageTitle: "Powertracker"}),
		charts.WithTitleOpts(opts.Title{
			Title:    "Average power " + aggrTitle,
			Subtitle: "In watts",
		}),
		charts.WithTooltipOpts(opts.Tooltip{Show: true}),
		charts.WithLegendOpts(opts.Legend{Right: "80%"}),
	)

	xData := make([]string, 0)
	barData := make([]opts.BarData, 0)
	for _, ex := range entries {
		e := ex
		xData = append(xData, e.Created.Format(timeFormat))
		barData = append(barData, opts.BarData{Value: toFixed(e.CurrentUsage, 2)})
	}

	// Put data into instance
	bar.SetXAxis(xData).
		AddSeries("Power (Wattage)", barData)

	out := new(bytes.Buffer)
	err := bar.Render(out)
	if err != nil {
		return nil, errors.Wrap(err, "rendering barchart")
	}
	return out.Bytes(), nil
}

func BarChartOverTime(entries []model.Entry) ([]byte, error) {
	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{PageTitle: "Powertracker"}),
		charts.WithTitleOpts(opts.Title{
			Title:    "Average power consumption grouped per hour",
			Subtitle: "In watts",
		}),
		charts.WithTooltipOpts(opts.Tooltip{Show: true}),
		charts.WithLegendOpts(opts.Legend{Show: true, Right: "80%"}),
	)

	xData := make([]string, 0)
	barData := make([]opts.BarData, 0)

	// Example using go-linq's psuedo-generic API (uses reflection behind the scenes and is much slower)
	linq.From(entries).
		GroupByT(func(e model.Entry) string {
			return e.Created.Local().Format("15") + ":00"
		}, func(e model.Entry) float64 {
			return e.CurrentUsage
		}).
		SortT(func(i, j linq.Group) bool {
			return i.Key.(string) < j.Key.(string)
		}).
		ForEachT(func(group linq.Group) {
			xData = append(xData, group.Key.(string))
			barData = append(barData, opts.BarData{Value: toFixed(linq.From(group.Group).Average(), 2)})
		})

	bar.SetXAxis(xData).AddSeries("Accumulated hourly power (Watts)", barData)

	out := new(bytes.Buffer)
	err := bar.Render(out)
	if err != nil {
		return nil, errors.Wrap(err, "rendering over time barchart")
	}
	return out.Bytes(), nil
}

func LineChart(entries []model.Entry) ([]byte, error) {
	aggrTitle, timeFormat := resolveAggregation(entries)

	xData := make([]string, 0)
	lineData := make([]opts.LineData, 0)
	for _, ex := range entries {
		e := ex
		xData = append(xData, e.Created.Format(timeFormat))
		lineData = append(lineData, opts.LineData{Value: toFixed(e.CurrentUsage, 2)})
	}

	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeRoma, PageTitle: "Powertracker"}),
		charts.WithTitleOpts(opts.Title{
			Title:    "Average power " + aggrTitle,
			Subtitle: "In watts",
		}),
		charts.WithTooltipOpts(opts.Tooltip{Show: true, Trigger: "axis"}),
	)

	line.SetXAxis(xData).
		AddSeries("Power usage (Watts)", lineData).
		SetSeriesOptions(
			charts.WithMarkPointNameTypeItemOpts(
				opts.MarkPointNameTypeItem{Name: "Maximum", Type: "max"},
				opts.MarkPointNameTypeItem{Name: "Average", Type: "average"},
				opts.MarkPointNameTypeItem{Name: "Minimum", Type: "min"},
			),
			charts.WithLineChartOpts(opts.LineChart{Smooth: true}),
			charts.WithMarkPointStyleOpts(
				opts.MarkPointStyle{
					Symbol:     []string{"diamond"},
					Label:      &opts.Label{Show: true},
					SymbolSize: 30.0,
				}),
		)

	out := new(bytes.Buffer)
	err := line.Render(out)
	if err != nil {
		return nil, errors.Wrap(err, "rendering line chart")
	}
	return out.Bytes(), nil
}
func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}
func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func resolveAggregation(entries []model.Entry) (string, string) {
	first := entries[0].Created
	last := entries[len(entries)-1].Created
	timeBetween := last.Sub(first)
	minutesPerEntry := timeBetween.Minutes() / float64(len(entries))
	if minutesPerEntry < 6 {
		return "per 5 minutes", "060102 15:04"
	} else if minutesPerEntry < 62 {
		return "per hour", "060102 15:04"
	} else if minutesPerEntry < 60*24+24 {
		return "per day", "060102"
	} else {
		return "per month", "0601"
	}
}
