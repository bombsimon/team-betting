package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bombsimon/team-betting/pkg"
	"github.com/jinzhu/gorm"

	// Import MySQL dialects for side effects.
	_ "github.com/go-sql-driver/mysql"
)

// ErrorType represents a type for database errors.
type ErrorType int

// Known errors from the database (MySQL)
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
	if dsn == "" {
		dsn = "betting:betting@/betting?charset=utf8&parseTime=True&loc=Local"
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	gdb, err := gorm.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	gdb.SingularTable(true)

	if strings.EqualFold(os.Getenv("LOG_LEVEL"), "debug") {
		logger := log.New(os.Stdout, "[DEBUG] ", log.LstdFlags)

		gdb.LogMode(true)
		gdb.SetLogger(logger)
	}

	return &pkg.Database{
		Gorm: gdb,
		DB:   db,
	}
}

// ErrType will return an ErrorType from a regular error from the database.
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
