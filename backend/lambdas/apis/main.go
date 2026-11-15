package main

import (
	"os"

	"apis/cmd"
	"apis/cmd/local"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	if os.Getenv("APP_ENV") == "local" {
		local.Handler()
	} else {
		lambda.Start(cmd.Entry)
	}
}
