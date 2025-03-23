package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ShareFrame/posting-service/atproto"
	"github.com/ShareFrame/posting-service/handler"
	"github.com/ShareFrame/posting-service/models"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sirupsen/logrus"
)

type CreatePostInput struct {
	AuthToken string   `json:"authToken"`
	DID       string   `json:"did"`
	Text      string   `json:"text,omitempty"`
	ImageUris []string `json:"imageUris,omitempty"`
	VideoUris []string `json:"videoUris,omitempty"`
}

var client = atproto.NewATProtoService(http.DefaultClient)

type LambdaUnitPayload struct {
	Body string `json:"body"`
}

func handlerFunc(ctx context.Context, event LambdaUnitPayload) (models.PostResponse, error) {
	var input CreatePostInput
	if err := json.Unmarshal([]byte(event.Body), &input); err != nil {
		logrus.WithError(err).Error("Failed to parse request body")
		return models.PostResponse{}, fmt.Errorf("invalid input")
	}

	post := models.ShareFrameFeedPost{
		NSID:      "social.shareframe.feed.post",
		Text:      input.Text,
		ImageUris: input.ImageUris,
		VideoUris: input.VideoUris,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
		SourceApp: "ShareFrame",
	}

	payload := models.RequestPayload{
		AuthToken: input.AuthToken,
		DID:       input.DID,
		Post:      post,
	}

	resp, err := handler.PostHandler(ctx, client, payload)
	if err != nil {
		logrus.WithError(err).Error("PostHandler failed")
		return models.PostResponse{}, err
	}

	return *resp, nil
}

func main() {
	lambda.Start(handlerFunc)
}
