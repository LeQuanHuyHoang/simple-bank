package api

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	db "simple-bank/db/sqlc"
)

type transferRequest struct {
	FromAccountID string `json:"from_account_id" binding:"required"`
	ToAccountID   string `json:"to_account_id" binding:"required"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (s *Server) createTransfer(ctx *gin.Context) {
	var req transferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	if !s.validAccount(ctx, uuid.MustParse(req.FromAccountID), req.Currency) {
		return
	}
	if !s.validAccount(ctx, uuid.MustParse(req.ToAccountID), req.Currency) {
		return
	}

	args := db.TransferTxParams{
		FromAccountID: uuid.MustParse(req.FromAccountID),
		ToAccountID:   uuid.MustParse(req.ToAccountID),
		Amount:        req.Amount,
	}

	account, err := s.store.TransferTx(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

func (s *Server) validAccount(ctx *gin.Context, accountID uuid.UUID, currency string) bool {
	account, err := s.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return false
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return false
	}

	if account.Currency != currency {
		err := fmt.Errorf("accout [%d] currency mismatch: %s vs %s", account.ID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return false
	}
	return true

}
