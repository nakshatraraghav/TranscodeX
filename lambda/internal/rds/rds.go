package rdsservice

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/nakshatraraghav/transcodex/lambda/config"
)

type RDSService struct {
	client *rds.RDS
}

func NewRDSService(session *session.Session) *RDSService {
	client := rds.New(session)

	return &RDSService{
		client: client,
	}
}

func (r *RDSService) GetRDSEndpoint() (string, error) {
	env := config.Getenv()

	response, err := r.client.DescribeDBInstances(&rds.DescribeDBInstancesInput{
		DBInstanceIdentifier: aws.String(env.DATABASE_INSTANCE_IDENTIFIER),
	})
	if err != nil {
		return "", err
	}

	return *response.DBInstances[0].Endpoint.Address, nil
}

func ConstructConnectionString(username, password, host string) string {
	port := "5432"
	database := "transcodex"

	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", username, password, host, port, database)
}
