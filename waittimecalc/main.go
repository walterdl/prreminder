package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/walterdl/prremind/notifiertypes"
)

func LambdaHandler(ctx context.Context, input notifiertypes.NotifierPayload) (notifiertypes.NotifierPayload, error) {
	waitingTime, err := calcWaitingTime(input)
	if err != nil {
		return input, err
	}
	input.WaitingTime = int(waitingTime.Minutes())
	return input, nil
}

func main() {
	lambda.Start(LambdaHandler)
}
