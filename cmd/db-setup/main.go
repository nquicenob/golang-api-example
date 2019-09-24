package main

import (
	"crypto/rand"
	"log"
	"os"
	"time"

	"github.com/oklog/ulid"
	"nquicenob.com/golang-api-example/internal/config"
	"nquicenob.com/golang-api-example/internal/db"
	"nquicenob.com/golang-api-example/internal/models"
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
	entropy := rand.Reader
	concept := "account origination"

	txs := []*models.Transaction{
		&models.Transaction{
			ID:      ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String(),
			Concept: concept,
		},
		&models.Transaction{
			ID:      ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String(),
			Concept: concept,
		},
		&models.Transaction{
			ID:      ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String(),
			Concept: concept,
		},
	}
	db.LoadTransactionsData(dbConn, txs)

	currencies := []*models.Currency{
		&models.Currency{
			Symbol: models.CurrencyTypeEUR,
			Name:   "Euro",
		},
	}
	db.LoadCurrenciesData(dbConn, currencies)
	db.LoadAccountsData(dbConn, []*models.Account{
		&models.Account{
			ID:      ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String(),
			Name:    "Nicolas Quiceno",
			Email:   "nquicenob@gmail.com",
			Balance: "0",
			LedgerRecords: []*models.LedgerRecord{
				&models.LedgerRecord{
					Type:             models.LedgerRecordCredit,
					Balance:          "0",
					Amount:           "0",
					PreviousBanlance: "0",
					CurrencyID:       currencies[0].ID,
					TransactionID:    txs[0].ID,
				},
			},
		},
		&models.Account{
			ID:      ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String(),
			Name:    "Karen Mercado",
			Email:   "kr@gmail.com",
			Balance: "0",
			LedgerRecords: []*models.LedgerRecord{
				&models.LedgerRecord{
					Type:             models.LedgerRecordCredit,
					Balance:          "0",
					Amount:           "0",
					PreviousBanlance: "0",
					CurrencyID:       currencies[0].ID,
					TransactionID:    txs[1].ID,
				},
			},
		},
		&models.Account{
			ID:      ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String(),
			Name:    "Santiago Mercado",
			Email:   "santy@gmail.com",
			Balance: "0",
			LedgerRecords: []*models.LedgerRecord{
				&models.LedgerRecord{
					Type:             models.LedgerRecordCredit,
					Balance:          "0",
					Amount:           "0",
					PreviousBanlance: "0",
					CurrencyID:       currencies[0].ID,
					TransactionID:    txs[2].ID,
				},
			},
		},
	})
	log.Println("Data was loaded")

	log.Println("<------ end ----->")
}
