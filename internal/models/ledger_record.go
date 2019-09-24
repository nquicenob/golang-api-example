package models

import "github.com/jinzhu/gorm"

type LedgerRecordType string

const (
	LedgerRecordDebit  LedgerRecordType = "DEBIT"
	LedgerRecordCredit LedgerRecordType = "CREDIT"
)

type LedgerRecord struct {
	gorm.Model

	Type             LedgerRecordType `gorm:"not null"`
	Balance          string           `gorm:"type:numeric(15,4);not null"`
	Amount           string           `gorm:"type:numeric(15,4);not null"`
	PreviousBanlance string           `gorm:"type:numeric(15,4);not null"`

	TransactionID string `gorm:"type:varchar(30) REFERENCES transactions(id);not null"`
	AccountID     string `gorm:"type:varchar(30) REFERENCES accounts(id);not null"`
	CurrencyID    uint   `gorm:"type:bigint REFERENCES currencies(id);not null"`
}
