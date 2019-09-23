package db

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"gopkg.in/matryer/try.v1"
	"nquicenob.com/golang-api-example/internal/config"
	"nquicenob.com/golang-api-example/internal/models"
)

func New(c config.Specification) (db *gorm.DB, err error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s sslmode=%s password=%s connect_timeout=%s",
		c.DBHost,
		c.DBPort,
		c.DBUser,
		c.DBName,
		c.DBSSLMode,
		c.DBPass,
		c.DBConnectTimeout,
	)
	try.MaxRetries = c.DBConnRetries
	err = try.Do(func(attempt int) (bool, error) {
		db, err = gorm.Open("postgres", connStr)
		if err != nil {
			nextRetry := time.Duration(attempt*2) * time.Duration(rand.Int31n(1000)) * time.Millisecond
			log.Println(err, "try again in: ", nextRetry)
			time.Sleep(nextRetry)
		}
		return attempt < c.DBConnRetries, err
	})

	return db, err
}

//TODO: err check
func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(
		&models.User{},
		&models.Account{},
		&models.Transaction{},
		&models.LedgerRecord{},
	)
}

func DropSchema(db *gorm.DB) {
	db.Exec("DROP SCHEMA public CASCADE; CREATE SCHEMA public;")
}
