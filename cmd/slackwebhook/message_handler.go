package main

import (
	"encoding/json"

	"github.com/walterdl/prremind/lib/slack"
)

func handleMessageEvent(rawEvent string) (string, error) {
	var msg slack.BaseSlackMessageEvent
	err := json.Unmarshal([]byte(rawEvent), &msg)
	if err != nil {
		return "", err
	}

	if !slack.IsRootMessage(msg) && !slack.IsRootMessageEdition(msg) {
		return unknownEvent, nil
	}

	err = queueMsg(msg)
	if err != nil {
		return "", err
	}

	return "Message received", nil
}
