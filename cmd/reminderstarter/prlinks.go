package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/walterdl/prremind/lib/notifiertypes"
	"github.com/walterdl/prremind/lib/slack"
)

func prLinks(msg slack.BaseSlackMessageEvent) []notifiertypes.PRLink {
	rawLinks := abstractFromMessage(msg)
	if rawLinks == nil {
		return nil
	}
	result := make([]notifiertypes.PRLink, len(rawLinks))

	for i, match := range rawLinks {
		result[i] = notifiertypes.PRLink{
			URL:       match[0],
			Namespace: match[1],
			Project:   match[2],
			PRID:      match[3],
		}
	}

	return result
}

// abstractFromMessage extracts PR links from a Slack message.
// It returns a slice of slices of strings. Each sub-slice of strings contains
// the URL, namespace, project, and PR ID of a PR link.
func abstractFromMessage(msg slack.BaseSlackMessageEvent) [][]string {
	// Replaces escaped slashes with regular slashes.
	// From https:\/\/gitlab.com\/... to https://gitlab.com/...
	txt := text(msg)
	txt = strings.ReplaceAll(txt, `\/`, "/")
	urlSegment := `([a-zA-Z0-9\-_.~]+)`
	prPattern := fmt.Sprintf(`https://gitlab.com/%[1]s/%[1]s/-/merge_requests/%[1]s`, urlSegment)
	re := regexp.MustCompile(prPattern)

	return re.FindAllStringSubmatch(txt, -1)
}

func text(msg slack.BaseSlackMessageEvent) string {
	if slack.IsRootMessageEdition(msg) {
		return msg.Event.Message.Text
	}

	return msg.Event.Text
}
