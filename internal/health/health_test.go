package health

import (
	"fmt"
	"testing"
)

func TestDeCamelCase(t *testing.T) {
	var tests = []struct {
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
	var tests = []struct {
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
