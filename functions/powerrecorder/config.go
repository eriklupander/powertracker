package main

import (
	"errors"
	"fmt"
)

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
