package storage

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

const filePath = "./data/musicstore.db?_foreign_keys=on"

func New() (*sql.DB, error) {
	storage, err := sql.Open("sqlite3", filePath)
	if err != nil {
		return nil, err
	}

	return storage, nil
}
