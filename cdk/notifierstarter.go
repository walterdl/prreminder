package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	eventsources "github.com/aws/aws-cdk-go/awscdk/v2/awslambdaeventsources"
	"github.com/aws/aws-cdk-go/awscdk/v2/awssqs"
	sfn "github.com/aws/aws-cdk-go/awscdk/v2/awsstepfunctions"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type notifierStarterProps struct {
	slackMessagesQueue awssqs.Queue
	stateMachine       sfn.StateMachine
}

func newNotifierStarter(scope constructs.Construct, props *notifierStarterProps) {
	starterFn := awslambda.NewFunction(scope, jsii.String("ReminderStarter"), &awslambda.FunctionProps{
		FunctionName: jsii.String("PRReminder-ReminderStarter"),
		Code:         awslambda.Code_FromAsset(cmdPath("reminderstarter"), nil),
		Runtime:      awslambda.Runtime_PROVIDED_AL2(),
		Handler:      jsii.String("bootstrap"),
		Architecture: awslambda.Architecture_ARM_64(),
	})
	props.slackMessagesQueue.GrantConsumeMessages(starterFn)
	queueEventSource := eventsources.NewSqsEventSource(props.slackMessagesQueue, &eventsources.SqsEventSourceProps{
		BatchSize:         jsii.Number(1),
		Enabled:           jsii.Bool(true),
		MaxBatchingWindow: awscdk.Duration_Minutes(jsii.Number(5)),
	})
	starterFn.AddEventSource(queueEventSource)
	props.stateMachine.GrantStartExecution(starterFn)
	props.stateMachine.GrantRead(starterFn)
	starterFn.AddEnvironment(jsii.String("STATE_MACHINE_ARN"), props.stateMachine.StateMachineArn(), nil)
}
