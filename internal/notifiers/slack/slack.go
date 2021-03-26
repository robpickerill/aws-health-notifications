package slack

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/rpickerill/aws-health-to-slack/internal/health"
)

type SlackClient struct {
	WebHookUrl string
	UserName   string
	Client     *http.Client
}

type SlackBlockMessage struct {
	Blocks []SlackBlock `json:"blocks"`
}
type SlackBlock struct {
	Type string          `json:"type"`
	Text *SlackBlockText `json:"text,omitempty"`
}

type SlackBlockText struct {
	Type  string `json:"type,omitempty"`
	Text  string `json:"text,omitempty"`
	Emoji bool   `json:"emoji,omitempty"`
}

func NewSlackClient(ctx context.Context, webhook string, username string, timeout time.Duration) *SlackClient {
	httpClient := http.Client{
		Timeout: timeout,
	}

	return &SlackClient{
		WebHookUrl: webhook,
		UserName:   username,
		Client:     &httpClient,
	}
}

func (s *SlackClient) Notify(ctx context.Context, wg *sync.WaitGroup, errorChan chan<- error, event health.HealthEvent) {
	defer wg.Done()

	message, err := s.parseMessage(event)
	if err != nil {
		errorChan <- err
		return
	}

	err = s.writeHTTPRequest(message)
	if err != nil {
		errorChan <- err
		return
	}
}

func (s *SlackClient) writeHTTPRequest(message SlackBlockMessage) error {
	body, _ := json.Marshal(message)
	log.Printf("%+v", message)

	req, err := http.NewRequest(http.MethodPost, s.WebHookUrl, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := s.Client.Do(req)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return err
	}
	if buf.String() != "ok" {
		return fmt.Errorf("%s HTTP response returned from Slack with body: %s", resp.Status, buf.String())
	}

	return nil
}

func (s *SlackClient) parseMessage(event health.HealthEvent) SlackBlockMessage {
	issueCode := humanReadableString(event.Detail.EventTypeCode, "title")
	issueType := humanReadableString(event.Detail.EventTypeCategory, "title")

	block := []SlackBlock{
		{
			Type: "header",
			Text: &SlackBlockText{
				Text: fmt.Sprintf("%s %s | %s", event.DetailType, issueType, issueCode),
				Type: "plain_text",
			},
		}, {
			Type: "section",
			Text: &SlackBlockText{
				Text: fmt.Sprintf("AWS Account ID: `%s` in AWS Region: `%s`", event.AccountID, event.Region),
				Type: "mrkdwn",
			},
		}, {
			Type: "divider",
		},
	}

	for _, value := range event.Detail.EventDescription {
		block = append(block, SlackBlock{
			Type: "section",
			Text: &SlackBlockText{
				Text:  fmt.Sprintf("Lang: %s | Description: %s", value.Language, value.LatestDescription),
				Type:  "mrkdwn",
				Emoji: false,
			},
		})
	}

	return SlackBlockMessage{
		Blocks: block,
	}
}

func humanReadableString(s string, c string) string {
	s = strings.Replace(s, "_", " ", -1)

	switch c {
	case "upper":
		return strings.ToUpper(s)
	case "lower":
		return strings.ToLower(s)
	case "title":
		return strings.Title(strings.ToLower(s))
	default:
		return s
	}

}
