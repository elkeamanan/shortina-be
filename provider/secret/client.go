package secret

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type SecretsManager interface {
	GetSecret(ctx context.Context, secretName string) (string, error)
}

type AWSSecretsManager struct {
	client *secretsmanager.Client
}

func NewAWSSecretsManager(ctx context.Context, region string) (SecretsManager, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := secretsmanager.NewFromConfig(cfg)
	return &AWSSecretsManager{client: client}, nil
}

func (s *AWSSecretsManager) GetSecret(ctx context.Context, secretName string) (string, error) {
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	}

	result, err := s.client.GetSecretValue(ctx, input)
	if err != nil {
		return "", fmt.Errorf("failed to get secret %s: %w", secretName, err)
	}

	if result.SecretString == nil {
		return "", fmt.Errorf("secret %s has no string value", secretName)
	}

	return *result.SecretString, nil
}
