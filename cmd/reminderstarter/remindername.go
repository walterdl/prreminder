package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/walterdl/prremind/lib/slack"
)

type reminderNameInput struct {
	msg slack.SlackMessageEvent
	// onlyPrefix is true when the predictable part of the name is needed.
	onlyPrefix bool
}

// reminderName generates a unique name for the state machine execution
// based on the current time in RFC3339 format.
// Provided that the name for a state machine cannot contain colons, this function replaces them with dashes.
func reminderName(input reminderNameInput) *string {
	result := fmt.Sprintf("%s--%s", input.msg.Channel, input.msg.Ts)
	if input.onlyPrefix {
		return &result
	}

	suffix := time.Now().Format(time.DateTime)
	suffix = strings.Replace(suffix, " ", "-", -1)
	suffix = strings.Replace(suffix, ":", "-", -1)
	result = fmt.Sprintf("%s--%s", result, suffix)

	return &result
}
