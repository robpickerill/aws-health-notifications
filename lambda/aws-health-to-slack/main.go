package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rpickerill/aws-health-to-slack/internal/health"
	"github.com/rpickerill/aws-health-to-slack/internal/notifiers/slack"
)

type Notifier interface {
	Notify(event health.HealthEventDetail) error
}

func LambdaHandler(ctx context.Context, event health.HealthEvent) {
	s, err := slack.New()
	if err != nil {
		log.Fatalf("%s", err)
	}
	s.Notify(event)
}

func main() {
	lambda.Start(LambdaHandler)
}
