package main

import (
	"fmt"
	"log"
	"os"

	"nquicenob.com/golang-api-example/internal/config"
	"nquicenob.com/golang-api-example/internal/db"
	"nquicenob.com/golang-api-example/internal/handlers"
	"nquicenob.com/golang-api-example/internal/models"
	"nquicenob.com/golang-api-example/internal/server"
	"nquicenob.com/golang-api-example/internal/services"
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
	sr := models.NewTransactionsRepository(dbConn, c)
	ar := models.NewAccountRepository(dbConn)
	st := services.NewTransactionsService(sr, ar)
	h := handlers.NewTransactionHandler(st)

	e.POST("/accounts/:account_id/movemoney", h.Create)
	e.GET("/accounts/:account_id", h.Find)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", c.Port)))
}
