package models

import "time"

type CurrencyType string

const (
	CurrencyTypeEUR CurrencyType = "EUR"
)

type Account struct {
	ID            string `gorm:"type:varchar(30);primary_key"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time      `sql:"index"`
	Currency      CurrencyType    `gorm:"not null"`
	UserID        uint            `gorm:"type:bigint REFERENCES users(id);not null"`
	LedgerRecords []*LedgerRecord `gorm:"foreignkey:AccountID"`
}
