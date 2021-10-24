package main

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
)

const configKey = "prod/tibber_config"
const influxDbTokenKey = "influxdb_cloud_api_token"

var influxTbToken string

func configure() {
	configJSON, err := getSecret(configKey)
	if err != nil {
		logrus.WithError(err).Fatal("error getting secrets from AWS Secrets Manager")
	}
	config := make(map[string]interface{})
	if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
		logrus.WithError(err).Fatal("error unmarshalling config JSON")
	}
	var ok bool
	influxTbToken, ok = config[influxDbTokenKey].(string)
	if !ok {
		logrus.Fatalf("unable to resolve %s from JSON", influxDbTokenKey)
	}
}
