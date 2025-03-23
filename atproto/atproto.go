package atproto

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ShareFrame/posting-service/models"
	"github.com/sirupsen/logrus"
)

type ATProtoClient interface {
	PostToFeed(post models.ShareFrameFeedPost, authToken, did string) (*models.PostResponse, error)
}

type ATProtoService struct {
	client *http.Client
}

func NewATProtoService(client *http.Client) *ATProtoService {
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}
	return &ATProtoService{client: client}
}

func (s *ATProtoService) PostToFeed(post models.ShareFrameFeedPost, authToken, did string) (*models.PostResponse, error) {
	const postURL = "https://shareframe.social/xrpc/com.atproto.repo.createRecord"

	payload, err := json.Marshal(models.CreateRecordRequest{
		Repo:       did,
		Collection: "social.shareframe.feed.post",
		Record:     post,
	})

	if err != nil {
		logrus.WithError(err).Error("Failed to marshal JSON payload")
		return nil, fmt.Errorf("failed to marshal request payload: %w", err)
	}

	req, err := http.NewRequest("POST", postURL, bytes.NewReader(payload))
	if err != nil {
		logrus.WithError(err).Error("Failed to create HTTP request")
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+authToken)

	resp, err := s.client.Do(req)
	if err != nil {
		logrus.WithError(err).Error("HTTP request failed")
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.WithError(err).Error("Failed to read response body")
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		logrus.WithField("status", resp.StatusCode).Error("Failed to post to feed")
		return nil, fmt.Errorf("failed to post: %s", string(body))
	}

	var postResponse models.PostResponse
	if err := json.Unmarshal(body, &postResponse); err != nil {
		logrus.WithError(err).Error("Failed to parse response JSON")
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &postResponse, nil
}
