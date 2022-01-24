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
	db.MustExec(schema)

	return &LinksStoragePG{
		db: db,
	}, nil
}

// SetLink stores link into database.
func (l LinksStoragePG) SetLink(ctx context.Context, id, link, userUUID string) {
	tx := l.db.MustBegin()
	tx.MustExecContext(ctx, "insert into links (user_id, link_id, original_link) VALUES ($1, $2, $3)", userUUID, id, link)
	_ = tx.Commit()
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

// Close database handler.
func (l LinksStoragePG) Close() error {
	return l.db.Close()
}
