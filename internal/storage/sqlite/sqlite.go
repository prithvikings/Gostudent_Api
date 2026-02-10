package sqlite

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/prithvikings/Gostudent_Api/internal/config"
)

type Sqlite struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*Sqlite, error) {
	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS STUDENTS(
	id INTEGER PRIMARY KEY AUTOINCREANMENT,
	NAME TEXT,
	emai TEXT
	age INTEGER
	)`)
	if err != nil {
		return nil, err
	}
	return &Sqlite{Db: db}, nil
}
