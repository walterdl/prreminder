package notifiertypes

import "github.com/walterdl/prremind/lib/slack"

type PRLink struct {
	URL       string `json:"url"`
	Namespace string `json:"namespace"`
	Project   string `json:"project"`
	PRID      string `json:"prID"`
}

type PRApprovalStatus struct {
	Approved          bool `json:"approved"`
	ApprovalsRequired int  `json:"approvalsRequired"`
	ApprovalsLeft     int  `json:"approvalsLeft"`
}

// NotifierPayload is the data used across the entire state machine.
type NotifierPayload struct {
	PR                PRLink                      `json:"pr"`
	PRApprovalStatus  PRApprovalStatus            `json:"approvalStatus"`
	PRNotFound        bool                        `json:"prNotFound"`
	Msg               slack.BaseSlackMessageEvent `json:"slackMessage"`
	WaitingTimeInSecs int                         `json:"waitingTimeInSeconds"`
	ExecCount         int                         `json:"executionsCount"`
}
