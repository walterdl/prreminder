package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/slack-go/slack"
	"github.com/walterdl/prremind/lib/notifiertypes"
	prReminderSlack "github.com/walterdl/prremind/lib/slack"
)

var api = slack.New(os.Getenv("SLACK_BOT_TOKEN"))
var msgTemplate = "Hi!\n\n<!here> This PR is still waiting for {approvalsInfo}. Please review it.\n\n{PR}"

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
	result := strings.ReplaceAll(msgTemplate, "{PR}", input.PR.URL)

	approvalsLeft := fmt.Sprintf(
		"%d approval%s",
		input.PRApprovalStatus.ApprovalsLeft,
		pluralSuffix(input.PRApprovalStatus.ApprovalsLeft),
	)

	return strings.ReplaceAll(result, "{approvalsInfo}", approvalsLeft)
}

func pluralSuffix(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}

func threadTS(input notifiertypes.NotifierPayload) string {
	if prReminderSlack.IsRootMessageEdition(input.Msg) {
		return input.Msg.Event.Message.ThreadTS
	}

	return input.Msg.Event.TS
}
