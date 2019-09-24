package models

import "time"

type Account struct {
	ID            string `gorm:"type:varchar(30);primary_key"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time `sql:"index"`
	Name          string
	Email         string          `gorm:"type:varchar(100);unique_index"`
	Balance       string          `gorm:"type:numeric(15,4);not null"`
	LedgerRecords []*LedgerRecord `gorm:"foreignkey:AccountID"`
}
