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

func PostToFeed(post models.ShareFrameFeedPost, authToken, did string) (string, error) {
	const postURL = "https://shareframe.social/xrpc/com.atproto.repo.createRecord"

	payload, err := json.Marshal(map[string]interface{}{
		"repo":       did,
		"collection": "social.shareframe.feed.post",
		"record":     post,
	})
	if err != nil {
		logrus.WithError(err).Error("Failed to marshal JSON payload")
		return "", fmt.Errorf("failed to marshal request payload: %w", err)
	}

	req, err := http.NewRequest("POST", postURL, bytes.NewReader(payload))
	if err != nil {
		logrus.WithError(err).Error("Failed to create HTTP request")
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+authToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		logrus.WithError(err).Error("HTTP request failed")
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.WithError(err).Error("Failed to read response body")
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		logrus.WithField("status", resp.StatusCode).Error("Failed to post to feed")
		return "", fmt.Errorf("failed to post: %s", string(body))
	}

	return string(body), nil
}
