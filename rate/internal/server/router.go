package server

import (
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"rate-service/internal/server/controllers"
)

type Server struct {
	rateController controllers.IRateController
	router         *gin.Engine
}

func NewServer(rateController controllers.IRateController) *Server {
	server := &Server{
		router:         gin.Default(),
		rateController: rateController,
	}
	server.routes()
	return server
}

func (s *Server) routes() {
	api := s.router.Group("/api")
	{
		api.GET("/rate", s.rateController.Get)
	}
}

func (s *Server) Run() {

	if err := s.router.Run(os.Getenv("PORT")); err != nil {
		log.Fatalf("Failed to run server %v", err.Error())
	}
}
