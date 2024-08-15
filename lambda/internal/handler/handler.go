package handler

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
)

func Handler(ctx context.Context, event events.SQSEvent) error {
	return nil
}
