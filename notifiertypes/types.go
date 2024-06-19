package notifiertypes

type SlackMessage struct {
	Type    string `json:"type"`
	Text    string `json:"text"`
	TS      string `json:"ts"`
	Channel string `json:"channel"`
}

type PRLink struct {
	URL       string `json:"url"`
	Namespace string `json:"namespace"`
	Project   string `json:"project"`
	PRID      string `json:"prID"`
}

// NotifierInput is the input for the state machine
type NotifierInput struct {
	PRs []PRLink     `json:"prs"`
	Msg SlackMessage `json:"slackMessage"`
}
