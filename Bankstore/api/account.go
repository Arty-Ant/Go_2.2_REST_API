package api

import (
	db "Bankstore/db/sqlc"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (server *Server) CreateAccount(ctx *gin.Context) {
	var req CreateAccountRequest
	// десериализация входящего JSON'a
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// на основе CreateAccountRequest создаём CreateAccountParams
	arg := db.CreateAccountParams{
		Owner:    req.Owner,
		Currency: req.Currency,
		Balance:  0,
	}
	// на основе arg создаём новый аккаунт
	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, account)
}
