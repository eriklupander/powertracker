package main

import (
	"github.com/eriklupander/powertracker/functions/exporter/model"
)

type DataSource interface {
	GetAll(from, to string) ([]model.Entry, error)
}
