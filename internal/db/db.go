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
		&models.Currency{},
		&models.Account{},
		&models.Transaction{},
		&models.LedgerRecord{},
	)
}

func DropSchema(db *gorm.DB) {
	db.Exec("DROP SCHEMA public CASCADE; CREATE SCHEMA public;")
}

func LoadAccountsData(db *gorm.DB, accounts []*models.Account) {
	for _, account := range accounts {
		if err := db.Create(account).Error; err != nil {
			log.Println(err)
			log.Fatalln("Unexpected ERROR while it's saving the following account -> ", account.ID, " | ", account.Email)
		}
	}
}

func LoadCurrenciesData(db *gorm.DB, currencies []*models.Currency) {
	for _, currency := range currencies {
		if err := db.Create(currency).Error; err != nil {
			log.Println(err)
			log.Fatalln("Unexpected ERROR while it's saving the following currency -> ", currency.ID, " | ", currency.Symbol)
		}
	}
}

func LoadTransactionsData(db *gorm.DB, txs []*models.Transaction) {
	for _, tx := range txs {
		if err := db.Create(tx).Error; err != nil {
			log.Println(err)
			log.Fatalln("Unexpected ERROR while it's saving the following currency -> ", tx.ID, " | ", tx.Concept)
		}
	}
}
