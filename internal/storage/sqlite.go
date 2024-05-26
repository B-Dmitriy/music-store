package storage

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

const filePath = "./data/musicstore.db"

func New() (*sql.DB, error) {
	storage, err := sql.Open("sqlite", filePath)
	if err != nil {
		return nil, err
	}

	return storage, nil
}
