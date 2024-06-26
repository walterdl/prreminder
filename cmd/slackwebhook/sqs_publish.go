package main

import (
	"context"
	"encoding/json"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/go-errors/errors"
	"github.com/walterdl/prremind/lib/slack"
)

func publishToSQS(msg slack.SlackMessageEvent) error {
	client, err := sqsClient()
	if err != nil {
		return errors.New(err)
	}

	sqsBody, err := marshalMsg(msg)
	if err != nil {
		return errors.New(err)
	}
	queueURL := os.Getenv("NEW_MESSAGE_QUEUE_URL")

	_, err = client.SendMessage(context.TODO(), &sqs.SendMessageInput{
		MessageBody: sqsBody,
		QueueUrl:    &queueURL,
	})

	if err != nil {
		return errors.New(err)
	}

	return nil
}

func marshalMsg(msg slack.SlackMessageEvent) (*string, error) {
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		return nil, errors.New(err)
	}

	result := string(jsonMsg)
	return &result, nil
}

var clientCache *sqs.Client

func sqsClient() (*sqs.Client, error) {
	if clientCache != nil {
		return clientCache, nil
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	clientCache = sqs.NewFromConfig(cfg)
	return clientCache, nil
}
