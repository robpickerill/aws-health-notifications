package slack

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"
)

type SlackClient struct {
	WebHookUrl string
	UserName   string
	Channel    string
	TimeOut    time.Duration
}

type SlackMessage struct {
	Username    string       `json:"username,omitempty"`
	IconEmoji   string       `json:"icon_emoji,omitempty"`
	Channel     string       `json:"channel,omitempty"`
	Text        string       `json:"text,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
}

type Attachment struct {
	Color         string      `json:"color,omitempty"`
	Fallback      string      `json:"fallback,omitempty"`
	CallbackID    string      `json:"callback_id,omitempty"`
	ID            int         `json:"id,omitempty"`
	AuthorID      string      `json:"author_id,omitempty"`
	AuthorName    string      `json:"author_name,omitempty"`
	AuthorSubname string      `json:"author_subname,omitempty"`
	AuthorLink    string      `json:"author_link,omitempty"`
	AuthorIcon    string      `json:"author_icon,omitempty"`
	Title         string      `json:"title,omitempty"`
	TitleLink     string      `json:"title_link,omitempty"`
	Pretext       string      `json:"pretext,omitempty"`
	Text          string      `json:"text,omitempty"`
	ImageURL      string      `json:"image_url,omitempty"`
	ThumbURL      string      `json:"thumb_url,omitempty"`
	MarkdownIn    []string    `json:"mrkdwn_in,omitempty"`
	Ts            json.Number `json:"ts,omitempty"`
}

func New() (*SlackClient, error) {
	webhook, exists := os.LookupEnv("SLACK_WEBHOOK")
	if !exists {
		return &SlackClient{}, errors.New("unable to find env var: SLACK_WEBHOOK")
	}

	username, exists := os.LookupEnv("SLACK_USERNAME")
	if !exists {
		username = "AWS Health Notifications"
	}

	client := SlackClient{
		WebHookUrl: webhook,
		UserName:   username,
		TimeOut:    time.Duration(5 * time.Second),
	}

	channel, exists := os.LookupEnv("SLACK_CHANNEL")
	if exists {
		client.Channel = channel
	}

	return &client, nil
}

func (s *SlackClient) Notify(event health.HealthEventDetail) error {
	slackRequest := SlackMessage{
		Text:      message,
		Username:  s.UserName,
		IconEmoji: ":hammer and wrench",
		Channel:   s.Channel,
	}

	err := s.sendHTTPRequest(slackRequest)
	if err != nil {
		return err
	}

	return nil
}

func (s *SlackClient) sendHTTPRequest(slackRequest SlackMessage) error {
	slackBody, _ := json.Marshal(slackRequest)
	req, err := http.NewRequest(http.MethodPost, s.WebHookUrl, bytes.NewBuffer(slackBody))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: s.TimeOut}
	resp, err := client.Do(req)
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
