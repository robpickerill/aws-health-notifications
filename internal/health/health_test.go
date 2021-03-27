package health

import (
	"fmt"
	"testing"
	"time"
)

func TestSeverity(t *testing.T) {
	timeNow := time.Now()
	timeFuture := timeNow.Add(6 * time.Hour)

	healthEventUrgent := HealthEvent{
		Version:    "0",
		ID:         "121345678-1234-1234-1234-123456789012",
		DetailType: "AWS Health Event",
		Source:     "aws.health",
		AccountID:  "123456789012",
		Time:       timeNow,
		Region:     "us-east-1",
		Detail: HealthEventDetail{
			EventArn:          "arn:aws:health:us-east-1::event/AWS_EC2_INSTANCE_STORE_DRIVE_PERFORMANCE_DEGRADED_90353408594353980",
			Service:           "EC2",
			EventTypeCode:     "AWS_EC2_INSTANCE_STORE_DRIVE_PERFORMANCE_DEGRADED",
			EventTypeCategory: "issue",
			StartTime:         timeNow,
			EndTime:           timeFuture,
		},
	}

	healthEventInformation := HealthEvent{
		Version:    "0",
		ID:         "121345678-1234-1234-1234-123456789012",
		DetailType: "AWS Health Event",
		Source:     "aws.health",
		AccountID:  "123456789012",
		Time:       timeNow,
		Region:     "us-east-1",
		Detail: HealthEventDetail{
			EventArn:          "arn:aws:health:us-east-1::event/AWS_EC2_INSTANCE_STORE_DRIVE_PERFORMANCE_DEGRADED_90353408594353980",
			Service:           "EC2",
			EventTypeCode:     "AWS_EC2_INSTANCE_STORE_DRIVE_PERFORMANCE_DEGRADED",
			EventTypeCategory: "scheduledChange",
			StartTime:         timeNow,
			EndTime:           timeFuture,
		},
	}

	tests := []struct {
		name  string
		input HealthEvent
		want  Severity
	}{
		{name: "test urgent", input: healthEventUrgent, want: URGENT},
		{name: "test information", input: healthEventInformation, want: INFORMATION},
	}

	for _, tt := range tests {
		testName := fmt.Sprintf("%s, %v", tt.name, tt.want)

		t.Run(testName, func(t *testing.T) {
			got := GetSeverity(tt.input)

			if got != tt.want {
				t.Errorf("got: %v, want: %v", got, tt.want)
			}
		})
	}
}

func TestDeCamelCase(t *testing.T) {
	tests := []struct {
		input, want string
	}{
		{"testTest", "test Test"},
		{"awsServiceName", "aws Service Name"},
		{"anAwsServiceName", "an Aws Service Name"},
		{"ec2", "ec2"},
		{"anInstanceOnEc2", "an Instance On Ec2"},
	}

	for _, tt := range tests {
		testName := fmt.Sprintf("%s, %s", tt.input, tt.want)

		t.Run(testName, func(t *testing.T) {
			got := DeCamelCase(tt.input)

			if got != tt.want {
				t.Errorf("got: %s, want: %s", got, tt.want)
			}
		})
	}
}

func TestToTitle(t *testing.T) {
	tests := []struct {
		input, want string
	}{
		{"test", "Test"},
		{"aws service name", "AWS Service Name"},
		{"an aws service name", "An AWS Service Name"},
		{"ec2", "EC2"},
		{"an instance on ec2", "An Instance on EC2"},
	}

	for _, tt := range tests {
		testName := fmt.Sprintf("%s, %s", tt.input, tt.want)

		t.Run(testName, func(t *testing.T) {
			got := ToTitle(tt.input)

			if got != tt.want {
				t.Errorf("got: %s, want: %s", got, tt.want)
			}
		})
	}
}
