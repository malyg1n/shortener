package pgsql

import (
	"context"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/malyg1n/shortener/model"
	"github.com/malyg1n/shortener/pkg/config"
	"github.com/malyg1n/shortener/pkg/errs"
	"github.com/malyg1n/shortener/storage"
	dbModel "github.com/malyg1n/shortener/storage/pgsql/model"
)

var (
	_      storage.LinksStorage = (*LinksStoragePG)(nil)
	schema                      = `create table if not exists links (
		id bigserial not null,
		user_id uuid null,
		link_id uuid not null unique,
		original_link text not null 
	);`
)

type LinksStoragePG struct {
	db *sqlx.DB
}

// NewLinksStoragePG creates new LinksStoragePG instance.
func NewLinksStoragePG(ctx context.Context) (*LinksStoragePG, error) {
	cfg := config.GetConfig()
	db, err := sqlx.Open("postgres", cfg.DatabaseDSN)
	if err != nil {
		return nil, err
	}
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(schema)
	if err != nil {
		return nil, err
	}

	return &LinksStoragePG{
		db: db,
	}, nil
}

// SetLink stores link into database.
func (l LinksStoragePG) SetLink(ctx context.Context, id, link, userUUID string) error {
	_, err := l.db.ExecContext(ctx, "insert into links (user_id, link_id, original_link) values ($1, $2, $3)", userUUID, id, link)
	return err
}

// GetLink returns link from database by id.
func (l LinksStoragePG) GetLink(ctx context.Context, id string) (string, error) {
	var link string
	err := l.db.GetContext(ctx, &link, "select original_link from links where link_id = $1", id)

	return link, err
}

// GetLinksByUser returns links by user.
func (l LinksStoragePG) GetLinksByUser(ctx context.Context, userUUID string) ([]model.Link, error) {
	var dbLinks []dbModel.Link
	var links []model.Link
	err := l.db.SelectContext(ctx, &dbLinks, "select * from links where user_id = $1", userUUID)
	if err != nil {
		return nil, err
	}

	if len(dbLinks) == 0 {
		return nil, errs.ErrNotFound
	}

	for _, dbLink := range dbLinks {
		links = append(links, dbLink.ToCanonical())
	}

	return links, err
}

// GetLinkByOriginal returns links by original link.
func (l LinksStoragePG) GetLinkByOriginal(ctx context.Context, url string) (string, error) {
	var link string
	err := l.db.GetContext(ctx, &link, "select link_id from links where original_link = $1", url)

	return link, err
}

// SetBatchLinks set links from collection.
func (l LinksStoragePG) SetBatchLinks(ctx context.Context, links []model.Link, userUUID string) error {
	tx, err := l.db.Begin()
	if err != nil {
		return err
	}
	for _, link := range links {
		_, err = tx.ExecContext(
			ctx,
			"insert into links (user_id, link_id, original_link) values ($1, $2, $3)",
			userUUID,
			link.ShortURL,
			link.OriginalURL,
		)

		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// Close database handler.
func (l LinksStoragePG) Close() error {
	return l.db.Close()
}

// Ping database.
func (l LinksStoragePG) Ping() error {
	return l.db.Ping()
}
