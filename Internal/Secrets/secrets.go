package Secrets

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type DBCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SecretsFetcher struct {
	client *secretsmanager.Client
}

func NewSecretsFetcher() (*SecretsFetcher, error) {
	awsCfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	return &SecretsFetcher{
		client: secretsmanager.NewFromConfig(awsCfg),
	}, nil
}

func (s *SecretsFetcher) FetchDBCredentials(secretARN string) (*DBCredentials, error) {
	result, err := s.client.GetSecretValue(context.TODO(), &secretsmanager.GetSecretValueInput{
		SecretId: &secretARN,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch secret: %w", err)
	}

	var creds DBCredentials
	if err := json.Unmarshal([]byte(*result.SecretString), &creds); err != nil {
		return nil, fmt.Errorf("failed to parse secret: %w", err)
	}

	return &creds, nil
}
