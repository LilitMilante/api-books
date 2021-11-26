package database

import (
	"database/sql"
)

type Storage struct {
	conn *sql.DB
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		conn: db,
	}
}
