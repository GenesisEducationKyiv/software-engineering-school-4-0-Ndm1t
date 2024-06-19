package server

import (
	"github.com/gin-gonic/gin"
	"gses4_project/internal/crons"
	"gses4_project/internal/server/controllers"
	"log"
	"os"
)

type Server struct {
	rateController         *controllers.RateController
	router                 *gin.Engine
	subscriptionController *controllers.SubscriptionController
	Scheduler              *crons.CronScheduler
}

func NewServer() *Server {
	server := &Server{
		router:                 gin.Default(),
		rateController:         controllers.NewRateController(),
		subscriptionController: controllers.NewSubscriptionController(),
		Scheduler:              crons.NewCronScheduler(),
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
