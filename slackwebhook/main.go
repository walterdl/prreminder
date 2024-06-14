package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type SlackEvent struct {
	Type string `json:"type"`
}

func main() {
	lambda.Start(LambdaHandler)
}

func LambdaHandler(ctx context.Context, ev events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	var slackEvent SlackEvent
	err := json.Unmarshal([]byte(ev.Body), &slackEvent)
	if err != nil {
		log.Println("Error unmarshalling Slack event:", err)
		return errorResponse(err)
	}

	if slackEvent.Type == "url_verification" {
		var authEvent AuthSlackEvent
		err := json.Unmarshal([]byte(ev.Body), &authEvent)
		if err != nil {
			log.Println("Error unmarshalling Slack auth event:", err)
			return errorResponse(err)
		}

		return handleAuthRequest(authEvent)
	}

	return errorResponse(errors.New("unknown Slack event type"))
}

type AuthSlackEvent struct {
	Type      string
	Token     string
	Challenge string
}

func handleAuthRequest(ev AuthSlackEvent) (events.APIGatewayV2HTTPResponse, error) {
	body, _ := json.Marshal(map[string]string{
		"challenge": ev.Challenge,
	})
	return events.APIGatewayV2HTTPResponse{
		StatusCode: 200,
		Body:       string(body),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func errorResponse(err error) (events.APIGatewayV2HTTPResponse, error) {
	return events.APIGatewayV2HTTPResponse{
		StatusCode: 500,
		Body:       err.Error(),
	}, nil
}
