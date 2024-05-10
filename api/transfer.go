package api

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"net/http"
	db "simple-bank/db/sqlc"
	"simple-bank/token"
)

type transferRequest struct {
	FromAccountID string `json:"from_account_id" binding:"required"`
	ToAccountID   string `json:"to_account_id" binding:"required"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req transferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	fromAccount, valid := server.validAccount(ctx, uuid.MustParse(req.FromAccountID), req.Currency)
	if !valid {
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if fromAccount.Owner != authPayload.Username {
		err := errors.New("fromAccount doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errResponse(err))
		return
	}

	_, valid = server.validAccount(ctx, uuid.MustParse(req.ToAccountID), req.Currency)
	if !valid {
		return
	}

	args := db.TransferTxParams{
		FromAccountID: uuid.MustParse(req.FromAccountID),
		ToAccountID:   uuid.MustParse(req.ToAccountID),
		Amount:        req.Amount,
	}

	account, err := server.store.TransferTx(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

func (server *Server) validAccount(ctx *gin.Context, accountID uuid.UUID, currency string) (db.Account, bool) {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return account, false
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return account, false
	}

	if account.Currency != currency {
		err := fmt.Errorf("accout [%d] currency mismatch: %s vs %s", account.ID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return account, false
	}
	return account, true

}
