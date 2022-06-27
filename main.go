package main

import (
	_ "github.com/go-sql-driver/mysql"
	"log"
	"mehmet.com/database"
	"mehmet.com/person"
	"mehmet.com/storage"
	"net/http"
)

const basePath = "/api"

func main() {

	database.SetupDatabase()
	person.SetupRoutes(basePath)
	storage.SetupRoutes(basePath)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
