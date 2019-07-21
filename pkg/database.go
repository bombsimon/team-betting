package pkg

import (
	"database/sql"

	"github.com/doug-martin/goqu/v7"
	"github.com/pkg/errors"
)

// Querier is the interface that implements From() which is use for
// goqu.Database and goqu.TxDatabase.
type Querier interface {
	From(from ...interface{}) *goqu.Dataset
}

// Database is a type that holds the goqu.Database and the underlying sql.DB
// which was used to create the goqu.Database.
type Database struct {
	Gq *goqu.Database
	DB *sql.DB
}

// Transaction generats a new transaction.
func (d *Database) Transaction() (*goqu.TxDatabase, error) {
	return d.Gq.Begin()
}

// CommitOrRollback takes a transaction and an error. If the error is nil, the
// transaction will be committed, otherwise it will be rolled back. If the
// commit or rollback fails a new error is returned, wrapping the original error
// if that's not nil.
func CommitOrRollback(tx *goqu.TxDatabase, err error) error {
	if err != nil {
		if tErr := tx.Rollback(); tErr != nil {
			return errors.Wrap(err, tErr.Error())
		}

		return err
	}

	if tErr := tx.Commit(); tErr != nil {
		return tErr
	}

	return nil
}
