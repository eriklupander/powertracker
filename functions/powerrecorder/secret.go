package main

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

func getSecret(secretName string) (string, error) {
	s := session.Must(session.NewSession())
	sm := secretsmanager.New(s)
	output, err := sm.GetSecretValue(&secretsmanager.GetSecretValueInput{SecretId: &secretName})
	if err != nil {
		return "", err
	}
	return *output.SecretString, nil
}
