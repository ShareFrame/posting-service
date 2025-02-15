package models

type ShareFrameFeedPost struct {
	Text      string   `json:"text,omitempty"`
	ImageUris []string `json:"imageUris,omitempty"`
	VideoUris []string `json:"videoUris,omitempty"`
	CreatedAt string   `json:"createdAt"`
	NSID      string   `json:"nsid"`
}

type RequestPayload struct {
	AuthToken string             `json:"authToken"`
	DID       string             `json:"did"`
	Post      ShareFrameFeedPost `json:"post"`
}
