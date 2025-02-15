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

type PostResponse struct {
	URI              string `json:"uri"`
	CID              string `json:"cid"`
	Commit          Commit `json:"commit"`
	ValidationStatus string `json:"validationStatus"`
}

type Commit struct {
	CID string `json:"cid"`
	Rev string `json:"rev"`
}
