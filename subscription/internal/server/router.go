package server

import (
	"fmt"
	"github.com/VictoriaMetrics/metrics"
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
	s.router.Use(s.metricsMiddleware())

	api := s.router.Group("/api")
	{
		api.POST("/subscribe", s.subscriptionController.Subscribe)
		api.GET("/subscriptions", s.subscriptionController.ListSubscribed)
		api.POST("/unsubscribe", s.subscriptionController.Unsubscribe)
		api.GET("/metrics", func(ctx *gin.Context) {
			metrics.WritePrometheus(ctx.Writer, true)
			return
		})
	}

}

func (s *Server) Run() {

	if err := s.router.Run(viper.GetString("PORT")); err != nil {
		log.Fatalf("Failed to run server %v", err.Error())
	}
}

func (s *Server) metricsMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		s := fmt.Sprintf(`requests_total{path=%q}`, ctx.Request.URL.Path)
		metrics.GetOrCreateCounter(s).Inc()
		ctx.Next()
		statusCode := ctx.Writer.Status()
		if statusCode >= 400 {
			s = fmt.Sprintf(`request_failed{path=%q}`, ctx.Request.URL.Path)
			metrics.GetOrCreateCounter(s).Inc()
		} else {
			s = fmt.Sprintf(`request_success{path=%q}`, ctx.Request.URL.Path)
			metrics.GetOrCreateCounter(s).Inc()
		}
	}
}
