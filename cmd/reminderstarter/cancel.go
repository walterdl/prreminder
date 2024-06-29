package main

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/sfn"
	"github.com/aws/aws-sdk-go-v2/service/sfn/types"
	sfnClient "github.com/walterdl/prremind/lib/sfn"
	"github.com/walterdl/prremind/lib/slack"
)

var stateMachineArn = os.Getenv("STATE_MACHINE_ARN")
var errReminderNotFound = errors.New("reminder not found")

func cancelCurrentReminder(msg slack.BaseSlackMessageEvent) error {
	client, err := sfnClient.New()
	if err != nil {
		return err
	}

	execution, err := currentExecution(client, msg)
	if err != nil {
		return err
	}
	cause := "User edited original message, potentially including different PRs"
	_, err = client.StopExecution(context.TODO(), &sfn.StopExecutionInput{
		ExecutionArn: execution.ExecutionArn,
		Cause:        &cause,
	})
	return err
}

func currentExecution(client *sfn.Client, msg slack.BaseSlackMessageEvent) (*types.ExecutionListItem, error) {
	// Only try twice to find the reminder execution.
	// This is because the expected use case is that the user edits the message shortly after typing it.
	// Therefore, it's expected that the reminder will be found among the firsts executions.
	maxTries := 2
	var nexToken *string
	for i := 0; i < maxTries; i++ {
		output, err := fetchExecutions(client, nexToken)
		if err != nil {
			return nil, err
		}

		for _, execution := range output.Executions {
			hasName := strings.HasPrefix(
				*execution.Name, *reminderName(reminderNameInput{msg: msg.Event, onlyPrefix: true}),
			)
			if hasName && execution.Status == types.ExecutionStatusRunning {
				return &execution, nil
			}
		}

		if output.NextToken != nil {
			nexToken = output.NextToken
		}
	}

	return nil, errReminderNotFound
}

func fetchExecutions(client *sfn.Client, nextToken *string) (*sfn.ListExecutionsOutput, error) {
	return client.ListExecutions(context.TODO(), &sfn.ListExecutionsInput{
		StateMachineArn: &stateMachineArn,
		MaxResults:      10,
		NextToken:       nextToken,
	})
}
