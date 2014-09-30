package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"os"
)

type SlackMsg struct {
	Channel   string `json:"channel"`
	Username  string `json:"username,omitempty"`
	Text      string `json:"text"`
	Parse     string `json:"parse"`
	IconEmoji string `json:"icon_emoji,omitempty"`
}

func (m SlackMsg) Encode() (string, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (m SlackMsg) Post(WebhookURL string) error {
	encoded, err := m.Encode()
	if err != nil {
		return err
	}

	resp, err := http.PostForm(WebhookURL, url.Values{"payload": {encoded}})
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("Not OK")
	}
	return nil
}

func main() {
	if len(os.Args) != 4 {
		log.Fatalln("Args failed")
		return
	}

	token := os.Args[1]
	subject := os.Args[2]
	message := os.Args[3]

	channel := "#devops"
	username := "zabbix"
	webhookUrl := "https://mst365.slack.com/services/hooks/incoming-webhook?token=" + token

	var emoji string
	switch subject {
	case "RECOVERY":
		emoji = ":smile:"
	case "PROBLEM":
		emoji = ":frowning:"
	default:
		emoji = ":ghost:"
	}

	// was there a message on the command line? If so use it.
	msg := SlackMsg{
		Channel:   channel,
		Username:  username,
		Parse:     "full",
		Text:      subject + ": " + message,
		IconEmoji: emoji,
	}

	err := msg.Post(webhookUrl)
	if err != nil {
		log.Fatalf("Post failed: %v", err)
	}
	return
}
