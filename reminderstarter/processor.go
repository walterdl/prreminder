package main

import (
	"fmt"
	"regexp"
	"strings"
)

type SlackMessage struct {
	Type    string `json:"type"`
	Text    string `json:"text"`
	TS      string `json:"ts"`
	Channel string `json:"channel"`
}

type PRLink struct {
	URL       string
	Namespace string
	Project   string
	PRID      string
}

func processSlackMessage(msg SlackMessage) error {
	// Replaces escaped slashes with regular slashes.
	// From https:\/\/gitlab.com\/... to https://gitlab.com/...
	msg.Text = strings.ReplaceAll(msg.Text, `\/`, "/")
	fmt.Println(msg.Text)
	urlSegment := `([a-zA-Z0-9\-_.~]+)`
	prPattern := fmt.Sprintf(`https://gitlab.com/%[1]s/%[1]s/-/merge_requests/%[1]s`, urlSegment)
	re := regexp.MustCompile(prPattern)
	matches := re.FindAllStringSubmatch(msg.Text, -1)
	if matches == nil || len(matches) == 0 {
		// No PR links found. Nothing to do.
		return nil
	}

	prLinks := make([]PRLink, len(matches))
	for i, match := range matches {
		prLinks[i] = PRLink{
			URL:       match[0],
			Namespace: match[1],
			Project:   match[2],
			PRID:      match[3],
		}
	}
	fmt.Println(prLinks)
	return nil
}
