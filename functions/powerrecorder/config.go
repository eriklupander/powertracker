package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
)

func configure() {
	tibberConfigJSON, err := getSecret(tibberConfigKey)
	if err != nil {
		logrus.WithError(err).Fatal("error getting secrets from AWS Secrets Manager")
	}
	config := make(map[string]interface{})
	if err := json.Unmarshal([]byte(tibberConfigJSON), &config); err != nil {
		logrus.WithError(err).Fatal("error unmarshalling tibber config JSON")
	}
	var ok bool
	tibberApiKey, ok = config[tibberApiKeyKey].(string)
	if !ok {
		logrus.Fatal("unable to resolve tibber_api_key from JSON")
	}
	tibberHomeId, ok = config[tibberHomeIdKey].(string)
	if !ok {
		logrus.Fatal("unable to resolve tibber_home_id from JSON")
	}
}

func validateConfig() error {
	if err := validate(tibberApiKeyKey, tibberApiKey); err != nil {
		return err
	}
	if err := validate(tibberHomeIdKey, tibberHomeId); err != nil {
		return err
	}
	return nil
}
func validate(key, value string) error {
	if value == "" {
		errMsg := fmt.Sprintf("No value configured in AWS Secrets Manager using key '%s'! Cannot execute.", key)
		return errors.New(errMsg)
	}
	return nil
}
