package slack

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/rpickerill/aws-health-to-slack/internal/health"
)

type SlackClient struct {
	WebHookUrl string
	UserName   string
	Channel    string
	Client     *http.Client
}

type SlackMessage struct {
	Username  string `json:"username,omitempty"`
	IconEmoji string `json:"icon_emoji,omitempty"`
	Channel   string `json:"channel,omitempty"`
	Text      string `json:"text,omitempty"`
}

func NewSlackClient(ctx context.Context, webhook string, username string, timeout time.Duration) *SlackClient {
	httpClient := http.Client{
		Timeout: timeout,
	}

	slackClient := SlackClient{
		WebHookUrl: webhook,
		UserName:   username,
		Client:     &httpClient,
	}

	channel, exists := os.LookupEnv("SLACK_CHANNEL")
	if exists {
		slackClient.Channel = channel
	}

	return &slackClient
}

func (s *SlackClient) Notify(ctx context.Context, wg *sync.WaitGroup, event health.HealthEvent) error {
	defer wg.Done()

	sev := health.GetSeverity(event)

	var message, emoji string
	if sev == health.URGENT {
		emoji = ":red_circle"
	} else {
		emoji = ":hammer_and_wrench:"
	}

	slackRequest := SlackMessage{
		Text:      message,
		Username:  s.UserName,
		IconEmoji: emoji,
		Channel:   s.Channel,
	}

	s.writeHTTPRequest(slackRequest)

	return nil
}

func (s *SlackClient) writeHTTPRequest(message SlackMessage) error {
	body, _ := json.Marshal(message)

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
		return errors.New("non-ok response returned from Slack")
	}

	return nil
}
