package aorm

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
)

type Tx struct {
	ctx      context.Context
	tx       *sql.Tx
	sqlSlice []string
}

func newTx(ctx context.Context, tx *sql.Tx) *Tx {
	return &Tx{
		ctx: ctx,
		tx:  tx,
	}
}
func (s *Tx) Add(in TxModel) *Tx {
	s.sqlSlice = append(s.sqlSlice, in.SQL())
	return s
}

func (s *Tx) Commit() error {
	for _, query := range s.sqlSlice {
		_, err := s.tx.ExecContext(s.ctx, query)
		if err != nil {
			_ = s.tx.Rollback()
			return errors.WithStack(err)
		}
	}
	s.sqlSlice = nil
	return s.tx.Commit()
}
