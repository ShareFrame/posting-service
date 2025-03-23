package models

type ShareFrameFeedPost struct {
	Text              string                   `json:"text,omitempty"`
	ImageUris         []string                 `json:"imageUris,omitempty"`
	VideoUris         []string                 `json:"videoUris,omitempty"`
	CreatedAt         string                   `json:"createdAt"`
	Likes             int                      `json:"likes,omitempty"`
	Shares            int                      `json:"shares,omitempty"`
	Comments          int                      `json:"comments,omitempty"`
	Rewatches         int                      `json:"rewatches,omitempty"`
	Saves             int                      `json:"saves,omitempty"`
	WatchTime         int                      `json:"watchTime,omitempty"`
	LocationString    string                   `json:"locationString,omitempty"`
	City              string                   `json:"city,omitempty"`
	Region            string                   `json:"region,omitempty"`
	Country           string                   `json:"country,omitempty"`
	TimeZone          string                   `json:"timeZone,omitempty"`
	Geohash           string                   `json:"geohash,omitempty"`
	TrendingScore     float64                  `json:"trendingScore,omitempty"`
	IsStory           bool                     `json:"isStory,omitempty"`
	ExpiresAt         string                   `json:"expiresAt,omitempty"`
	Language          string                   `json:"language,omitempty"`
	Tags              []string                 `json:"tags,omitempty"`
	Keywords          []string                 `json:"keywords,omitempty"`
	ReplyTo           string                   `json:"replyTo,omitempty"`
	QuoteOf           string                   `json:"quoteOf,omitempty"`
	AuthorDisplayName string                   `json:"authorDisplayName,omitempty"`
	AuthorHandle      string                   `json:"authorHandle,omitempty"`
	ImageMetadata     map[string]interface{}   `json:"imageMetadata,omitempty"`
	VideoMetadata     map[string]interface{}   `json:"videoMetadata,omitempty"`
	EditHistory       []map[string]interface{} `json:"editHistory,omitempty"`
	SourceApp         string                   `json:"sourceApp,omitempty"`
	NSID              string                   `json:"nsid"`
}

type CreateRecordRequest struct {
	Repo       string             `json:"repo"`
	Collection string             `json:"collection"`
	Record     ShareFrameFeedPost `json:"record"`
}

type RequestPayload struct {
	AuthToken string             `json:"authToken"`
	DID       string             `json:"did"`
	Post      ShareFrameFeedPost `json:"post"`
}

type PostResponse struct {
	URI              string `json:"uri"`
	CID              string `json:"cid"`
	Commit           Commit `json:"commit"`
	ValidationStatus string `json:"validationStatus"`
}

type Commit struct {
	CID string `json:"cid"`
	Rev string `json:"rev"`
}
