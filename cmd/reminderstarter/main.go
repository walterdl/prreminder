package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/walterdl/prremind/lib/notifiertypes"
)

func LambdaHandler(ctx context.Context, sqsEvent events.SQSEvent) error {
	var slackMessage notifiertypes.SlackMessage
	// The process receives one SQS message at a time. Thus, it can safely retrieve just the first element.
	err := json.Unmarshal([]byte(sqsEvent.Records[0].Body), &slackMessage)
	if err != nil {
		return fmt.Errorf("invalid sqs message body: %v", err)
	}

	processSlackMessage(slackMessage)

	return nil
}

func main() {
	lambda.Start(LambdaHandler)
}
