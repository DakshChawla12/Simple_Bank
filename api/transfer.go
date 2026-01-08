package api

import (
	"database/sql"
	"errors"
	"net/http"

	db "github.com/DakshChawla/simplebank/db/sqlc"
	"github.com/DakshChawla/simplebank/token"
	"github.com/gin-gonic/gin"
)

type transferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) validateAccount(c *gin.Context, accountID int64, currency string) (db.Account, bool) {
	account, err := server.store.GetAccount(c, accountID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return account, false
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return account, false
	}

	if account.Currency != currency {
		err = errors.New("invalid currency")
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return account, false
	}

	return account, true
}

func (server *Server) createTransfer(c *gin.Context) {
	var req transferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	fromAccount, valid := server.validateAccount(c, req.FromAccountID, req.Currency)
	if !valid {
		return
	}

	_, valid = server.validateAccount(c, req.ToAccountID, req.Currency)
	if !valid {
		return
	}

	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Username != fromAccount.Owner {
		c.JSON(http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	result, err := server.store.TransferTx(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, result)
}
