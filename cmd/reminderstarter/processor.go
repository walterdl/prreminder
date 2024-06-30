package main

import (
	"errors"
	"log"

	"github.com/walterdl/prremind/lib/slack"
)

func processSlackMessage(msg slack.BaseSlackMessageEvent) error {
	if slack.IsRootMessageEdition(msg) {
		err := cancelCurrentReminders(msg)
		if err != nil && !errors.Is(err, errRemindersNotFound) {
			return err
		}
	}

	prs := prLinks(msg)
	if prs == nil {
		return nil
	}

	err := startReminders(prs, msg)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
