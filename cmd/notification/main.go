package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/walterdl/prremind/lib/notifiertypes"
)

func LambdaHandler(_ context.Context, input notifiertypes.NotifierPayload) (notifiertypes.NotifierPayload, error) {
	err := sendNotification(input)
	if err != nil {
		panic(err)
	}

	return input, nil
}

func main() {
	lambda.Start(LambdaHandler)
}
