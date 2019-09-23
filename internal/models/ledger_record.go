package models

import "github.com/jinzhu/gorm"

type LedgerRecordType string

const (
	LedgerRecordDebit   LedgerRecordType = "DEBIT"
	LedgerRecordCreedit LedgerRecordType = "CREDIT"
)

type LedgerRecord struct {
	gorm.Model
	Amount           string           `gorm:"not null"`
	Type             LedgerRecordType `gorm:"not null"`
	Balance          string           `gorm:"not null"`
	PreviousBanlance string           `gorm:"not null"`
	TransactionID    string           `gorm:"type:varchar(30) REFERENCES transactions(id);not null"`
	AccountID        uint             `gorm:"type:varchar(30) REFERENCES accounts(id);not null"`
}
