package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigatewayv2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigatewayv2integrations"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awssqs"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type slackWebhookProps struct {
	newMessageQueue awssqs.Queue
}

func newSlackWebhook(scope constructs.Construct, props *slackWebhookProps) {
	httpApi := awsapigatewayv2.NewHttpApi(scope, jsii.String("PRReminderHttpApi"), &awsapigatewayv2.HttpApiProps{
		ApiName: jsii.String("PRReminderHttpApi"),
	})

	slackWebhookFn := awslambda.NewFunction(scope, jsii.String("SlackWebhook"), &awslambda.FunctionProps{
		FunctionName: jsii.String("PRReminder-SlackWebhook"),
		Code:         awslambda.Code_FromAsset(cmdPath("slackwebhook"), nil),
		Runtime:      awslambda.Runtime_PROVIDED_AL2(),
		Handler:      jsii.String("bootstrap"),
		Architecture: awslambda.Architecture_ARM_64(),
		Environment: &map[string]*string{
			"NEW_MESSAGE_QUEUE_URL": props.newMessageQueue.QueueUrl(),
		},
	})
	props.newMessageQueue.GrantSendMessages(slackWebhookFn)

	endpointIntegration := awsapigatewayv2integrations.NewHttpLambdaIntegration(
		jsii.String("slackWebhookHTTPIntegration"),
		slackWebhookFn,
		&awsapigatewayv2integrations.HttpLambdaIntegrationProps{},
	)
	httpApi.AddRoutes(&awsapigatewayv2.AddRoutesOptions{
		Path: jsii.String("/slack"),
		Methods: &[]awsapigatewayv2.HttpMethod{
			"POST",
		},
		Integration: endpointIntegration,
	})
}
