package handler

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/ShareFrame/posting-service/atproto"
	"github.com/ShareFrame/posting-service/models"
	"github.com/sirupsen/logrus"
)

func PostHandler(ctx context.Context, client atproto.ATProtoClient, request models.RequestPayload) (*models.PostResponse, error) {
	if request.AuthToken == "" || request.DID == "" {
		err := errors.New("invalid request: missing 'authToken' or 'did'")
		logrus.Error(err)
		return nil, err
	}

	if err := validatePost(request.Post); err != nil {
		logrus.WithError(err).WithField("NSID", request.Post.NSID).Error("Validation failed")
		return nil, fmt.Errorf("invalid post: %w", err)
	}

	postResponse, err := client.PostToFeed(request.Post, request.AuthToken, request.DID)
	if err != nil {
		logrus.WithError(err).WithField("DID", request.DID).Error("Failed to post to feed")
		return nil, fmt.Errorf("posting to feed failed: %w", err)
	}

	logrus.WithField("DID", request.DID).Info("Post successfully created")
	return postResponse, nil
}

func validatePost(post models.ShareFrameFeedPost) error {
	if post.NSID != "social.shareframe.feed.post" {
		return errors.New("invalid NSID: only social.shareframe.feed.post is allowed")
	}
	if len(post.Text) > 300 {
		return errors.New("post text must be 300 characters or fewer")
	}
	if len(post.ImageUris) == 0 && len(post.VideoUris) == 0 {
		return errors.New("at least one image or video is required")
	}

	allowedImages := []string{".jpg", ".jpeg", ".png", ".gif", ".heic", ".heif"}
	allowedVideos := []string{".mp4", ".mov", ".webm"}

	for _, uri := range post.ImageUris {
		if !isValidExtension(uri, allowedImages) {
			return fmt.Errorf("invalid image format: %s (only jpg, jpeg, png, gif, heic, heif allowed)", filepath.Ext(uri))
		}
	}

	for _, uri := range post.VideoUris {
		if !isValidExtension(uri, allowedVideos) {
			return fmt.Errorf("invalid video format: %s (only mp4, mov, webm allowed)", filepath.Ext(uri))
		}
	}

	if _, err := time.Parse(time.RFC3339, post.CreatedAt); err != nil {
		return errors.New("invalid datetime format")
	}

	return nil
}

func isValidExtension(uri string, allowed []string) bool {
	ext := strings.ToLower(filepath.Ext(uri))
	for _, validExt := range allowed {
		if ext == validExt {
			return true
		}
	}
	return false
}
