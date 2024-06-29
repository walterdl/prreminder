package main

import (
	"errors"
	"log"

	"github.com/walterdl/prremind/lib/slack"
)

func processSlackMessage(msg slack.BaseSlackMessageEvent) error {
	if slack.IsRootMessageEdition(msg) {
		err := cancelCurrentReminder(msg)
		if err != nil && !errors.Is(err, errReminderNotFound) {
			return err
		}
	}

	prs := prLinks(msg)
	if prs == nil {
		return nil
	}

	err := startReminder(prs, msg.Event)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
