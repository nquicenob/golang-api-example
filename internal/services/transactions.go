package services

import (
	"nquicenob.com/golang-api-example/internal/models"
)

type TransactionsService interface {
	CreateTransaction(originAccountID string, t *Transaction) (r *TransactionResult, err error)
	GetAccountAndBalance(originAccountID string) (b *AccountBalance, err error)
}

func NewTransactionsService(tr models.TransactionsRepository, ac models.AccountRepository) TransactionsService {
	return transactionsService{
		transactionsRepo: tr,
		accountRepo:      ac,
	}
}

type transactionsService struct {
	transactionsRepo models.TransactionsRepository
	accountRepo      models.AccountRepository
}

type Amount struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
	Type     string `json:"-"`
}

type Transaction struct {
	Amount          Amount `json:"amount"`
	Concept         string `json:"concept"`
	TargetAccountID string `json:"target_account_id"`
}

type TransactionResult struct {
	ID               uint   `json:"id"`
	Amount           Amount `json:"amount"`
	Concept          string `json:"concept"`
	TargetAccountID  string `json:"target_account_id"`
	PreviousBanlance Amount `json:"previous_banlance"`
	Banlance         Amount `json:"banlance"`
}

type Owner struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type AccountBalance struct {
	AccountNumber string `json:"account_number"`
	Owner         Owner  `json:"owner"`
	Balance       Amount `json:"balance"`
}

func (ts transactionsService) CreateTransaction(originAccountID string, t *Transaction) (r *TransactionResult, err error) {

	t.Amount.Type = string(models.LedgerRecordDebit)
	accountOrigin := &models.Account{ID: originAccountID}
	if err := ts.accountRepo.FindAccountID(accountOrigin); err != nil {
		return r, err
	}
	accountTarget := &models.Account{ID: t.TargetAccountID}
	if err := ts.accountRepo.FindAccountID(accountTarget); err != nil {
		return r, err
	}

	currency := &models.Currency{Symbol: models.CurrencyType(t.Amount.Currency)}
	if err := ts.transactionsRepo.FindCurrencyBySymbol(currency); err != nil {
		return r, err
	}

	ledgerRecordOrigin := &models.LedgerRecord{
		Type:       models.LedgerRecordType(t.Amount.Type),
		Amount:     t.Amount.Value,
		CurrencyID: currency.ID,
		AccountID:  originAccountID,
	}

	targetTypeOP := models.LedgerRecordCredit
	if models.LedgerRecordType(t.Amount.Type) == targetTypeOP {
		targetTypeOP = models.LedgerRecordDebit
	}

	ledgerRecordTarget := &models.LedgerRecord{
		Type:       targetTypeOP,
		Amount:     t.Amount.Value,
		CurrencyID: currency.ID,
		AccountID:  t.TargetAccountID,
	}

	if err := ts.transactionsRepo.CreateTransaction(ledgerRecordOrigin, ledgerRecordTarget, t.Concept); err != nil {
		return r, err
	}

	r = &TransactionResult{
		ID:              ledgerRecordOrigin.ID,
		Amount:          t.Amount,
		Concept:         t.Concept,
		TargetAccountID: t.TargetAccountID,
		PreviousBanlance: Amount{
			Value:    ledgerRecordOrigin.PreviousBanlance,
			Currency: t.Amount.Currency,
		},
		Banlance: Amount{
			Value:    ledgerRecordOrigin.Balance,
			Currency: t.Amount.Currency,
		},
	}

	return r, err
}

func (ts transactionsService) GetAccountAndBalance(originAccountID string) (b *AccountBalance, err error) {
	account := &models.Account{ID: originAccountID}
	if err := ts.accountRepo.FindAccountID(account); err != nil {
		return b, err
	}

	return &AccountBalance{
		AccountNumber: account.ID,
		Owner: Owner{
			Name:  account.Name,
			Email: account.Email,
		},
		Balance: Amount{
			Value:    account.Balance,
			Currency: "EUR",
		},
	}, nil
}
