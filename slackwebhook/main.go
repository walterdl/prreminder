package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
)

func init() {
	log.Println("Hello 2!")
}

func LambdaHandler(ctx context.Context) (int, error) {
	return 7, nil
}

func main() {
	lambda.Start(LambdaHandler)
}
