package api

import (
	"github.com/gin-gonic/gin"
	db "simple-bank/db/sqlc"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
}

func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/accounts", server.createAccount)
	router.GET("/list-accounts", server.listAccount)
	router.GET("/get-account/:id", server.getAccount)
	router.PATCH("/update-account", server.addAccountBalance)
	router.DELETE("/delete-account", server.deleteAccount)

	server.router = router
	return server
}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func errResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
