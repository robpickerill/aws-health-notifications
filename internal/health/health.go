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
