package pkg

import (
	"database/sql"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// Database is a type that holds the DB/ORM flavour and the underlying sql.DB
// which can be used to create additional DBs or transactions.
type Database struct {
	Gorm *gorm.DB
	DB   *sql.DB
}

// Transaction generates a new transaction.
func (d *Database) Transaction() (*gorm.DB, error) {
	tx := d.Gorm.Begin()

	return tx, tx.Error
}

// CommitOrRollback takes a transaction and an error. If the error is nil, the
// transaction will be committed, otherwise it will be rolled back. If the
// commit or rollback fails a new error is returned, wrapping the original error
// if that's not nil.
func CommitOrRollback(tx *gorm.DB, err error) error {
	if err != nil {
		if tErr := tx.Rollback().Error; tErr != nil {
			return errors.Wrap(err, tErr.Error())
		}

		return err
	}

	if tErr := tx.Commit().Error; tErr != nil {
		return tErr
	}

	return nil
}
