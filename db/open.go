package db

import (
	"database/sql"

	"github.com/afteralec/grpc-user/db/query"
)

type OpenOutput struct {
	DB      *sql.DB
	Queries *query.Queries
}

func Open(url string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", url)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("PRAGMA JOURNAL_MODE = WAL;")
	if err != nil {
		return nil, err
	}
	_, err = db.Exec("PRAGMA SYNCHRONOUS = normal;")
	if err != nil {
		return nil, err
	}
	_, err = db.Exec("PRAGMA TEMP_STORE = memory;")
	if err != nil {
		return nil, err
	}
	_, err = db.Exec("PRAGMA MMAP_SIZE = 30000000000;")
	if err != nil {
		return nil, err
	}
	_, err = db.Exec("PRAGMA PAGE_SIZE = 32768;")
	if err != nil {
		return nil, err
	}
	_, err = db.Exec("PRAGMA FOREIGN_KEYS = ON;")
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
