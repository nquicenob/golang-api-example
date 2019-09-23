package main

import (
	"fmt"
	"log"
	"os"

	"nquicenob.com/golang-api-example/internal/config"
	"nquicenob.com/golang-api-example/internal/db"
	"nquicenob.com/golang-api-example/internal/server"
)

func main() {
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

	db.AutoMigrate(dbConn)

	e := server.New(c)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", c.Port)))
}
