package main

import (
	"context"
	"errors"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go-v2/service/sfn"
	"github.com/aws/aws-sdk-go-v2/service/sfn/types"
	sfnClient "github.com/walterdl/prremind/lib/sfn"
	"github.com/walterdl/prremind/lib/slack"
)

var stateMachineArn = os.Getenv("STATE_MACHINE_ARN")
var errRemindersNotFound = errors.New("reminder not found")

func cancelCurrentReminders(msg slack.BaseSlackMessageEvent) error {
	client, err := sfnClient.New()
	if err != nil {
		return err
	}

	remindersPrefix, err := reminderName(reminderNameInput{msg: msg, onlyPrefix: true})
	if err != nil {
		return err
	}

	executions, err := currentExecutions(client, *remindersPrefix)
	if err != nil {
		return err
	}

	stopExecutions(client, executions)
	return nil
}

func currentExecutions(client *sfn.Client, remindersPrefix string) ([]types.ExecutionListItem, error) {
	// Only try twice to find the reminder execution.
	// This is because the expected use case is that the user edits the message shortly after typing it.
	// Therefore, it's expected that the reminders will be found among the firsts executions.
	maxTries := 2
	var nexToken *string
	result := make([]types.ExecutionListItem, 0)

	for i := 0; i < maxTries; i++ {
		output, err := fetchExecutions(client, nexToken)
		if err != nil {
			return nil, err
		}

		for _, execution := range output.Executions {
			hasName := strings.HasPrefix(*execution.Name, remindersPrefix)
			if hasName && execution.Status == types.ExecutionStatusRunning {
				result = append(result, execution)
			}
		}

		if output.NextToken != nil {
			nexToken = output.NextToken
		}
	}

	if len(result) > 0 {
		return result, nil
	}

	return nil, errRemindersNotFound
}

func fetchExecutions(client *sfn.Client, nextToken *string) (*sfn.ListExecutionsOutput, error) {
	return client.ListExecutions(context.TODO(), &sfn.ListExecutionsInput{
		StateMachineArn: &stateMachineArn,
		MaxResults:      10,
		NextToken:       nextToken,
	})
}

func stopExecutions(client *sfn.Client, execs []types.ExecutionListItem) {
	wg := sync.WaitGroup{}

	for _, exec := range execs {
		wg.Add(1)
		go func(client *sfn.Client, exec types.ExecutionListItem) {
			defer wg.Done()
			cause := "User edited original message, potentially including different PRs"
			_, err := client.StopExecution(context.TODO(), &sfn.StopExecutionInput{
				ExecutionArn: exec.ExecutionArn,
				Cause:        &cause,
			})

			if err != nil {
				// Reminders that could not be stopped are left to run.
				// I.e. the error is not propagated but logged.
				// Although not ideal, extra reminders are not harmful.
				log.Println(err)
			}
		}(client, exec)
	}

	wg.Wait()
}
