package server

import (
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"rate-service/internal/crons"
	"rate-service/internal/server/controllers"
)

type Server struct {
	rateController controllers.IRateController
	router         *gin.Engine
	Scheduler      crons.ICronScheduler
}

func NewServer(rateController controllers.IRateController, scheduler crons.ICronScheduler) *Server {
	server := &Server{
		router:         gin.Default(),
		rateController: rateController,
		Scheduler:      scheduler,
	}

	server.Scheduler.Setup()
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
