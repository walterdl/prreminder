package main

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/walterdl/prremind/lib/slack"
)

type reminderNameInput struct {
	msg slack.BaseSlackMessageEvent
	// onlyPrefix is true when the predictable part of the name is needed.
	onlyPrefix bool
}

// reminderName generates a unique name for the state machine execution
// The prefix is {channelName}-{messageTs}.
// The suffix is a random UUID.
// State machine executions for a message can be identified by the prefix.
func reminderName(input reminderNameInput) (*string, error) {
	prefix := fmt.Sprintf("%s-%s", input.msg.Event.Channel, input.msg.Event.Ts)
	if input.onlyPrefix {
		return &prefix, nil
	}

	uuidVal, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	suffix := uuidVal.String()
	result := fmt.Sprintf("%s--%s", prefix, suffix)

	return &result, nil
}
