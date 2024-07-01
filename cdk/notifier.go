package main

import (
	"os"
	"strconv"

	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	sfn "github.com/aws/aws-cdk-go/awscdk/v2/awsstepfunctions"
	sfnTasks "github.com/aws/aws-cdk-go/awscdk/v2/awsstepfunctionstasks"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type notifierProps struct {
	prChecker          awslambda.IFunction
	waitTimeCalc       awslambda.IFunction
	notificationSender awslambda.IFunction
}
type Notifier struct {
	stateMachine sfn.StateMachine
}

func NewNotifier(scope constructs.Construct, props notifierProps) *Notifier {
	waitTimeCalcStep := sfnTasks.NewLambdaInvoke(scope, jsii.String("WaitTimeCalcTask"), &sfnTasks.LambdaInvokeProps{
		LambdaFunction: props.waitTimeCalc,
		InputPath:      jsii.String("$"),
		OutputPath:     jsii.String("$.Payload"),
	})
	prCheckerStep := sfnTasks.NewLambdaInvoke(scope, jsii.String("PRCheckerTask"), &sfnTasks.LambdaInvokeProps{
		LambdaFunction: props.prChecker,
		InputPath:      jsii.String("$"),
		OutputPath:     jsii.String("$.Payload"),
	})
	notificationSenderStep := sfnTasks.NewLambdaInvoke(scope, jsii.String("NotificationSenderTask"), &sfnTasks.LambdaInvokeProps{
		LambdaFunction: props.notificationSender,
		InputPath:      jsii.String("$"),
		OutputPath:     jsii.String("$.Payload"),
	})
	waitStep := sfn.NewWait(scope, jsii.String("WaitTask"), &sfn.WaitProps{
		Time:      sfn.WaitTime_SecondsPath(jsii.String("$.waitingTimeInSeconds")),
		StateName: jsii.String("WaitForApproval"),
	})
	endStep := sfn.NewSucceed(scope, jsii.String("EndState"), &sfn.SucceedProps{})

	definition := waitTimeCalcStep.Next(
		waitStep.Next(
			prCheckerStep.Next(
				sfn.NewChoice(scope, jsii.String("IsApproved"), &sfn.ChoiceProps{
					Comment:    jsii.String("Check if the PR is approved. If it is, the state machine ends. Otherwise, continues to send a reminder."),
					StateName:  jsii.String("IsApproved"),
					InputPath:  jsii.String("$"),
					OutputPath: jsii.String("$"),
				}).When(
					sfn.Condition_BooleanEquals(
						jsii.String("$.approvalStatus.approved"),
						jsii.Bool(true),
					),
					endStep,
					nil,
				).Otherwise(
					// Loop back to the start of the state machine.
					notificationSenderStep.Next(
						sfn.NewChoice(scope, jsii.String("MaxNotifications"), &sfn.ChoiceProps{
							Comment:    jsii.String("Check if the maximum number of notifications has been sent."),
							StateName:  jsii.String("MaxNotifications"),
							InputPath:  jsii.String("$"),
							OutputPath: jsii.String("$"),
						}).When(
							sfn.Condition_NumberLessThan(
								jsii.String("$.executionsCount"),
								jsii.Number(maxNotifications()),
							),
							waitTimeCalcStep,
							nil,
						).Otherwise(endStep),
					),
				),
			),
		),
	)

	stateMachine := sfn.NewStateMachine(scope, jsii.String("StateMachine"), &sfn.StateMachineProps{
		DefinitionBody:   sfn.DefinitionBody_FromChainable(definition),
		StateMachineType: sfn.StateMachineType_STANDARD,
		StateMachineName: jsii.String("PRReminder-SFNMachine"),
	})

	return &Notifier{
		stateMachine,
	}
}

func maxNotifications() int {
	defaultVal := 2
	rawVal := os.Getenv("MAX_NOTIFICATIONS")
	if rawVal == "" {
		return defaultVal
	}

	result, err := strconv.Atoi(rawVal)
	if err != nil {
		return defaultVal
	}

	return result
}
