package database

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/bombsimon/team-betting/pkg"

	"github.com/doug-martin/goqu/v7"

	_ "github.com/doug-martin/goqu/v7/dialect/mysql"
	_ "github.com/go-sql-driver/mysql"
)

type ErrorType int

const (
	ErrUnknown ErrorType = iota
	ErrDuplicateKey
	ErrForeignKeyConstraint
	ErrUnknownColumn
	ErrNoDefaultValue
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

func ErrType(err error) ErrorType {
	errMap := map[string]ErrorType{
		"1062": ErrDuplicateKey,
		"1216": ErrForeignKeyConstraint, // Cannot add or update a child row
		"1452": ErrForeignKeyConstraint,
		"1054": ErrUnknownColumn,
		"1364": ErrNoDefaultValue,
	}

	for id, e := range errMap {
		if strings.Contains(err.Error(), fmt.Sprintf("Error %s:", id)) {
			return e
		}
	}

	return ErrUnknown
}
