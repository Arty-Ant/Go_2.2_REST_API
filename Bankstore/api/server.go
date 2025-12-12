package api
import (
    db "Bankstore/db/sqlc"
    "github.com/gin-gonic/gin"
)
type Server struct {
    store *db.Store
    router *gin.Engine
}
func NewServer(store *db.Store) *Server {
    server := &Server{store: store}
    router := gin.Default()
    // accounts routes
    router.POST("/accounts", server.CreateAccount)
    router.GET("/accounts/:id", server.GetAccount)
    router.GET("/accounts", server.ListAccounts)
    router.DELETE("/accounts/:id", server.DeleteAccount)
    
    // users routes
    router.POST("/users", server.CreateUser)
    server.router = router
    return server
}
// errorResponse return gin.H -> map[string]interface{}
func errorResponse(err error) gin.H {
    return gin.H{"error": err.Error()}
}
// Start server method
func (server *Server) Start(address string) error {
    return server.router.Run(address)
}