package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/nakshatraraghav/transcodex/lambda/config"
	"github.com/nakshatraraghav/transcodex/lambda/internal/handler"
)

func main() {
	config.LoadEnv()
	lambda.Start(handler.HandlerLocal)
}
