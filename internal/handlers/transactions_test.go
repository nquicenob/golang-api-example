package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"nquicenob.com/golang-api-example/internal/config"
	"nquicenob.com/golang-api-example/internal/server"
	"nquicenob.com/golang-api-example/internal/services"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"nquicenob.com/golang-api-example/internal/db"
	"nquicenob.com/golang-api-example/internal/handlers"
	"nquicenob.com/golang-api-example/internal/models"
)

var currencies = []*models.Currency{
	&models.Currency{
		// ID:     uint(1),
		Symbol: models.CurrencyTypeEUR,
		Name:   "Euro",
	},
}

var txs = []*models.Transaction{
	&models.Transaction{
		ID:      "11DNFR6T85VAKDMKG6WVQ5SBYZ",
		Concept: "account origination",
	},
	&models.Transaction{
		ID:      "21DNFR6T85VAKDMKG6WVQ5SBYZ",
		Concept: "account origination",
	},
	&models.Transaction{
		ID:      "31DNFR6T85VAKDMKG6WVQ5SBYZ",
		Concept: "account origination",
	},
}

var accounts = []*models.Account{
	&models.Account{
		ID:      "01DNFR6T85VAKDMKG6WVQ5SBYZ",
		Name:    "Nicolas Quiceno",
		Email:   "nquicenob@gmail.com",
		Balance: "0",
		LedgerRecords: []*models.LedgerRecord{
			&models.LedgerRecord{
				Type:             models.LedgerRecordCredit,
				Balance:          "0",
				Amount:           "0",
				PreviousBanlance: "0",
				CurrencyID:       1,
				TransactionID:    txs[0].ID,
			},
		},
	},
	&models.Account{
		ID:      "01DNFR6T85S36AAE7RTCPVZYAV",
		Name:    "Karen Mercado",
		Email:   "kr@gmail.com",
		Balance: "0",
		LedgerRecords: []*models.LedgerRecord{
			&models.LedgerRecord{
				Type:             models.LedgerRecordCredit,
				Balance:          "0",
				Amount:           "0",
				PreviousBanlance: "0",
				CurrencyID:       1,
				TransactionID:    txs[1].ID,
			},
		},
	},
	&models.Account{
		ID:      "01DNFR6T85C3N7ZA9BZ8KKWHKK",
		Name:    "Santiago Mercado",
		Email:   "santy3@gmail.com",
		Balance: "0",
		LedgerRecords: []*models.LedgerRecord{
			&models.LedgerRecord{
				Type:             models.LedgerRecordCredit,
				Balance:          "0",
				Amount:           "0",
				PreviousBanlance: "0",
				CurrencyID:       1,
				TransactionID:    txs[2].ID,
			},
		},
	},
	&models.Account{
		ID:      "41DNFR6T85C3N7ZA9BZ8KKWHKK",
		Name:    "Santiago Mercado",
		Email:   "santy4@gmail.com",
		Balance: "0",
		LedgerRecords: []*models.LedgerRecord{
			&models.LedgerRecord{
				Type:             models.LedgerRecordCredit,
				Balance:          "0",
				Amount:           "0",
				PreviousBanlance: "0",
				CurrencyID:       1,
				TransactionID:    txs[2].ID,
			},
		},
	},
}

var dbConn *gorm.DB

var cf config.Specification

func TestMain(m *testing.M) {
	cf, err := config.New()
	if err != nil {
		os.Exit(1)
	}
	dbConn, err = db.New(cf)
	if err != nil {
		os.Exit(1)
	}
	db.DropSchema(dbConn)
	db.AutoMigrate(dbConn)
	db.LoadCurrenciesData(dbConn, currencies)
	db.LoadTransactionsData(dbConn, txs)
	db.LoadAccountsData(dbConn, accounts)

	code := m.Run()

	// db.DropSchema(dbConn)
	dbConn.Close()
	os.Exit(code)
}

// 404
// 409 - no money

type transactionInput struct {
	Name   string
	Result map[string]string
	Params map[string]string
	Body   string
}

func createTransaction(ti *transactionInput, server *echo.Echo, h handlers.TransactionHandler) (rec *httptest.ResponseRecorder, err error) {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(ti.Body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()

	c := server.NewContext(req, rec)
	c.SetPath("/accounts/:account_id/movemoney")
	c.SetParamNames("account_id")
	c.SetParamValues(ti.Params["account_id"])

	return rec, h.Create(c)
}

func getAccountAndBalance(accountID string, server *echo.Echo, h handlers.TransactionHandler) (rec *httptest.ResponseRecorder, err error) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()

	c := server.NewContext(req, rec)
	c.SetPath("/accounts/:account_id")
	c.SetParamNames("account_id")
	c.SetParamValues(accountID)

	return rec, h.Find(c)
}

func TestCreateTransaction(t *testing.T) {
	server := server.New(cf)
	sr := models.NewTransactionsRepository(dbConn, cf)
	ar := models.NewAccountRepository(dbConn)
	st := services.NewTransactionsService(sr, ar)
	h := handlers.NewTransactionHandler(st)

	rec, err := createTransaction(
		&transactionInput{
			Params: map[string]string{"account_id": "01DNFR6T85C3N7ZA9BZ8KKWHKK"},
			Body:   `{"data":{"amount":{"value":"1.99","currency":"EUR"},"concept":"fiirst movement","target_account_id":"41DNFR6T85C3N7ZA9BZ8KKWHKK"}}`,
		},
		server,
		h,
	)

	// Assertions
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		result := new(handlers.TransactionOutput)
		err = json.NewDecoder(rec.Body).Decode(result)
		assert.NoError(t, err)
		assert.NotNil(t, result.Data.ID)
		assert.Equal(t, result.Data.PreviousBanlance.Value, "0.0000")
		assert.Equal(t, result.Data.Banlance.Value, "-1.9900")
	}
}

func TestCreateTransactionConcurrency(t *testing.T) {
	server := server.New(cf)
	sr := models.NewTransactionsRepository(dbConn, cf)
	ar := models.NewAccountRepository(dbConn)
	st := services.NewTransactionsService(sr, ar)
	h := handlers.NewTransactionHandler(st)

	t.Run("Concurrency balance", func(t *testing.T) {
		testCases := []*transactionInput{
			&transactionInput{
				Name:   "test 1",
				Params: map[string]string{"account_id": "01DNFR6T85VAKDMKG6WVQ5SBYZ"},
				Body:   `{"data":{"amount":{"value":"1","currency":"EUR"},"concept":"tx1 parallel","target_account_id":"01DNFR6T85S36AAE7RTCPVZYAV"}}`,
			},
			&transactionInput{
				Name:   "test 2",
				Params: map[string]string{"account_id": "01DNFR6T85VAKDMKG6WVQ5SBYZ"},
				Body:   `{"data":{"amount":{"value":"2","currency":"EUR"},"concept":"tx2 parallel","target_account_id":"01DNFR6T85S36AAE7RTCPVZYAV"}}`,
			},
			&transactionInput{
				Name:   "test 3",
				Params: map[string]string{"account_id": "01DNFR6T85S36AAE7RTCPVZYAV"},
				Body:   `{"data":{"amount":{"value":"1","currency":"EUR"},"concept":"tx3 parallel","target_account_id":"01DNFR6T85VAKDMKG6WVQ5SBYZ"}}`,
			},
			&transactionInput{
				Name:   "test 4",
				Params: map[string]string{"account_id": "01DNFR6T85S36AAE7RTCPVZYAV"},
				Body:   `{"data":{"amount":{"value":"2","currency":"EUR"},"concept":"tx4 parallel","target_account_id":"01DNFR6T85VAKDMKG6WVQ5SBYZ"}}`,
			},
		}
		for _, tc := range testCases {
			tc := tc
			t.Run(tc.Name, func(t *testing.T) {
				t.Parallel()
				rec, err := createTransaction(
					tc,
					server,
					h,
				)
				if assert.NoError(t, err) {
					assert.Equal(t, http.StatusCreated, rec.Code)
					result := new(handlers.TransactionOutput)
					err = json.NewDecoder(rec.Body).Decode(result)
					assert.NoError(t, err)
					assert.NotNil(t, result.Data.ID)
				}
			})
		}
	})

	t.Run("Get results", func(t *testing.T) {
		for _, accountID := range []string{"01DNFR6T85VAKDMKG6WVQ5SBYZ", "01DNFR6T85S36AAE7RTCPVZYAV"} {
			rec1, err := getAccountAndBalance(accountID, server, h)
			if assert.NoError(t, err) {
				result := &handlers.AccountBalance{}
				err := json.Unmarshal([]byte(rec1.Body.String()), &result)
				assert.NoError(t, err)
				assert.Equal(
					t,
					"0.0000",
					result.Data.Balance.Value,
				)
			}
		}
	})

}
