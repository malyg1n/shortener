package model

import (
	"database/sql"
	"github.com/malyg1n/shortener/model"
)

// Link is a DB representation of model.Link canonical model.
type Link struct {
	ID          uint           `db:"id"`
	UserID      sql.NullString `db:"user_id"`
	LinkID      string         `db:"link_id"`
	OriginalURL string         `db:"original_link"`
}

// ToCanonical converts a DB object to canonical model.
func (l Link) ToCanonical() model.Link {
	return model.Link{
		ShortURL:    l.LinkID,
		OriginalURL: l.OriginalURL,
	}
}
