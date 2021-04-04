package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sirupsen/logrus"
)

const tibberConfigKey = "prod/tibber_config"
const tibberHomeIdKey = "tibber_home_id"
const tibberApiKeyKey = "tibber_api_key"

var tibberApiKey, tibberHomeId string

// handler is the function called when the lambda is invoked, i.e. by the Event Bridge event in our case.
func handler(ctx context.Context) error {
	if err := validateConfig(); err != nil {
		return err
	}

	// connect to watty
	if err := recordPowerUsageFromWatty(tibberApiKey, tibberHomeId); err != nil {
		logrus.WithError(err).Error("error recording power usage from Watty")
	}

	return nil
}

// main is called when a new lambda starts, so don't
// expect to have something done for every query here.
func main() {
	logrus.Info("init power recorder")

	// load secrets etc, will panic on errors.
	configure()

	lambda.StartWithContext(context.Background(), handler)
}
