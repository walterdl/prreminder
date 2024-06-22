package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-errors/errors"
)

func main() {
	lambda.Start(LambdaHandler)
}

func LambdaHandler(ctx context.Context, ev events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	res, err := handleSlackEvent(ev.Body)
	if err != nil {
		return errorResponse(err)
	}

	return successResponse(res)
}

func errorResponse(err error) (events.APIGatewayV2HTTPResponse, error) {
	if errWithStack, ok := err.(*errors.Error); ok {
		log.Println(errWithStack.ErrorStack())
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: 500,
		Body:       err.Error(),
	}, nil
}

func successResponse(body string) (events.APIGatewayV2HTTPResponse, error) {
	return events.APIGatewayV2HTTPResponse{
		StatusCode: 200,
		Body:       body,
		Headers: map[string]string{
			"Content-Type": "plain/text",
		},
	}, nil
}
