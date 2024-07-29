package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsssm"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

func newNotificationSender(scope constructs.Construct) awslambda.Function {
	slackBotToken := awsssm.StringParameter_FromStringParameterAttributes(scope, jsii.String("SlackBothToken"), &awsssm.StringParameterAttributes{
		ParameterName: jsii.String("/prreminder/slack/bot-token"),
		ValueType:     awsssm.ParameterValueType_STRING,
	})

	return awslambda.NewFunction(scope, jsii.String("NotificationSender"), &awslambda.FunctionProps{
		FunctionName: jsii.String("PRReminder-NotificationSender"),
		Code:         awslambda.Code_FromAsset(cmdPath("notification"), nil),
		Runtime:      awslambda.Runtime_PROVIDED_AL2(),
		Handler:      jsii.String("bootstrap"),
		Architecture: awslambda.Architecture_ARM_64(),
		Environment: &map[string]*string{
			"SLACK_BOT_TOKEN": slackBotToken.StringValue(),
		},
	})
}
