package main

import (
	"context"
	"net/http"

	"github.com/ShareFrame/posting-service/atproto"
	"github.com/ShareFrame/posting-service/handler"
	"github.com/ShareFrame/posting-service/models"
	"github.com/aws/aws-lambda-go/lambda"
)

func lambdaHandler(ctx context.Context, request models.RequestPayload) (*models.PostResponse, error) {
	client := atproto.NewATProtoService(&http.Client{})
	return handler.PostHandler(ctx, client, request)
}

func main() {
	lambda.Start(lambdaHandler)
}
