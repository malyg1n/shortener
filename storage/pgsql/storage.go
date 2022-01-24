package pgsql

import (
	"context"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/malyg1n/shortener/model"
	"github.com/malyg1n/shortener/pkg/config"
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

func NewLinkStoragePG() (*LinksStoragePG, error) {
	cfg := config.GetConfig()
	db, err := sqlx.Open("postgres", cfg.DatabaseDSN)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	db.MustExec(schema)

	return &LinksStoragePG{
		db: db,
	}, nil
}

func (l LinksStoragePG) SetLink(ctx context.Context, id, link, userUUID string) {
	tx := l.db.MustBegin()
	tx.MustExec("insert into links (user_id, link_id, original_link) VALUES ($1, $2, $3)", userUUID, id, link)
	_ = tx.Commit()
}

func (l LinksStoragePG) GetLink(ctx context.Context, id string) (string, error) {
	var link string
	err := l.db.Get(&link, "select original_link from links where link_id = $1", id)

	return link, err
}

func (l LinksStoragePG) GetLinksByUser(ctx context.Context, userUUID string) ([]model.Link, error) {
	var dbLinks []dbModel.Link
	var links []model.Link
	err := l.db.Select(&dbLinks, "select * from links where user_id = $1", userUUID)
	if err != nil {
		return nil, err
	}

	for _, dbLink := range dbLinks {
		links = append(links, dbLink.ToCanonical())
	}

	return links, err
}

func (l LinksStoragePG) Close() error {
	return l.db.Close()
}
