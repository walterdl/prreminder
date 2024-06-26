package main

import (
	"encoding/json"

	"github.com/walterdl/prremind/lib/slack"
)

func handleAuthEvent(rawEvent string) (string, error) {
	var authEvent slack.AuthSlackEvent
	err := json.Unmarshal([]byte(rawEvent), &authEvent)
	if err != nil {
		return "", err
	}

	return authEvent.Challenge, nil
}
