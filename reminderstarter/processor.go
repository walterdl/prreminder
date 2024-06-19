package main

import (
	"log"

	"github.com/walterdl/prremind/notifiertypes"
)

func processSlackMessage(msg notifiertypes.SlackMessage) error {
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
