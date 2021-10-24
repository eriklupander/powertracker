package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sirupsen/logrus"
)

// handler is the function called when the lambda is invoked, i.e. by the Event Bridge event in our case.
func handler(ctx context.Context) error {

	bucket := "chargerstatus"
	org := "" // ENTER EMAIL HERE!
	client := NewInfluxWriter(bucket, org)

	// scrape chargefinder API
	if err := ScrapeChargers(client); err != nil {
		logrus.WithError(err).Error("error recording charger status")
	}

	return nil
}

// main is called when a new lambda starts, so don't
// expect to have something done for every query here.
func main() {
	logrus.Info("init charger status recorder")

	// load secrets etc, will panic on errors.
	configure()

	lambda.StartWithContext(context.Background(), handler)
}
