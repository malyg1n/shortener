package model

// Link base model.
type Link struct {
	ShortURL    string
	OriginalURL string
	UserUUID    string
	IsDeleted   bool
}
