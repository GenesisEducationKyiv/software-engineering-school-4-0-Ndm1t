package server

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
	"subscription-service/internal/server/controllers"
)

type Server struct {
	router                 *gin.Engine
	subscriptionController controllers.ISubscriptionController
}

func NewServer(
	subscriptionController controllers.ISubscriptionController) *Server {
	server := &Server{
		router:                 gin.Default(),
		subscriptionController: subscriptionController,
	}

	server.routes()
	return server
}

func (s *Server) routes() {
	api := s.router.Group("/api")
	{
		api.POST("/subscribe", s.subscriptionController.Subscribe)
		api.GET("/subscriptions", s.subscriptionController.ListSubscribed)
		api.POST("/unsubscribe", s.subscriptionController.Unsubscribe)
	}
}

func (s *Server) Run() {

	if err := s.router.Run(viper.GetString("PORT")); err != nil {
		log.Fatalf("Failed to run server %v", err.Error())
	}
}
