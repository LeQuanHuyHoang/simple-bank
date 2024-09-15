package handler

import (
	"Go_Learn/pkg/model"
	"Go_Learn/pkg/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type UserHandler struct {
	UserSrv service.IUserService
}

type IUserHandler interface {
	Login(ctx *gin.Context)
	SearchFile(ctx *gin.Context)
}

func NewUserHandler(srv service.IUserService) IUserHandler {
	return &UserHandler{
		UserSrv: srv,
	}
}

func (h *UserHandler) Login(ctx *gin.Context) {
	req := model.LoginRequest{}
	ctx.ShouldBindJSON(&req)

	//Get JWT
	token, err := h.UserSrv.GenJWTToken(1)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, token)
}

func (h *UserHandler) SearchFile(ctx *gin.Context) {
	authHeader := ctx.Request.Header.Get("Authorization")
	authHeader = strings.TrimPrefix(authHeader, "Bearer ")
	_, err := h.UserSrv.ValidateToken(authHeader)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	req := model.FilterRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {

		fmt.Println("req", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	fmt.Println("req", err)

	rs, err := h.UserSrv.SearchFile(req.Filters)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, rs)
}
