package server

import (
	"github.com/gin-gonic/gin"
	"gses4_project/internal/container"
	"gses4_project/internal/crons"
	"gses4_project/internal/server/controllers"
	"log"
	"os"
)

type Server struct {
	rateController         controllers.IRateController
	router                 *gin.Engine
	subscriptionController controllers.ISubscriptionController
	Scheduler              crons.ICronScheduler
	container              container.IContainer
}

func NewServer(container container.IContainer, rateController controllers.IRateController,
	subscriptionController controllers.ISubscriptionController, cronScheduler crons.ICronScheduler) *Server {
	server := &Server{
		router:                 gin.Default(),
		rateController:         rateController,
		subscriptionController: subscriptionController,
		Scheduler:              cronScheduler,
		container:              container,
	}

	server.Scheduler.Setup()
	server.routes()
	return server
}

func (s *Server) routes() {
	api := s.router.Group("/api")
	{
		api.GET("/rate", s.rateController.Get)
		api.POST("/subscribe", s.subscriptionController.Subscribe)
	}
}

func (s *Server) Run() {

	if err := s.router.Run(os.Getenv("PORT")); err != nil {
		log.Fatalf("Failed to run server %v", err.Error())
	}
}
