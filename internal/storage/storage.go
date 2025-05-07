package storage

//создание и управление бд

import (
	"database/sql"
	"log"
)

var DB *sql.DB

func NewPostgresDB() (*sql.DB, error) {
	dsn := "postgres://user:password@db:5432/user_service?sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	log.Println("Connected to the database")
	DB = db
	return db, nil
}
