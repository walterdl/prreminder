package main

import (
	"fmt"

	"github.com/walterdl/prremind/lib/slack"
)

// reminderName generates a unique name for the state machine execution
// based on the current time in RFC3339 format.
// Provided that the name for a state machine cannot contain colons, this function replaces them with dashes.
func reminderName(msg slack.SlackMessageEvent) *string {
	result := fmt.Sprintf("%s-%s", msg.Channel, msg.Ts)

	return &result
}
