package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kelseyhightower/envconfig"
	"github.com/rpickerill/aws-health-to-slack/internal/health"
	"github.com/rpickerill/aws-health-to-slack/internal/notifiers/slack"
)

type Config struct {
	Debug         bool
	SlackWebHook  string
	SlackUserName string
	SlackTimeOut  time.Duration
}

func LambdaHandler(ctx context.Context, event health.HealthEvent) {
	var conf Config
	envconfig.MustProcess("", &conf)

	s := slack.NewSlackClient(ctx, conf.SlackWebHook, conf.SlackUserName, conf.SlackTimeOut)

	errorChan := make(chan error)
	wgChan := make(chan bool)

	waitGroup := sync.WaitGroup{}

	waitGroup.Add(1)
	go s.Notify(ctx, &waitGroup, event)

	go func() {
		waitGroup.Wait()
		close(wgChan)
	}()

	select {
	case <-wgChan:
		log.Printf("Successfully sent notifications, cleaning up")
		break
	case err := <-errorChan:
		log.Printf("%s", err)
	}

	log.Println("Successfully Completed.")
}

func main() {
	lambda.Start(LambdaHandler)
}
