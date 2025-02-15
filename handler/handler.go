package handler

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ShareFrame/posting-service/atproto"
	"github.com/ShareFrame/posting-service/models"
	"github.com/sirupsen/logrus"
)

func PostHandler(ctx context.Context, request models.RequestPayload) (map[string]string, error) {
	if request.AuthToken == "" || request.DID == "" {
		err := errors.New("invalid request: missing 'authToken' or 'did'")
		logrus.Error(err)
		return nil, err
	}

	if err := validatePost(request.Post); err != nil {
		logrus.WithError(err).WithField("NSID", request.Post.NSID).Error("Validation failed")
		return nil, fmt.Errorf("invalid post: %w", err)
	}

	response, err := atproto.PostToFeed(request.Post, request.AuthToken, request.DID)
	if err != nil {
		logrus.WithError(err).WithField("DID", request.DID).Error("Failed to post to feed")
		return nil, fmt.Errorf("posting to feed failed: %w", err)
	}

	logrus.WithField("DID", request.DID).Info("Post successfully created")
	return map[string]string{"message": "Post created successfully", "response": response}, nil
}

func validatePost(post models.ShareFrameFeedPost) error {
	switch {
	case post.NSID != "social.shareframe.feed.post":
		return errors.New("invalid NSID: only social.shareframe.feed.post is allowed")
	case len(post.Text) > 300:
		return errors.New("post text must be 300 characters or fewer")
	case len(post.ImageUris) == 0 && len(post.VideoUris) == 0:
		return errors.New("at least one image or video is required")
	}

	if _, err := time.Parse(time.RFC3339, post.CreatedAt); err != nil {
		return errors.New("invalid datetime format")
	}

	return nil
}
