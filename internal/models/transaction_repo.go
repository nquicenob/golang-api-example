package models

import (
	"context"
	"crypto/rand"
	"database/sql"
	"fmt"
	"time"

	"nquicenob.com/golang-api-example/internal/config"

	"github.com/jinzhu/gorm"
	"github.com/oklog/ulid"
	"github.com/shopspring/decimal"
	"gopkg.in/matryer/try.v1"

	mrand "math/rand"
)

type TransactionsRepository interface {
	CreateTransaction(*LedgerRecord, *LedgerRecord, string) error
	FindLastLedgerRecordByAccountID(*LedgerRecord) error
	FindCurrencyBySymbol(*Currency) error
}

func NewTransactionsRepository(db *gorm.DB, cf config.Specification) TransactionsRepository {
	return transactionsRepository{
		db: db,
	}
}

type transactionsRepository struct {
	db *gorm.DB
	cf *config.Specification
}

func (tr transactionsRepository) executeTransactionStep(sqlTransaction *gorm.DB, ledgerRecord *LedgerRecord, txID string) (err error) {
	account := &Account{ID: ledgerRecord.AccountID}
	if err := sqlTransaction.Where("id = ?", ledgerRecord.AccountID).Last(account).Error; err != nil {
		return err
	}
	ledgerRecord.TransactionID = txID
	ledgerRecord.PreviousBanlance = account.Balance
	prevBalance, err := decimal.NewFromString(account.Balance)
	if err != nil {
		return err
	}
	var amountS string
	if ledgerRecord.Type == LedgerRecordCredit {
		amountS = ledgerRecord.Amount
	} else {
		amountS = fmt.Sprintf("-%s", ledgerRecord.Amount)
	}
	amount, err := decimal.NewFromString(amountS)
	if err != nil {
		return err
	}
	ledgerRecord.Balance = prevBalance.Add(amount).StringFixedBank(4)
	if err := sqlTransaction.Create(ledgerRecord).Error; err != nil {
		return err
	}
	if err := sqlTransaction.Model(account).Update("balance", ledgerRecord.Balance).Error; err != nil {
		return err
	}
	return nil
}

func onFail(sqlTransaction *gorm.DB, d time.Duration) {
	sqlTransaction.Rollback()
	time.Sleep(d)
}

func (tr transactionsRepository) CreateTransaction(ledgerRecordOrigin *LedgerRecord, ledgerRecordTarget *LedgerRecord, concept string) (err error) {
	entropy := rand.Reader
	tx := &Transaction{
		ID:      ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String(),
		Concept: concept,
	}

	err = try.Do(func(attempt int) (bool, error) {
		sqlTransaction := tr.db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
		defer func() {
			if r := recover(); r != nil {
				sqlTransaction.Rollback()
			}
		}()

		if err := sqlTransaction.Create(tx).Error; err != nil {
			onFail(sqlTransaction, time.Duration(mrand.Int31n(100))*time.Millisecond)
			return true, err
		}
		for _, ledgerRecord := range []*LedgerRecord{ledgerRecordOrigin, ledgerRecordTarget} {
			if err := tr.executeTransactionStep(sqlTransaction, ledgerRecord, tx.ID); err != nil {
				onFail(sqlTransaction, time.Duration(mrand.Int31n(100))*time.Millisecond)
				return true, err
			}
		}

		if err := sqlTransaction.Commit().Error; err != nil {
			onFail(sqlTransaction, time.Duration(mrand.Int31n(100))*time.Millisecond)
			return true, err
		}

		return false, nil
	})

	return err
}

func (tr transactionsRepository) FindLastLedgerRecordByAccountID(lr *LedgerRecord) (err error) {
	if err := tr.db.Where("account_id = ?", lr.AccountID).Last(lr).Error; err != nil {
		return err
	}
	return nil
}

func (tr transactionsRepository) FindCurrencyBySymbol(c *Currency) (err error) {
	if err := tr.db.Where("symbol = ?", c.Symbol).First(c).Error; err != nil {
		return err
	}
	return nil
}
