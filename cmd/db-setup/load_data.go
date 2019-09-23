package main

import (
	"crypto/rand"
	"log"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/oklog/ulid"
	"nquicenob.com/golang-api-example/internal/models"
)

func loadData(db *gorm.DB) {
	entropy := rand.Reader
	users := []*models.User{
		&models.User{
			Name:  "Nicolas Quiceno",
			Email: "nquicenob@gmail.com",
			Accounts: []*models.Account{
				&models.Account{
					ID:       ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String(),
					Currency: models.CurrencyTypeEUR,
				},
			},
		},
		&models.User{
			Name:  "Karen Mercado",
			Email: "kr@gmail.com",
			Accounts: []*models.Account{
				&models.Account{
					ID:       ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String(),
					Currency: models.CurrencyTypeEUR,
				},
				&models.Account{
					ID:       ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String(),
					Currency: models.CurrencyTypeEUR,
				},
			},
		},
		&models.User{
			Name:  "Santiago Mercado",
			Email: "santy@gmail.com",
		},
	}
	for _, user := range users {
		if err := db.Create(user).Error; err != nil {
			log.Println(err)
			log.Fatalln("Unexpected ERROR while it's saving the following user -> ", user.ID, " | ", user.Email)
		}
	}
}
