package main

import (
	"encoding/json"

	"github.com/go-errors/errors"
)

type BaseSlackEvent struct {
	Type string `json:"type"`
}

type SlackEventHandler func(rawEvent string) (string, error)

const unknownEvent = "Unknown event"

func handleSlackEvent(rawEvent string) (string, error) {
	var ev BaseSlackEvent
	err := json.Unmarshal([]byte(rawEvent), &ev)
	if err != nil {
		return "", errors.New(err)
	}

	handlersMap := map[string]SlackEventHandler{
		"url_verification": handleAuthEvent,
		"event_callback":   handleMessageEvent,
	}

	if handleEvent, ok := handlersMap[ev.Type]; ok {
		return handleEvent(rawEvent)
	}

	return unknownEvent, nil
}
