package main

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kelseyhightower/envconfig"
	"github.com/rpickerill/aws-health-to-slack/internal/health"
	"github.com/rpickerill/aws-health-to-slack/internal/notifiers/slack"
)

type Notifier interface {
	Notify(event health.HealthEventDetail) error
}

type Config struct {
	Debug         bool
	SlackWebohook string
	SlackUsername string
	SlackTimeout  time.Duration
}

func LambdaHandler(ctx context.Context, event health.HealthEvent) {
	var conf Config
	err := envconfig.Process("", &conf)
	if err != nil {
		log.Fatalf("%s", err)
	}

	s := slack.NewSlackClient(conf.SlackWebohook, conf.SlackUsername, conf.SlackTimeout)
	s.Notify(event)
}

func main() {
	lambda.Start(LambdaHandler)
}
