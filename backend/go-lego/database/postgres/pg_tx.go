package postgres

import (
	"github.com/jmoiron/sqlx"
)

type Tx struct {
	tx *sqlx.Tx
}

func newTx(tx *sqlx.Tx) *Tx {
	return &Tx{
		tx: tx,
	}
}

func (t *Tx) getTx() *sqlx.Tx {
	return t.tx
}

func (t *Tx) Commit() error {
	return t.tx.Commit()
}

func (t *Tx) Rollback() error {
	return t.tx.Rollback()
}
