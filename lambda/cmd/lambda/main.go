package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/nakshatraraghav/transcodex/lambda/internal/handler"
)

func main() {
	lambda.Start(handler.Handler)
}
