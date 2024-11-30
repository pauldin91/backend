package api

import (
	db "backend/db/sqlc"

	"github.com/gin-gonic/gin"
)

type Server struct {
	//config utils.Config
	store db.Store
	//tokenMaker token.Maker
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.GET("/accounts", server.listAccounts)
	router.GET("/accounts/:id", server.getAccount)
	router.POST("/accounts", server.createAccount)
	router.POST("/transfers", server.createTransfer)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
