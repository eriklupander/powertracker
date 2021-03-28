package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/labstack/gommon/log"
)

const tibberConfigKey = "prod/tibber_config"
const tibberHomeIdKey = "tibber_home_id"
const tibberApiKeyKey = "tibber_api_key"

var tibberApiKey, tibberHomeId string

// handler is the function called by the lambda.
func handler(ctx context.Context) error {
	if err := validateConfig(); err != nil {
		return err
	}

	// connect to watty
	if err := connectToWatty(tibberApiKey, tibberHomeId); err != nil {
		log.Error(err.Error())
	}

	return nil
}

// main is called when a new lambda starts, so don't
// expect to have something done for every query here.
func main() {
	fmt.Println("init power recorder")
	var err error
	tibberConfigJSON, err := getSecret(tibberConfigKey)
	if err != nil {
		panic(err.Error())
	}
	config := make(map[string]interface{})
	if err := json.Unmarshal([]byte(tibberConfigJSON), &config); err != nil {
		panic(err.Error())
	}
	var ok bool
	tibberApiKey, ok = config[tibberApiKeyKey].(string)
	if !ok {
		panic("unable to resolve tibber_api_key from JSON")
	}
	tibberHomeId, ok = config[tibberHomeIdKey].(string)
	if !ok {
		panic("unable to resolve tibber_home_id from JSON")
	}

	lambda.StartWithContext(context.Background(), handler)
}
