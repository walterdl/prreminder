package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsssm"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

func newPRChecker(scope constructs.Construct) awslambda.Function {
	gitlabAPIKey := awsssm.StringParameter_FromStringParameterAttributes(scope, jsii.String("GitlabAPIKey"), &awsssm.StringParameterAttributes{
		ParameterName: jsii.String("/prreminder/gitlab/api-key"),
		ValueType:     awsssm.ParameterValueType_STRING,
	})

	return awslambda.NewFunction(scope, jsii.String("PRChecker"), &awslambda.FunctionProps{
		FunctionName: jsii.String("PRReminder-PRChecker"),
		Code:         awslambda.Code_FromAsset(cmdPath("prchecker"), nil),
		Runtime:      awslambda.Runtime_PROVIDED_AL2(),
		Handler:      jsii.String("bootstrap"),
		Architecture: awslambda.Architecture_ARM_64(),
		Environment: &map[string]*string{
			"GITLAB_API_KEY": gitlabAPIKey.StringValue(),
		},
	})
}
