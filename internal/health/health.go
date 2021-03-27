package health

import (
	"encoding/json"
	"strings"
	"time"
	"unicode"

	"golang.org/x/text/language"
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
	Resources  []string          `json:"resources,omitempty"`
	Detail     HealthEventDetail `json:"detail"`
}

type HealthEventDetail struct {
	EventArn          string                        `json:"eventArn"`
	Service           string                        `json:"service"`
	EventTypeCode     string                        `json:"eventTypeCode"`
	EventTypeCategory string                        `json:"eventTypeCategory"`
	StartTime         time.Time                     `json:"startTime,omitempty"`
	EndTime           time.Time                     `json:"endTime,omitempty"`
	EventDescription  []HealthEventDescription      `json:"eventDescription"`
	AffectedEntities  []HealthEventAffectedEntities `json:"affectedEntities,omitempty"`
}

type HealthEventDescription struct {
	Language          language.Tag `json:"language"`
	LatestDescription string       `json:"latestDescription"`
}

type HealthEventAffectedEntities struct {
	EntityValue string          `json:"entityValue"`
	Tags        json.RawMessage `json:"tags"`
}

func GetSeverity(event HealthEvent) Severity {
	switch event.Detail.EventTypeCategory {
	case "issue":
		return URGENT
	default:
		return INFORMATION
	}
}

// -----
// Helper functions for health events

func DeCamelCase(s string) string {
	var b strings.Builder

	priorLowOrNum := false
	for _, v := range s {
		if priorLowOrNum && unicode.IsUpper(v) {
			b.WriteString(" ")
		}
		b.WriteRune(v)
		priorLowOrNum = unicode.IsLower(v) || unicode.IsNumber(v)
	}

	return b.String()
}

func ToTitle(s string) string {
	// using a map for lookup efficiency
	capitalisedAcroynms := map[string]bool{
		"aws": true,
		"ec2": true,
		"rds": true,
		"api": true,
	}

	lowercaseWords := map[string]bool{
		"a":   true,
		"an":  true,
		"on":  true,
		"the": true,
		"to":  true,
	}

	words := strings.Fields(s)

	for i, v := range words {
		lower := strings.ToLower(v)

		if capitalisedAcroynms[lower] {
			words[i] = strings.ToUpper(lower)
			continue
		}

		if lowercaseWords[lower] && i != 0 {
			words[i] = lower
			continue
		}

		words[i] = strings.Title(lower)
	}

	return strings.Join(words, " ")
}
