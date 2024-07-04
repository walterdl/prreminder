package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
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

	NewAppStack(app, "PRReminderStack", &AppStackProps{
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

	newMessageQueue := awssqs.NewQueue(stack, jsii.String("PRReminder-NewSlackMessage"), &awssqs.QueueProps{
		QueueName:         jsii.String("PRReminder-NewSlackMessage"),
		VisibilityTimeout: awscdk.Duration_Seconds(jsii.Number(20)),
	})
	newSlackWebhook(stack, &slackWebhookProps{
		newMessageQueue,
	})
	notifier := NewNotifier(stack, notifierProps{
		prChecker:          newPRChecker(stack),
		waitTimeCalc:       newWaitTimeCalc(stack),
		notificationSender: newNotificationSender(stack),
	})
	newNotifierStarter(stack, &notifierStarterProps{
		newMessageQueue,
		notifier.stateMachine,
	})

	return stack
}
