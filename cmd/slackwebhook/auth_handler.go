package main

import (
	"encoding/json"
)

type AuthSlackEvent struct {
	Type      string `json:"type"`
	Token     string `json:"token"`
	Challenge string `json:"challenge"`
}

func handleAuthEvent(rawEvent string) (string, error) {
	var authEvent AuthSlackEvent
	err := json.Unmarshal([]byte(rawEvent), &authEvent)
	if err != nil {
		return "", err
	}

	return authEvent.Challenge, nil
}
