package main

import (
	"context"
	"errors"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/walterdl/prremind/lib/notifiertypes"
)

func LambdaHandler(ctx context.Context, input notifiertypes.NotifierPayload) (notifiertypes.NotifierPayload, error) {
	approvalStatus, err := checkPR(input)
	input.ExecCount++

	if err != nil {
		if errors.Is(err, errPRNotFound) {
			input.PRNotFound = true
			return input, nil
		}

		panic(err)
	}

	input.PRApprovalStatus = approvalStatus
	return input, nil
}

func main() {
	lambda.Start(LambdaHandler)
}
