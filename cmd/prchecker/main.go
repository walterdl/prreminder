package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/walterdl/prremind/lib/notifiertypes"
)

func LambdaHandler(ctx context.Context, input notifiertypes.NotifierPayload) (notifiertypes.NotifierPayload, error) {
	approvalStatus, err := checkPR(input)
	if err != nil {
		panic(err)
	}

	input.PRApprovalStatus = approvalStatus
	return input, nil
}

func main() {
	lambda.Start(LambdaHandler)
}
