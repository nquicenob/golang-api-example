package main

import (
	"log"
	"os"

	"nquicenob.com/golang-api-example/internal/config"
	"nquicenob.com/golang-api-example/internal/db"
)

func main() {
	log.Println("<------ start ----->")
	c, err := config.New()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	dbConn, err := db.New(c)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer dbConn.Close()

	log.Println("Dropping schema")
	db.DropSchema(dbConn)
	log.Println("The schema was dropped")

	log.Println("Migrating schema")
	db.AutoMigrate(dbConn)
	log.Println("The schema was migrated")

	log.Println("Loading data")
	loadData(dbConn)
	log.Println("Data was loaded")

	log.Println("<------ end ----->")
}
