package health

import (
	"encoding/json"
	"time"
)

type Severity int

const (
	INFORMATION Severity = iota
	URGENT
)

type HealthEvent struct {
	Version    string            `json:"version"`
	ID         string            `json:"id"`
	DetailType string            `json:"detail-type"`
	Source     string            `json:"source"`
	AccountID  string            `json:"account"`
	Time       time.Time         `json:"time"`
	Region     string            `json:"region"`
	Resources  []string          `json:"resources"`
	Detail     HealthEventDetail `json:"detail"`
}

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

func GetSeverity(event HealthEvent) Severity {
	switch event.Detail.EventTypeCategory {
	case "issue":
		return URGENT
	case "accountNotification":
		return INFORMATION
	case "scheduledChange":
		return INFORMATION
	default:
		return INFORMATION
	}
}

func GetEnglishNotification(event HealthEvent) (string, error) {

}
