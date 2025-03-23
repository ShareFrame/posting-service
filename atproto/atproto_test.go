package atproto

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/ShareFrame/posting-service/models"
	"github.com/stretchr/testify/assert"
)

type mockTransport struct {
	roundTripFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.roundTripFunc(req)
}

func TestPostToFeed(t *testing.T) {
	now := time.Now().Format(time.RFC3339)

	tests := []struct {
		name           string
		post           models.ShareFrameFeedPost
		authToken      string
		did            string
		mockResponse   string
		mockStatusCode int
		mockErr        error
		expectErr      bool
		expectedResp   *models.PostResponse
	}{
		{
			name: "Successful post",
			post: models.ShareFrameFeedPost{
				NSID:      "social.shareframe.feed.post",
				Text:      "Valid Post",
				ImageUris: []string{"https://example.com/image.jpg"},
				CreatedAt: now,
			},
			authToken: "valid_token",
			did:       "did:example:123",
			mockResponse: `{
				"uri": "at://did:example:123/social.shareframe.feed.post/xyz",
				"cid": "bafyre123456",
				"commit": { "cid": "commit123", "rev": "rev123" },
				"validationStatus": "unknown"
			}`,
			mockStatusCode: http.StatusOK,
			expectErr:      false,
			expectedResp: &models.PostResponse{
				URI:              "at://did:example:123/social.shareframe.feed.post/xyz",
				CID:              "bafyre123456",
				Commit:           models.Commit{CID: "commit123", Rev: "rev123"},
				ValidationStatus: "unknown",
			},
		},
		{
			name: "Failed to marshal request payload",
			post: models.ShareFrameFeedPost{
				NSID: "social.shareframe.feed.post",
				Text: string(make([]byte, 1<<20)),
			},
			authToken: "valid_token",
			did:       "did:example:123",
			expectErr: true,
		},
		{
			name: "HTTP request failure",
			post: models.ShareFrameFeedPost{
				NSID:      "social.shareframe.feed.post",
				Text:      "Hello World!",
				ImageUris: []string{"https://example.com/image.jpg"},
				CreatedAt: now,
			},
			authToken: "valid_token",
			did:       "did:example:123",
			mockErr:   errors.New("network error"),
			expectErr: true,
		},
		{
			name: "Non-200 response",
			post: models.ShareFrameFeedPost{
				NSID:      "social.shareframe.feed.post",
				Text:      "Hello World!",
				ImageUris: []string{"https://example.com/image.jpg"},
				CreatedAt: now,
			},
			authToken:      "valid_token",
			did:            "did:example:123",
			mockResponse:   `{"error":"InvalidToken","message":"Token could not be verified"}`,
			mockStatusCode: http.StatusUnauthorized,
			expectErr:      true,
		},
		{
			name: "Invalid response JSON",
			post: models.ShareFrameFeedPost{
				NSID:      "social.shareframe.feed.post",
				Text:      "Hello World!",
				ImageUris: []string{"https://example.com/image.jpg"},
				CreatedAt: now,
			},
			authToken:      "valid_token",
			did:            "did:example:123",
			mockResponse:   `{invalid_json}`,
			mockStatusCode: http.StatusOK,
			expectErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTransport := &mockTransport{
				roundTripFunc: func(req *http.Request) (*http.Response, error) {
					if tt.mockErr != nil {
						return nil, tt.mockErr
					}
					res := &http.Response{
						StatusCode: tt.mockStatusCode,
						Header:     make(http.Header),
						Body:       io.NopCloser(strings.NewReader(tt.mockResponse)),
					}
					return res, nil
				},
			}

			mockClient := &http.Client{Transport: mockTransport}
			service := NewATProtoService(mockClient)

			resp, err := service.PostToFeed(tt.post, tt.authToken, tt.did)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResp, resp)
			}
		})
	}
}
