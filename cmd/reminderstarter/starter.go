package main

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
	"github.com/walterdl/prremind/lib/notifiertypes"
	"github.com/walterdl/prremind/lib/slack"
)

func startReminder(prs []notifiertypes.PRLink, msg slack.SlackMessageEvent) error {
	client, err := sfnClient()
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
		Name:            stateMachineName(),
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

// stateMachineName generates a unique name for the state machine execution
// based on the current time in RFC3339 format.
// Provided that the name for a state machine cannot contain colons, this function replaces them with dashes.
func stateMachineName() *string {
	currentTime := time.Now().UTC()
	isoFormat := currentTime.Format(time.RFC3339)
	result := strings.ReplaceAll(isoFormat, ":", "-")
	return &result
}

var clientCache *sfn.Client

func sfnClient() (*sfn.Client, error) {
	if clientCache != nil {
		return clientCache, nil
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	clientCache = sfn.NewFromConfig(cfg)
	return clientCache, nil
}
