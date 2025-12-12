package api

import (
	db "Bankstore/db/sqlc"

	"github.com/gin-gonic/gin"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
}

func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	// TODO: add routes to router
	router.POST("/accounts", server.CreateAccount)
	router.GET("/accounts/:id", server.GetAccount)
	router.GET("/accounts", server.ListAccounts)

	server.router = router
	return server
}

type CreateAccountRequest struct {
	Owner    string      `json:"owner" binding:"required"`
	Currency db.Currency `json:"currency" binding:"required,oneof=USD EUR"`
}

// errorResponse return gin.H -> map[string]interface{}
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

// Start server method
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
