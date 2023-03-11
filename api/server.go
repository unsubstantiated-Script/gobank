package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/techschool/simplebank/db/sqlc"
)

// Serves up all HTTP requests for our banking service
type Server struct {
	store  *db.Store
	router *gin.Engine
}

//NewServer creates a new HTTP server and setup routing.
func NewServer(store *db.Store) *Server {
	//Setting things up
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/accounts", server.createAccount)
	// add routes to router

	server.router = router
	return server
}

// Start runs the HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
