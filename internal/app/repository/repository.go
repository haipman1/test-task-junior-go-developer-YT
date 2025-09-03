package repository

//Layer between App and DataBase

import (
	"database/sql"
	"fmt"
	"log"
)

func NewDB(databaseName string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseName)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to the database!")
	return db, nil
}
