package main

import (
	"fmt"
	"github.com/eriklupander/powertracker/functions/exporter/aggregator"
	"github.com/eriklupander/powertracker/functions/exporter/csv"
	"github.com/eriklupander/powertracker/functions/exporter/graph"
	"github.com/eriklupander/powertracker/functions/exporter/model"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
	"net/http"
)

func setupRouter(source DataSource) *chi.Mux {

	r := chi.NewRouter()
	r.Use(middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: logrus.New()}))

	r.Get("/", handle(source))
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(fmt.Sprintf("The requested path %s was not found", r.RequestURI)))
	})
	return r
}

func handle(source DataSource) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		output := r.URL.Query().Get("output")
		fromStr := r.URL.Query().Get("from")
		toStr := r.URL.Query().Get("to")
		aggregate := r.URL.Query().Get("aggregate")
		graphType := r.URL.Query().Get("graph")

		if output == "" {
			output = "csv"
		}
		if graphType == "" {
			graphType = "hist"
		}

		entries, err := source.GetAll(fromStr, toStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		entries, err = aggregator.Aggregate(entries, aggregate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		switch output {
		case "png":
			exportPNG(w, entries, graphType)
		case "csv":
			fallthrough
		default:
			exportCSV(w, entries)
		}
		logrus.Infof("exported %d entries in %s format between %s and %s\n", len(entries), output, fromStr, toStr)
	}
}

func exportCSV(w http.ResponseWriter, entries []model.Entry) {
	data, err := csv.Export(entries)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

func exportPNG(w http.ResponseWriter, entries []model.Entry, graphType string) {
	var data []byte
	var err error
	switch graphType {
	case "lineplot":
		data, err = graph.ExportLinePlot(entries)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case "hist":
		fallthrough
	default:
		data, err = graph.ExportHist(entries)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "image/png")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}
