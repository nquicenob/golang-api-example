package models

import "github.com/jinzhu/gorm"

type CurrencyType string

const (
	CurrencyTypeEUR CurrencyType = "EUR"
)

type Currency struct {
	gorm.Model
	Symbol        CurrencyType `gorm:"type:varchar(5);unique_index"`
	Name          string
	LedgerRecords []*LedgerRecord `gorm:"foreignkey:CurrencyID"`
}
