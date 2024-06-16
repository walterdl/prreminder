package main

import (
	"encoding/json"
)

type MessageSlackEvent struct {
	Event SlackMessage `json:"event"`
}

type SlackMessage struct {
	Type    string `json:"type"`
	Text    string `json:"text"`
	Ts      string `json:"ts"`
	Channel string `json:"channel"`
}

func handleMessageEvent(rawEvent string) (string, error) {
	var authEvent MessageSlackEvent
	err := json.Unmarshal([]byte(rawEvent), &authEvent)
	if err != nil {
		return "", err
	}

	if authEvent.Event.Type != "message" {
		return unknownEvent, nil
	}

	err = publishToSQS(authEvent.Event)
	if err != nil {
		return "", err
	}

	return "Message received", nil
}
