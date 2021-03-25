package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Notifier interface {
	Notify(message string)
}

func LambdaHandler(ctx context.Context, event events.CloudWatchEvent) error {

	return nil
}

func main() {
	lambda.Start(LambdaHandler)
}
