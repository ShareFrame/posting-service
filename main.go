package main

import (
	"github.com/ShareFrame/posting-service/handler"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handler.PostHandler)
}