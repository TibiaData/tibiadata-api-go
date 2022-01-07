package structs

// Fansites stores a list of promoted and supported fansites
type Fansites struct {
	PromotedFansites  []Fansite `json:"promoted"`
	SupportedFansites []Fansite `json:"supported"`
}

// Fansite stores a fansite information
type Fansite struct {
	Name           string      `json:"name"`
	LogoURL        string      `json:"logo_url"`
	Homepage       string      `json:"homepage"`
	Contact        string      `json:"contact"`
	ContentType    ContentType `json:"content_type"`
	SocialMedia    SocialMedia `json:"social_media"`
	Languages      []string    `json:"languages"`
	Specials       []string    `json:"specials"`
	FansiteItem    bool        `json:"fansite_item"`
	FansiteItemURL string      `json:"fansite_item_url"`
}

// ContentType stores a fansite content type information
type ContentType struct {
	Statistics bool `json:"statistics"`
	Texts      bool `json:"texts"`
	Tools      bool `json:"tools"`
	Wiki       bool `json:"wiki"`
}

// SocialMedia tells which social media a fansite has
type SocialMedia struct {
	Discord   bool `json:"discord"`
	Facebook  bool `json:"facebook"`
	Instagram bool `json:"instagram"`
	Reddit    bool `json:"reddit"`
	Twitch    bool `json:"twitch"`
	Twitter   bool `json:"twitter"`
	Youtube   bool `json:"youtube"`
}
