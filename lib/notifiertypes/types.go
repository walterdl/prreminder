package notifiertypes

import "github.com/walterdl/prremind/lib/slack"

type PRLink struct {
	URL       string `json:"url"`
	Namespace string `json:"namespace"`
	Project   string `json:"project"`
	PRID      string `json:"prID"`
}

// NotifierPayload is the data used across the entire state machine.
type NotifierPayload struct {
	PRs         []PRLink                `json:"prs"`
	Msg         slack.SlackMessageEvent `json:"slackMessage"`
	WaitingTime int                     `json:"waitingTime"`
}
