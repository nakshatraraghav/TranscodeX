package handler

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/nakshatraraghav/transcodex/lambda/config"
	ecsservice "github.com/nakshatraraghav/transcodex/lambda/internal/ecs"
	rdsservice "github.com/nakshatraraghav/transcodex/lambda/internal/rds"
)

func Handler(ctx context.Context, event events.SQSEvent) error {
	env := config.Getenv()

	session, err := session.NewSession(&aws.Config{
		Region: aws.String(env.REGION_STRING),
	})
	if err != nil {
		return err
	}

	rds := rdsservice.NewRDSService(session)
	ecs := ecsservice.NewECSService(session)

	// Step 1: Retrieve the connection string from aws secrets manager
	host, err := rds.GetRDSEndpoint()
	if err != nil {
		return err
	}

	CONNECTION_STRING := rdsservice.ConstructConnectionString(
		env.RDS_DATABASE_USERNAME,
		env.RDS_DATABASE_PASSWORD,
		host,
	)

	for _, record := range event.Records {
		ecs.RunTask(record, CONNECTION_STRING)
	}

	return nil
}
