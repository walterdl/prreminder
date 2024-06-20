package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	sfn "github.com/aws/aws-cdk-go/awscdk/v2/awsstepfunctions"
	sfnTasks "github.com/aws/aws-cdk-go/awscdk/v2/awsstepfunctionstasks"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type Notifier struct {
	stateMachine sfn.StateMachine
}

func NewNotifier(scope constructs.Construct) *Notifier {
	waitTimeCalcFn := awslambda.NewFunction(scope, jsii.String("WaitTimeCalc"), &awslambda.FunctionProps{
		FunctionName: jsii.String("PRReminder-WaitTimeCalc"),
		Code:         awslambda.Code_FromAsset(jsii.String("../waittimecalc/dist"), nil),
		Runtime:      awslambda.Runtime_PROVIDED_AL2(),
		Handler:      jsii.String("bootstrap"),
		Architecture: awslambda.Architecture_ARM_64(),
	})

	waitTimeCalcStep := sfnTasks.NewLambdaInvoke(scope, jsii.String("WaitTimeCalcTask"), &sfnTasks.LambdaInvokeProps{
		LambdaFunction: waitTimeCalcFn,
		OutputPath:     jsii.String("$"),
		InputPath:      jsii.String("$"),
	})

	endStep := sfn.NewSucceed(scope, jsii.String("EndState"), &sfn.SucceedProps{})

	definition := waitTimeCalcStep.Next(endStep)

	stateMachine := sfn.NewStateMachine(scope, jsii.String("StateMachine"), &sfn.StateMachineProps{
		Definition:       definition,
		StateMachineType: sfn.StateMachineType_STANDARD,
		StateMachineName: jsii.String("PRReminder-SFNMachine"),
	})

	return &Notifier{
		stateMachine,
	}
}
