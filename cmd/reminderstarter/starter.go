package main

import (
	"context"
	"encoding/json"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/sfn"
	"github.com/walterdl/prremind/lib/notifiertypes"
	sfnClient "github.com/walterdl/prremind/lib/sfn"
	"github.com/walterdl/prremind/lib/slack"
)

func startReminder(prs []notifiertypes.PRLink, msg slack.SlackMessageEvent) error {
	client, err := sfnClient.New()
	if err != nil {
		return err
	}

	arn := os.Getenv("STATE_MACHINE_ARN")
	input, err := stateMachineInput(prs, msg)
	if err != nil {
		return err
	}

	_, err = client.StartExecution(context.TODO(), &sfn.StartExecutionInput{
		StateMachineArn: &arn,
		Input:           input,
		Name:            reminderName(reminderNameInput{msg: msg, onlyPrefix: false}),
	})
	if err != nil {
		return err
	}

	return nil
}

func stateMachineInput(prs []notifiertypes.PRLink, msg slack.SlackMessageEvent) (*string, error) {
	input := notifiertypes.NotifierPayload{
		PRs: prs,
		Msg: msg,
	}
	jsonInput, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}
	result := string(jsonInput)

	return &result, nil
}
