package handlers

import (
	"net/http"

	"github.com/labstack/echo"
	"nquicenob.com/golang-api-example/internal/services"
)

type TransactionHandler interface {
	Create(c echo.Context) (err error)
	Find(c echo.Context) (err error)
}

func NewTransactionHandler(ts services.TransactionsService) TransactionHandler {
	return transactionHandler{
		transactionsService: ts,
	}
}

type transactionHandler struct {
	transactionsService services.TransactionsService
}

type TransactionInput struct {
	Data *services.Transaction `json:"data"`
}

type TransactionOutput struct {
	Data *services.TransactionResult `json:"data"`
}

type AccountBalance struct {
	Data *services.AccountBalance `json:"data"`
}

func (handler transactionHandler) Create(c echo.Context) (err error) {
	originAccountID := c.Param("account_id")
	ti := new(TransactionInput)
	if err = c.Bind(ti); err != nil {
		return
	}

	result, err := handler.transactionsService.CreateTransaction(originAccountID, ti.Data)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "liada parda")
	}
	return c.JSON(
		http.StatusCreated,
		TransactionOutput{Data: result},
	)
}

func (handler transactionHandler) Find(c echo.Context) (err error) {
	originAccountID := c.Param("account_id")
	result, err := handler.transactionsService.GetAccountAndBalance(originAccountID)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "liada parda")
	}
	return c.JSON(
		http.StatusCreated,
		AccountBalance{Data: result},
	)
}
