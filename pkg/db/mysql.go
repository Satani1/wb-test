package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type Repository struct {
	db *sql.DB
}

func NewDB(url string) (*Repository, error) {
	db, err := sql.Open("mysql", url)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &Repository{
		db,
	}, nil
}

func (r *Repository) Close() {
	r.db.Close()
}
