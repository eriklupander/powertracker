package provider

import (
	"context"
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/eriklupander/powertracker/functions/statusrecorder/model"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	statusFree     = 2
	statusOccupied = 5
)

var noData = fmt.Errorf("no data")

// ChargeFinderProvider is responsible for retrieving and parsing charger data from the Chargefinder service.
type ChargeFinderProvider struct {
	client *http.Client
}

func NewChargeFinderProvider() *ChargeFinderProvider {
	return &ChargeFinderProvider{client: &http.Client{}}
}

func (cfc *ChargeFinderProvider) LoadSite(ctx context.Context, site model.Site) (model.Record, error) {

	// find CCS chargers
	ccsChargers := findCCSChargersOnSite([]byte(site.Data))

	if len(ccsChargers) == 0 {
		return model.Record{}, noData
	}

	req, err := http.NewRequest("GET", "https://adm.chargefinder.com/status/"+site.Id, nil)
	if err != nil {
		return model.Record{}, err
	}
	req = req.WithContext(ctx)

	resp, err := cfc.client.Do(req)
	if err != nil {
		return model.Record{}, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return model.Record{}, err
	}
	return parseChargers(ccsChargers, data)
}

func findCCSChargersOnSite(data []byte) []model.CCSCharger {

	out := make([]model.CCSCharger, 0)

	_, err := jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		plugType, _ := jsonparser.GetString(value, "plug")
		if plugType == "CCS" {
			// get the outlets
			_, err = jsonparser.ArrayEach(value, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
				identifier, _ := jsonparser.GetString(value, "identifier")
				name, _ := jsonparser.GetString(value, "name")
				capacity, _ := jsonparser.GetInt(value, "capacity")
				out = append(out, model.CCSCharger{
					Identifier: identifier,
					Name:       name,
					Capacity:   capacity,
				})
			}, "outlets")
		}
	}, "outletList")

	if err != nil {
		return nil
	}
	return out
}

func parseChargers(ccsChargers []model.CCSCharger, data []byte) (model.Record, error) {
	available := 0
	_, err := jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		id, _ := jsonparser.GetString(value, "id")
		for _, ccsCharger := range ccsChargers {
			if id == ccsCharger.Identifier || id == ccsCharger.Name {
				status, _ := jsonparser.GetInt(value, "status")
				if status == statusFree {
					available++
				}
				break
			}
		}
	})
	if err != nil {
		return model.Record{}, err
	}
	return model.Record{Available: available, Total: len(ccsChargers)}, nil
}

func loadSiteMetadata(id string, dump bool) ([]byte, error) {
	resp, err := http.Get("https://api.chargefinder.com/station/" + id)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// dump files to disk
	if dump {
		if err := ioutil.WriteFile(id+".json", data, os.FileMode(0755)); err != nil {
			logrus.WithError(err).Info("error dumping file to disk")
		}
	}
	return data, nil
}
