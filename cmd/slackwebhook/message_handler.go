package main

import (
	"encoding/json"
)

func handleMessageEvent(rawEvent string) (string, error) {
	var msgEvent MessageEvent
	err := json.Unmarshal([]byte(rawEvent), &msgEvent)
	if err != nil {
		return "", err
	}

	if !isPublishable(msgEvent) {
		return unknownEvent, nil
	}

	err = publishToSQS(msgEvent.Event)
	if err != nil {
		return "", err
	}

	return "Message received", nil
}

func isPublishable(msg MessageEvent) bool {
	if msg.Type != "event_callback" || msg.Event.Type != "message" {
		// Not a message event
		return false
	}

	if msg.Event.ThreadTS == "" && msg.Event.Subtype == "" && msg.Event.Text != "" {
		// It's a root message
		return true
	}

	if msg.Event.Subtype == "message_changed" && msg.Event.Message.ThreadTS == "" && msg.Event.Message.Text != "" {
		// It's a root message edition
		return true
	}

	if msg.Event.Subtype == "message_changed" && msg.Event.Message.Ts == msg.Event.Message.ThreadTS && msg.Event.Message.Text != "" {
		// It's a root message edition which has replies.
		return true
	}

	return false
}
