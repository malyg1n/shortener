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

// SetBatchLinkRequest model.
type SetBatchLinkRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// SetBatchLinkResponse model.
type SetBatchLinkResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}
