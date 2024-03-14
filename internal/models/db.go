package models

import (
	"database/sql"
	"fmt"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDB() error {
	// Construct the connection string
	connectionString := "postgres://tkaeaswb:XHTYs5BSzLXESRRScK5fArMk2RNFJau-@john.db.elephantsql.com/tkaeaswb"

	var err error
	db, err = sql.Open("postgres", connectionString)
	if err != nil {
		return err
	}
	// Ping the database to check if the connection is successful
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to database")
	return nil
}

func GetConnect() *sql.DB {
	return db
}

func CloseDB() {
	db.Close()
}
