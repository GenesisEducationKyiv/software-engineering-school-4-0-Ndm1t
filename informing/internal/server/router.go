package server

import (
	"fmt"
	"github.com/VictoriaMetrics/metrics"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"informing-service/internal/crons"
	"informing-service/internal/server/controllers"
	"log"
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
	s.router.Use(s.metricsMiddleware())
	api := s.router.Group("/api")
	{
		api.POST("/inform", s.informingController.SendInforming)
		api.GET("/metrics", func(ctx *gin.Context) {
			metrics.WritePrometheus(ctx.Writer, false)
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
