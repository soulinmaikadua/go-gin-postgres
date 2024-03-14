package models

import (
	"database/sql"
	"fmt"
)

var db *sql.DB

func InitDB(dataSourceName string) error {
	var err error
	db, err = sql.Open("postgres", dataSourceName)
	if err != nil {
		return err
	}
	if err := db.Ping(); err != nil {
		return err
	}
	fmt.Println("Connected to database")
	return nil
}

func CloseDB() {
	db.Close()
}
