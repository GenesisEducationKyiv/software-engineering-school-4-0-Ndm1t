package server

import (
	"github.com/gin-gonic/gin"
	"informing-service/internal/crons"
	"informing-service/internal/server/controllers"
	"log"
	"os"
)

type Server struct {
	informingController controllers.InformingControllerInterface
	router              *gin.Engine
	Scheduler           crons.ICronScheduler
}

func NewServer(
	informingController controllers.InformingControllerInterface, cronScheduler crons.ICronScheduler) *Server {
	server := &Server{
		informingController: informingController,
		router:              gin.Default(),
		Scheduler:           cronScheduler,
	}

	server.Scheduler.Setup()
	server.routes()
	return server
}

func (s *Server) routes() {
	api := s.router.Group("/api")
	{
		api.POST("/inform", s.informingController.SendInforming)
	}
}

func (s *Server) Run() {

	if err := s.router.Run(os.Getenv("PORT")); err != nil {
		log.Fatalf("Failed to run server %v", err.Error())
	}
}
