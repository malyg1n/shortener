package models

// SetLinkRequest model.
type SetLinkRequest struct {
	URL string `json:"url"`
}

// SetLinkResponse model.
type SetLinkResponse struct {
	Result string `json:"result"`
}

// LinkResponse model.
type LinkResponse struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
