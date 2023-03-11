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

	// add routes to router

	server.router = router
	return server
}
