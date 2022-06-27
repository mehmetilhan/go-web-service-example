package database

import (
	"database/sql"
	"log"
)

var DbConn *sql.DB

func SetupDatabase() {
	var err error
	DbConn, err = sql.Open("mysql", "root:secret@tcp(127.0.0.1:3306)/go_test")

	if err != nil {
		log.Fatal(err)
	}
}
