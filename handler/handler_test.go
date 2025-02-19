package handler

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ShareFrame/posting-service/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockATProtoClient struct {
	mock.Mock
}

func (m *MockATProtoClient) PostToFeed(post models.ShareFrameFeedPost, authToken, did string) (*models.PostResponse, error) {
	args := m.Called(post, authToken, did)
	if args.Get(0) != nil {
		return args.Get(0).(*models.PostResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestPostHandler(t *testing.T) {
	mockAtproto := new(MockATProtoClient)

	tests := []struct {
		name        string
		request     models.RequestPayload
		mockResp    *models.PostResponse
		mockErr     error
		expectErr   bool
		expectResp  *models.PostResponse
		mockCalled  bool
	}{
		{
			name: "Valid post request",
			request: models.RequestPayload{
				AuthToken: "valid_token",
				DID:       "did:example:123",
				Post: models.ShareFrameFeedPost{
					NSID:      "social.shareframe.feed.post",
					Text:      "Hello World!",
					ImageUris: []string{"https://example.com/image.jpg"},
					CreatedAt: time.Now().Format(time.RFC3339),
				},
			},
			mockResp: &models.PostResponse{
				URI:              "at://did:example:123/social.shareframe.feed.post/xyz",
				CID:              "bafyre123456",
				Commit:           models.Commit{CID: "commit123", Rev: "rev123"},
				ValidationStatus: "unknown",
			},
			mockErr:    nil,
			expectErr:  false,
			expectResp: &models.PostResponse{
				URI:              "at://did:example:123/social.shareframe.feed.post/xyz",
				CID:              "bafyre123456",
				Commit:           models.Commit{CID: "commit123", Rev: "rev123"},
				ValidationStatus: "unknown",
			},
			mockCalled: true,
		},
		{
			name: "Missing AuthToken",
			request: models.RequestPayload{
				AuthToken: "",
				DID:       "did:example:123",
				Post:      models.ShareFrameFeedPost{},
			},
			expectErr:  true,
			expectResp: nil,
			mockCalled: false,
		},
		{
			name: "Invalid post format",
			request: models.RequestPayload{
				AuthToken: "valid_token",
				DID:       "did:example:123",
				Post: models.ShareFrameFeedPost{
					NSID: "wrong.nsid",
					Text: "Invalid Post",
				},
			},
			expectErr:  true,
			expectResp: nil,
			mockCalled: false,
		},
		{
			name: "Failed to post due to API error",
			request: models.RequestPayload{
				AuthToken: "valid_token",
				DID:       "did:example:123",
				Post: models.ShareFrameFeedPost{
					NSID:      "social.shareframe.feed.post",
					Text:      "API failure test",
					ImageUris: []string{"https://example.com/image.jpg"},
					CreatedAt: time.Now().Format(time.RFC3339),
				},
			},
			mockResp:   nil,
			mockErr:    errors.New("failed to post to feed"),
			expectErr:  true,
			expectResp: nil,
			mockCalled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			if tt.mockCalled {
				mockAtproto.On("PostToFeed", mock.Anything, mock.Anything, mock.Anything).
					Return(tt.mockResp, tt.mockErr).Once()
			}

			resp, err := PostHandler(ctx, mockAtproto, tt.request)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectResp, resp)
			}

			if tt.mockCalled {
				mockAtproto.AssertExpectations(t)
			}
		})
	}
}

func TestValidatePost(t *testing.T) {
	tests := []struct {
		name      string
		post      models.ShareFrameFeedPost
		expectErr bool
	}{
		{
			name: "Valid post with image",
			post: models.ShareFrameFeedPost{
				NSID:      "social.shareframe.feed.post",
				Text:      "Hello World!",
				ImageUris: []string{"https://example.com/photo.jpg"},
				VideoUris: []string{},
				CreatedAt: time.Now().Format(time.RFC3339),
			},
			expectErr: false,
		},
		{
			name: "Valid post with video",
			post: models.ShareFrameFeedPost{
				NSID:      "social.shareframe.feed.post",
				Text:      "Check this out!",
				ImageUris: []string{},
				VideoUris: []string{"https://example.com/video.mp4"},
				CreatedAt: time.Now().Format(time.RFC3339),
			},
			expectErr: false,
		},
		{
			name: "Invalid NSID",
			post: models.ShareFrameFeedPost{
				NSID:      "invalid.nsid",
				Text:      "Wrong NSID",
				ImageUris: []string{"https://example.com/photo.jpg"},
				CreatedAt: time.Now().Format(time.RFC3339),
			},
			expectErr: true,
		},
		{
			name: "Exceeds max text length",
			post: models.ShareFrameFeedPost{
				NSID:      "social.shareframe.feed.post",
				Text:      string(make([]byte, 301)),
				ImageUris: []string{"https://example.com/photo.jpg"},
				CreatedAt: time.Now().Format(time.RFC3339),
			},
			expectErr: true,
		},
		{
			name: "Missing image and video",
			post: models.ShareFrameFeedPost{
				NSID:      "social.shareframe.feed.post",
				Text:      "No media attached",
				ImageUris: []string{},
				VideoUris: []string{},
				CreatedAt: time.Now().Format(time.RFC3339),
			},
			expectErr: true,
		},
		{
			name: "Invalid image format",
			post: models.ShareFrameFeedPost{
				NSID:      "social.shareframe.feed.post",
				Text:      "Invalid format",
				ImageUris: []string{"https://example.com/photo.pdf"},
				CreatedAt: time.Now().Format(time.RFC3339),
			},
			expectErr: true,
		},
		{
			name: "Invalid video format",
			post: models.ShareFrameFeedPost{
				NSID:      "social.shareframe.feed.post",
				Text:      "Invalid format",
				VideoUris: []string{"https://example.com/video.avi"},
				CreatedAt: time.Now().Format(time.RFC3339),
			},
			expectErr: true,
		},
		{
			name: "Invalid datetime format",
			post: models.ShareFrameFeedPost{
				NSID:      "social.shareframe.feed.post",
				Text:      "Wrong timestamp",
				ImageUris: []string{"https://example.com/photo.jpg"},
				CreatedAt: "invalid-date",
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePost(tt.post)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestIsValidExtension(t *testing.T) {
	tests := []struct {
		name      string
		uri       string
		allowed   []string
		expectRes bool
	}{
		{"Valid image - jpg", "https://example.com/photo.jpg", []string{".jpg", ".jpeg", ".png"}, true},
		{"Valid image - png", "https://example.com/photo.png", []string{".jpg", ".jpeg", ".png"}, true},
		{"Valid image - heic", "https://example.com/photo.heic", []string{".jpg", ".jpeg", ".png", ".heic", ".heif"}, true},
		{"Invalid image - pdf", "https://example.com/photo.pdf", []string{".jpg", ".jpeg", ".png"}, false},
		{"Valid video - mp4", "https://example.com/video.mp4", []string{".mp4", ".mov"}, true},
		{"Invalid video - avi", "https://example.com/video.avi", []string{".mp4", ".mov"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidExtension(tt.uri, tt.allowed)
			assert.Equal(t, tt.expectRes, result)
		})
	}
}
