# Chargescraper
Loads data from ChargeFinder API and stores in InfluxDB.

Purely a hobby project - not intended for commercial use!!!

## Building and running
In order to NOT infringe on any ChargeFinder terms of usage, the `sites.go` file containing site metadata is not commited to the repository.

To successfully compile this application, please create your own `sites.go` using the template below:
```go
package main

import "github.com/eriklupander/powertracker/functions/statusrecorder/model"

var chargeFinderSites = []model.Site{
	{
		Id:           "id-of-the-charging-site",
		Name:         "Name of the Charging Site",
		SiteProvider: "Not used",
		Address:      "Not used",
		Data:         `Fetch from ChargeFinder, using path param: https://api.chargefinder.com/station/<siteId>`,
	},
}
```

I suggest going to ChargeFinder.com and use the DevTools in your browser to see site-id's on requests that's executed when a site is clicked.