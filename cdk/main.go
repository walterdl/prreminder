package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigatewayv2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigatewayv2integrations"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awssqs"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type AppStackProps struct {
	awscdk.StackProps
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewAppStack(app, "PRRemindStack", &AppStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

func env() *awscdk.Environment {
	return nil
}

func NewAppStack(scope constructs.Construct, id string, props *AppStackProps) awscdk.Stack {
	var sprops awscdk.StackProps

	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	httpApi := awsapigatewayv2.NewHttpApi(stack, jsii.String("PRRemindHttpApi"), &awsapigatewayv2.HttpApiProps{
		ApiName: jsii.String("PRRemindHttpApi"),
	})

	newSlackMessageQueue := awssqs.NewQueue(stack, jsii.String("PRRemind-NewSlackMessage"), &awssqs.QueueProps{
		QueueName:         jsii.String("PRRemind-NewSlackMessage"),
		VisibilityTimeout: awscdk.Duration_Seconds(jsii.Number(20)),
	})
	slackWebhookFn := awslambda.NewFunction(stack, jsii.String("SlackWebhook"), &awslambda.FunctionProps{
		FunctionName: jsii.String("PRReminder-SlackWebhook"),
		Code:         awslambda.Code_FromAsset(jsii.String("../slackwebhook/dist"), nil),
		Runtime:      awslambda.Runtime_PROVIDED_AL2(),
		Handler:      jsii.String("bootstrap"),
		Architecture: awslambda.Architecture_ARM_64(),
		Environment: &map[string]*string{
			"NEW_MESSAGE_QUEUE_URL": newSlackMessageQueue.QueueUrl(),
		},
	})
	newSlackMessageQueue.GrantSendMessages(slackWebhookFn)

	slackWebhookIntegration := awsapigatewayv2integrations.NewHttpLambdaIntegration(jsii.String("slackWebhookHTTPIntegration"), slackWebhookFn, &awsapigatewayv2integrations.HttpLambdaIntegrationProps{})
	httpApi.AddRoutes(&awsapigatewayv2.AddRoutesOptions{
		Path: jsii.String("/slack"),
		Methods: &[]awsapigatewayv2.HttpMethod{
			"POST",
		},
		Integration: slackWebhookIntegration,
	})

	NewReminderStarter(stack, &reminderStarterProps{
		newSlackMessageQueue: &newSlackMessageQueue,
	})

	return stack
}

type reminderStarterProps struct {
	newSlackMessageQueue *awssqs.Queue
}

func NewReminderStarter(scope constructs.Construct, props *reminderStarterProps) {
	lambdaFn := awslambda.NewFunction(scope, jsii.String("ReminderStarter"), &awslambda.FunctionProps{
		FunctionName: jsii.String("PRReminder-ReminderStarter"),
		Code:         awslambda.Code_FromAsset(jsii.String("../reminderstarter/dist"), nil),
		Runtime:      awslambda.Runtime_PROVIDED_AL2(),
		Handler:      jsii.String("bootstrap"),
		Architecture: awslambda.Architecture_ARM_64(),
	})
	(*props.newSlackMessageQueue).GrantConsumeMessages(lambdaFn)
}
