package server

import (
	"gateway/internal/server/controllers"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

type (
	Server struct {
		rateController         controllers.RateControllerInterface
		subscriptionController controllers.SubscriptionControllerInterface
		router                 *gin.Engine
	}
)

func NewServer(rateController controllers.RateControllerInterface,
	subscriptionController controllers.SubscriptionControllerInterface) *Server {
	server := &Server{
		rateController:         rateController,
		subscriptionController: subscriptionController,
		router:                 gin.Default(),
	}
	server.routes()
	return server
}

func (s *Server) routes() {
	api := s.router.Group("/api")
	{
		api.POST("/subscribe", s.subscriptionController.Subscribe)
		api.GET("/rate", s.rateController.Get)
	}
}

func (s *Server) Run() {

	if err := s.router.Run(os.Getenv("PORT")); err != nil {
		log.Fatalf("Failed to run server %v", err.Error())
	}
}
