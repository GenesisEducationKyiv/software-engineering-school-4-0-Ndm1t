package server

import (
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

type Server struct {
	router *gin.Engine
}

func NewServer() *Server {
	server := &Server{
		router: gin.Default(),
	}

	server.routes()
	return server
}

func (s *Server) routes() {
	api := s.router.Group("/api")
	{
		api.GET("/rate", s.GetRate)
		api.POST("/subscribe", s.Subscribe)
	}
}

func (s *Server) Run() {
	if err := s.router.Run(os.Getenv("PORT")); err != nil {
		log.Fatalf("Failed to run server %v", err.Error())
	}
}
