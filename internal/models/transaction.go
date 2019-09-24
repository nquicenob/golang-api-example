package models

import "time"

type Transaction struct {
	ID            string `gorm:"type:varchar(30);primary_key"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time `sql:"index"`
	Concept       string
	LedgerRecords []*LedgerRecord `gorm:"foreignkey:TransactionID"`
}
