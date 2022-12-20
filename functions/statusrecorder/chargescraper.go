package main

import (
	"context"
	"github.com/eriklupander/powertracker/functions/statusrecorder/model"
	"github.com/eriklupander/powertracker/functions/statusrecorder/provider"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

func scrapeChargers(influxWriter RecordWriter) error {
	ctx, cfn := context.WithTimeout(context.Background(), time.Second*20)
	defer cfn()

	chargeFinderProvider := provider.NewChargeFinderProvider()

	wg := sync.WaitGroup{}
	wg.Add(len(chargeFinderSites))
	for _, site := range chargeFinderSites {

		go collect(ctx, site, &wg, influxWriter, chargeFinderProvider)

		// add a slight artificial stagger to avoid tripping any rate limiters at Chargefinder
		time.Sleep(time.Millisecond * 100)
	}
	wg.Wait()
	logrus.Info("scraping done!")

	return nil
}

func collect(ctx context.Context, site model.Site, wg *sync.WaitGroup, influxWriter RecordWriter, provider Provider) {
	st := time.Now()

	r, err := provider.LoadSite(ctx, site)
	if err != nil {
		logrus.WithError(err).Infof("error collecting data from site %v", site.Id)
		wg.Done()
		return
	}
	r.SiteId = site.Id
	r.SiteName = site.Name
	influxWriter.Write(r)
	logrus.WithFields(logrus.Fields{
		"site":      site.Name,
		"available": r.Available,
		"max":       r.Total,
		"duration":  time.Since(st).String(),
	}).Infof("scraped site %s", site.Id)
	wg.Done()
}

// Provider is a call-site declared interface for accessing site data
type Provider interface {
	LoadSite(ctx context.Context, site model.Site) (model.Record, error)
}
