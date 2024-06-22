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

// NotifierPayload is the data used across the entire state machine.
type NotifierPayload struct {
	PRs         []PRLink     `json:"prs"`
	Msg         SlackMessage `json:"slackMessage"`
	WaitingTime int          `json:"waitingTime"`
}
