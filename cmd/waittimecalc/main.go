package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/walterdl/prremind/lib/notifiertypes"
)

func LambdaHandler(_ context.Context, input notifiertypes.NotifierPayload) (notifiertypes.NotifierPayload, error) {
	waitingTime, err := calcWaitingTime()
	if err != nil {
		return input, err
	}
	input.WaitingTimeInSecs = int(waitingTime.Seconds())
	return input, nil
}

func main() {
	lambda.Start(LambdaHandler)
}
