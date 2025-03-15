package db

import (
	"context"
	"database/sql"
	"log/slog"
	"stress-tester/internal/dto"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	db    *sql.DB
	input chan *dto.Red
}

// NewDB initializes a new DB instance with the provided SQL database connection
// and input channel for *dto.Red. It ensures that the 'red' table exists in
// the database, creating it if necessary. The table includes fields for target,
// sent_at, received_at, status_code, and duration.

func NewDB(db *sql.DB, input chan *dto.Red) *DB {
	db.Exec("CREATE TABLE IF NOT EXISTS red (target text, sent_at timestamp, received_at timestamp, status_code int, duration int)")
	return &DB{
		db:    db,
		input: input,
	}
}

// Store takes a context and a channel of *dto.Red. It will consume all available
// values from the channel and insert them into the 'red' table in the database.
// It will stop consuming the channel when the context is canceled.
func (d *DB) Store(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			r := <-d.input
			_, err := d.db.Exec("INSERT INTO red (target, sent_at, received_at, status_code, duration) VALUES ( ?, ?, ?, ?, ?)", r.Target, r.SentAt, r.ReceivedAt, r.StatusCode, r.Duration)
			if err != nil {
				slog.Error("db.Store", "msg", err.Error())
			}
		}
	}
}

// getReds executes a query on the 'red' table and returns a slice of *dto.Red
// representing the results. If an error occurs while executing the query, or
// while scanning the results, an error is logged and the query continues to the
// next row.
func (d *DB) getReds(query string) []*dto.Red {
	rows, err := d.db.Query(query)
	if err != nil {
		slog.Error("db.getReds", "msg", err.Error())
	}
	defer rows.Close()
	var reds []*dto.Red
	for rows.Next() {
		r := &dto.Red{}
		err := rows.Scan(&r.Target, &r.SentAt, &r.ReceivedAt, &r.StatusCode, &r.Duration)
		if err != nil {
			slog.Error("db.getReds scan", "msg", err.Error())
		}
		reds = append(reds, r)
	}
	return reds
}

// GetAllReds retrieves all records from the 'red' table. It returns a slice of *dto.Red
// representing these records.
func (d *DB) GetAllReds() []*dto.Red {
	return d.getReds("SELECT target, sent_at, received_at, status_code, duration FROM red")
}

// GetRedsWithoutErrors retrieves all records from the 'red' table where the status code is 200,
// indicating that the request was successful. It returns a slice of *dto.Red representing these
// records.
func (d *DB) GetRedsWithoutErrors() []*dto.Red {
	return d.getReds("SELECT target, sent_at, received_at, status_code, duration FROM red where status_code = 200")
}

// GetRedWithErrors retrieves all records from the 'red' table where the status code is not 200,
// indicating that an error occurred. It returns a slice of *dto.Red representing these records.

func (d *DB) GetRedWithErrors() []*dto.Red {
	return d.getReds("SELECT target, sent_at, received_at, status_code, duration FROM red WHERE status_code != 200")
}

// Close closes the database connection. It will return any error it encounters.
func (d *DB) Close() error {
	return d.db.Close()
}

// ClearDatabase drops the 'red' table from the database. This action removes all 
// records stored in the 'red' table, effectively clearing its data.

func (d *DB) ClearDatabase() {
	d.db.Exec("drop table red")
}
