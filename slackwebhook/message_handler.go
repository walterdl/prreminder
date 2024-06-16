package main

import (
	"encoding/json"
)

type MessageSlackEvent struct {
	Event struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"event"`
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

	return "Message received", nil
}
