package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	sfn "github.com/aws/aws-cdk-go/awscdk/v2/awsstepfunctions"
	sfnTasks "github.com/aws/aws-cdk-go/awscdk/v2/awsstepfunctionstasks"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type notifierProps struct {
	prChecker    awslambda.IFunction
	waitTimeCalc awslambda.IFunction
}
type Notifier struct {
	stateMachine sfn.StateMachine
}

func NewNotifier(scope constructs.Construct, props notifierProps) *Notifier {
	waitTimeCalcStep := sfnTasks.NewLambdaInvoke(scope, jsii.String("WaitTimeCalcTask"), &sfnTasks.LambdaInvokeProps{
		LambdaFunction: props.waitTimeCalc,
		OutputPath:     jsii.String("$"),
		InputPath:      jsii.String("$"),
	})
	prCheckerStep := sfnTasks.NewLambdaInvoke(scope, jsii.String("PRCheckerTask"), &sfnTasks.LambdaInvokeProps{
		LambdaFunction: props.prChecker,
		OutputPath:     jsii.String("$"),
		InputPath:      jsii.String("$"),
	})
	endStep := sfn.NewSucceed(scope, jsii.String("EndState"), &sfn.SucceedProps{})

	definition := waitTimeCalcStep.Next(
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
			).Otherwise(endStep),
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
