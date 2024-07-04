package main

import (
	"os"
	"strings"

	"github.com/slack-go/slack"
	"github.com/walterdl/prremind/lib/notifiertypes"
	prReminderSlack "github.com/walterdl/prremind/lib/slack"
)

var api = slack.New(os.Getenv("SLACK_BOT_TOKEN"))
var msgTemplate = "Hi!\n\n<!here> This PR is still waiting for approvals. Please review it.\n\n{PR}"

func sendNotification(input notifiertypes.NotifierPayload) error {
	_, _, err := api.PostMessage(input.Msg.Event.Channel,
		slack.MsgOptionText(msg(input), false),
		slack.MsgOptionTS(threadTS(input)),
	)

	if err != nil {
		return err
	}

	return nil
}

func msg(input notifiertypes.NotifierPayload) string {
	return strings.Replace(msgTemplate, "{PR}", input.PR.URL, -1)
}

func threadTS(input notifiertypes.NotifierPayload) string {
	if prReminderSlack.IsRootMessageEdition(input.Msg) {
		return input.Msg.Event.Message.ThreadTS
	}

	return input.Msg.Event.TS
}
