package handler

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/nakshatraraghav/transcodex/lambda/config"
	ecsservice "github.com/nakshatraraghav/transcodex/lambda/internal/ecs"
)

func HandlerLocal(ctx context.Context, event events.SQSEvent) error {
	env := config.Getenv()

	session, err := session.NewSession(&aws.Config{
		Region: aws.String(env.REGION_STRING),
	})
	if err != nil {
		return err
	}

	ecs := ecsservice.NewECSService(session)

	for _, record := range event.Records {
		err := ecs.RunTask(record, env.CONNECTION_STRING)
		if err != nil {
			return err
		}
	}

	return nil
}
