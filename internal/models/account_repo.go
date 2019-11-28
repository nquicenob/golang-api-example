package models

import (
	"github.com/jinzhu/gorm"
)

type AccountRepository interface {
	FindAccountID(*Account) error
}

func NewAccountRepository(db *gorm.DB) AccountRepository {
	return &accountRepository{
		db: db,
	}
}

type accountRepository struct {
	db *gorm.DB
}

func (tr *accountRepository) FindAccountID(a *Account) (err error) {
	return tr.db.First(&a, "id = ?", a.ID).Error
}
