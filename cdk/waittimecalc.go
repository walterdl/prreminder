package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

func newWaitTimeCalc(scope constructs.Construct) awslambda.IFunction {
	return awslambda.NewFunction(scope, jsii.String("WaitTimeCalc"), &awslambda.FunctionProps{
		FunctionName: jsii.String("PRReminder-WaitTimeCalc"),
		Code:         awslambda.Code_FromAsset(cmdPath("waittimecalc"), nil),
		Runtime:      awslambda.Runtime_PROVIDED_AL2(),
		Handler:      jsii.String("bootstrap"),
		Architecture: awslambda.Architecture_ARM_64(),
		Environment: &map[string]*string{
			"TIMEZONE":                 jsii.String("America/New_York"),
			"DAYS":                     jsii.String("1,2,3,4,5"),
			"START_TIME":               jsii.String("8:0"),
			"PR_APPROVAL_WAIT_MINUTES": jsii.String("2880"),
		},
	})
}
