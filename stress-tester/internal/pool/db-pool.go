package pool

import (
	"database/sql"
	"log/slog"

	_ "github.com/mattn/go-sqlite3"
)

// GetDb initializes and returns a new in-memory SQLite database connection.
// It panics if the database cannot be opened.

func GetDb() *sql.DB {

	db, err := sql.Open("sqlite3", "file::memory:?")

	if err != nil {
		slog.Error("pool.GetDb", "msg", err.Error())
		panic(err)
	}
	return db
}
