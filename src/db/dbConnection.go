package main

import (
	"database/sql"
	"fmt"
)

// CreateCon for setting up mysql connection
func CreateCon() *sql.DB {
	db, err := sql.Open("mysql", "root:@tcp(db:3306)/gps")

	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(db.Ping())
	}

	return db
}
