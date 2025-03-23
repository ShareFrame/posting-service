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

var (
	allowedImageExts = map[string]struct{}{
		".jpg": {}, ".jpeg": {}, ".png": {}, ".gif": {}, ".heic": {}, ".heif": {},
	}
	allowedVideoExts = map[string]struct{}{
		".mp4": {}, ".mov": {}, ".webm": {},
	}
)

func PostHandler(ctx context.Context, client atproto.ATProtoClient, request models.RequestPayload) (*models.PostResponse, error) {
	if request.AuthToken == "" || request.DID == "" {
		err := errors.New("invalid request: missing 'authToken' or 'did'")
		logrus.Error(err)
		return nil, err
	}

	request.Post.SourceApp = "ShareFrame"

	if request.Post.IsStory && request.Post.ExpiresAt == "" {
		request.Post.ExpiresAt = time.Now().Add(24 * time.Hour).UTC().Format(time.RFC3339)
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

	if postResponse == nil {
		logrus.Error("PostHandler returned nil postResponse with no error")
		return nil, fmt.Errorf("no response returned from ATProto")
	}	

	return postResponse, nil
}

func validatePost(post models.ShareFrameFeedPost) error {
	if post.NSID != "social.shareframe.feed.post" {
		return errors.New("invalid NSID: only social.shareframe.feed.post is allowed")
	}

	if len(post.Text) > 300 {
		return errors.New("post text must be 300 characters or fewer")
	}

	for _, uri := range post.ImageUris {
		if !isValidExtension(uri, allowedImageExts) {
			return fmt.Errorf("invalid image format: %s", filepath.Ext(uri))
		}
	}

	for _, uri := range post.VideoUris {
		if !isValidExtension(uri, allowedVideoExts) {
			return fmt.Errorf("invalid video format: %s", filepath.Ext(uri))
		}
	}

	if !isRFC3339(post.CreatedAt) {
		return errors.New("invalid datetime format for createdAt")
	}

	if post.ExpiresAt != "" && !isRFC3339(post.ExpiresAt) {
		return errors.New("invalid datetime format for expiresAt")
	}

	return nil
}

func isValidExtension(uri string, allowed map[string]struct{}) bool {
	ext := strings.ToLower(filepath.Ext(uri))
	_, ok := allowed[ext]
	return ok
}

func isRFC3339(s string) bool {
	_, err := time.Parse(time.RFC3339, s)
	return err == nil
}
