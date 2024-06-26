package main

import (
	"log"

	"github.com/walterdl/prremind/lib/slack"
)

func processSlackMessage(msg slack.SlackMessageEvent) error {
	prs := prLinks(msg)
	if prs == nil {
		return nil
	}

	err := startReminder(prs, msg)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
