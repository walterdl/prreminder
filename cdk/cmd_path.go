package main

import (
	"fmt"

	"github.com/aws/jsii-runtime-go"
)

func cmdPath(cmd string) *string {
	return jsii.String(fmt.Sprintf("../cmd/%s/dist", cmd))
}
