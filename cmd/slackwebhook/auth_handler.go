package main

import (
	"encoding/json"
)

func handleAuthEvent(rawEvent string) (string, error) {
	var authEvent AuthSlackEvent
	err := json.Unmarshal([]byte(rawEvent), &authEvent)
	if err != nil {
		return "", err
	}

	return authEvent.Challenge, nil
}
