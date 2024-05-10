package api

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"net/http"
	db "simple-bank/db/sqlc"
	"simple-bank/token"
)

type createAccountRequest struct {
	Currency string `json:"currency" binding:"required,currency"`
}

type ListAccountsParams struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

type GetAccountParams struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type AddAccountBalanceParams struct {
	Amount int64  `json:"amount" binding:"required,min=1"`
	ID     string `json:"id" binding:"required"`
}

type DeleteAccountParams struct {
	ID string `json:"id" binding:"required,uuid"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	args := db.CreateAccountParams{
		Owner:    authPayload.Username,
		Currency: req.Currency,
		Balance:  0,
	}

	account, err := server.store.CreateAccount(ctx, args)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, errResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

func (server *Server) listAccount(ctx *gin.Context) {
	var req ListAccountsParams
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	args := db.ListAccountsParams{
		Owner:  authPayload.Username,
		Limit:  req.PageSize,
		Offset: req.PageID,
	}

	account, err := server.store.ListAccounts(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

func (server *Server) getAccount(ctx *gin.Context) {
	var req GetAccountParams
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}
	account, err := server.store.GetAccount(ctx, uuid.MustParse(req.ID))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Username != account.Owner {
		err := errors.New("account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errResponse(err))
	}

	ctx.JSON(http.StatusOK, account)
}

func (server *Server) addAccountBalance(ctx *gin.Context) {
	var req AddAccountBalanceParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	args := db.AddAccountBalanceParams{
		Amount: req.Amount,
		ID:     uuid.MustParse(req.ID),
	}

	account, err := server.store.AddAccountBalance(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

func (server *Server) deleteAccount(ctx *gin.Context) {
	var req DeleteAccountParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	err := server.store.DeleteAccount(ctx, uuid.MustParse(req.ID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, nil)
}
