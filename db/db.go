package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	db *sql.DB
}

func New(dbPath string) (*Database, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	if err := createTables(db); err != nil {
		return nil, err
	}
	return &Database{db: db}, nil
}

func createTables(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS products (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name VARCHAR,
			description VARCHAR,
			price DECIMAL,
			quantity INTEGER DEFAULT 0,
			image_url VARCHAR,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS cart_items (
			user_id INTEGER ,
			product_id INTEGER,
			quantity INTEGER NOT NULL DEFAULT 1,
			PRIMARY KEY (user_id, product_id),
			FOREIGN KEY (product_id) REFERENCES products(id)
		)`,
		`CREATE TABLE IF NOT EXISTS admins (
			user_id INTEGER PRIMARY KEY
		)`,
	}
	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return err
		}
	}
	return nil
}
