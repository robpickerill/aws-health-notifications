package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

var slackToken string

type HealthEventDetail struct {
	EventArn          string                        `json:"eventArn"`
	Service           string                        `json:"service"`
	EventTypeCode     string                        `json:"eventTypeCode"`
	EventTypeCategory string                        `json:"eventTypeCategory"`
	StartTime         string                        `json:"startTime"`
	EndTime           string                        `json:"endTime"`
	EventDescription  []HealthEventDescription      `json:"eventDescription"`
	AffectedEntities  []HealthEventAffectedEntities `json:"affectedEntities"`
}

type HealthEventDescription struct {
	Language          string `json:"language"`
	LatestDescription string `json:"latestDescription"`
}

type HealthEventAffectedEntities struct {
	EntityValue string          `json:"entityValue"`
	Tags        json.RawMessage `json:"tags"`
}

func GetParameter(p ssm.GetParameterInput) (string, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		Config:            aws.Config{},
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		return "", err
	}

	ssmsvc := ssm.New(sess, aws.NewConfig())
	param, err := ssmsvc.GetParameter(&p)
	if err != nil {
		return "", err
	}

	return *param.Parameter.Value, nil
}

func init() {
	const parameterEnvVar = "SLACK_WEBHOOK_TOKEN_PARAM_PATH"
	parameter, exists := os.LookupEnv(parameterEnvVar)
	if !exists {
		log.Fatalf("failed to find the env var: %s", parameterEnvVar)
	}

	var err error
	slackToken, err = GetParameter(ssm.GetParameterInput{
		Name:           aws.String(parameter),
		WithDecryption: aws.Bool(false),
	})
	if err != nil {
		log.Fatalf("failed to fetch slack webhook token: %s", err)
	}
}

func LambdaHandler(ctx context.Context, event events.CloudWatchEvent) error {

	return nil
}

func main() {
	lambda.Start(LambdaHandler)
}
