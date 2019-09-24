package models

import (
	"github.com/jinzhu/gorm"
)

type AccountRepository interface {
	FindAccountID(*Account) error
}

func NewAccountRepository(db *gorm.DB) AccountRepository {
	return accountRepository{
		db: db,
	}
}

type accountRepository struct {
	db *gorm.DB
}

func (tr accountRepository) FindAccountID(a *Account) (err error) {
	if err := tr.db.First(&a, "id = ?", a.ID).Error; err != nil {
		return err
	}
	return nil
}
