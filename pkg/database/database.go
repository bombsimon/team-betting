package database

import (
	"database/sql"

	"github.com/bombsimon/team-betting/pkg"
	"github.com/doug-martin/goqu"
)

// New will connect to a database based on passed DSN and return a new
// *pkg.Database type.
func New(dsn string) *pkg.Database {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	return &pkg.Database{
		Gq: goqu.New("mysql", db),
		DB: db,
	}
}
