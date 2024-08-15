package secretsmanager

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

type SecretsManagerService struct {
	client *secretsmanager.SecretsManager
}

func NewSecretsManagerService(session *session.Session) *SecretsManagerService {
	sm := secretsmanager.New(session)

	return &SecretsManagerService{
		client: sm,
	}
}

func (s *SecretsManagerService) Get(key string) (string, error) {
	response, err := s.client.GetSecretValue(&secretsmanager.GetSecretValueInput{
		SecretId: aws.String(key),
	})

	if err != nil {
		return "", err
	}

	secret := *response.SecretString

	return secret, nil

}

func (s *SecretsManagerService) Put(key, value string) error {
	_, err := s.client.PutSecretValue(&secretsmanager.PutSecretValueInput{
		SecretId:     aws.String(key),
		SecretString: aws.String(value),
	})

	if err != nil {
		return err
	}

	return nil

}
