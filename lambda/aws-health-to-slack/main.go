package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kelseyhightower/envconfig"
	"github.com/rpickerill/aws-health-to-slack/internal/health"
	"github.com/rpickerill/aws-health-to-slack/internal/notifiers/slack"
)

var conf Config

type Config struct {
	Debug             bool
	SlackWebHook      string        ``
	SlackUserName     string        `default:"AWS Health Notifications"`
	SlackTimeOut      time.Duration `default:"2s"`
	PagerdutyApiToken string
}

func init() {
	envconfig.Process("", &conf)
}

func LambdaHandler(ctx context.Context, event health.HealthEvent) error {

	errorChan := make(chan error)
	wgChan := make(chan bool)

	waitGroup := sync.WaitGroup{}

	if len(conf.SlackWebHook) > 0 {
		s := slack.NewSlackClient(ctx, conf.SlackWebHook, conf.SlackUserName, conf.SlackTimeOut)
		waitGroup.Add(1)
		go s.Notify(ctx, &waitGroup, errorChan, event)
	}

	go func() {
		waitGroup.Wait()
		close(wgChan)
	}()

	var errorCount = 0
	select {
	case <-wgChan:
		break
	case err := <-errorChan:
		log.Printf("%s", err)
		errorCount++
	}

	close(errorChan)
	switch errorCount {
	case 0:
		log.Printf("successfully completed")
	default:
		return fmt.Errorf("completed with %d errors, see logs", errorCount)
	}

	return nil
}

func main() {
	lambda.Start(LambdaHandler)
}
