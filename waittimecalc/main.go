package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
)

func LambdaHandler(ctx context.Context, ev interface{}) error {
	log.Println(ev)
	return nil
}

func main() {
	lambda.Start(LambdaHandler)
}
