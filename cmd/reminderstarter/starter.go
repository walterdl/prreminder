package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/sfn"
	"github.com/walterdl/prremind/lib/notifiertypes"
	sfnClient "github.com/walterdl/prremind/lib/sfn"
	"github.com/walterdl/prremind/lib/slack"
)

func startReminders(prs []notifiertypes.PRLink, msg slack.BaseSlackMessageEvent) error {
	client, err := sfnClient.New()
	if err != nil {
		return err
	}

	arn := os.Getenv("STATE_MACHINE_ARN")
	for _, pr := range prs {
		input, err := stateMachineInput(pr, msg)
		if err != nil {
			return err
		}

		name, err := reminderName(reminderNameInput{msg: msg, onlyPrefix: false})
		fmt.Println("name: ", *name)
		if err != nil {
			return err
		}

		_, err = client.StartExecution(context.TODO(), &sfn.StartExecutionInput{
			StateMachineArn: &arn,
			Input:           input,
			Name:            name,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func stateMachineInput(pr notifiertypes.PRLink, msg slack.BaseSlackMessageEvent) (*string, error) {
	input := notifiertypes.NotifierPayload{
		PR:  pr,
		Msg: msg,
	}
	jsonInput, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}
	result := string(jsonInput)

	return &result, nil
}
