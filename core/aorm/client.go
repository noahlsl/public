package aorm

import (
	"context"
	"github.com/doug-martin/goqu/v9"

	"github.com/jmoiron/sqlx"
	"github.com/noahlsl/public/core/dbx"
	"github.com/pkg/errors"
)

type Client struct {
	db *sqlx.DB
	g  goqu.DialectWrapper
}

func NewClient(dsn string, poolSize ...int) *Client {
	db := dbx.MustDB(dsn, poolSize...)
	g := goqu.Dialect("mysql")
	return &Client{
		db: db,
		g:  g,
	}
}

func (s *Client) Tx(ctx context.Context) (*Tx, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return newTx(ctx, tx), nil
}
func (s *Client) Close() error {
	return s.db.Close()
}
