package server

import (
	"fmt"
	"github.com/VictoriaMetrics/metrics"
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
	s.router.Use(s.metricsMiddleware())
	api := s.router.Group("/api")
	{
		api.GET("/rate", s.rateController.Get)
		api.GET("/metrics", func(ctx *gin.Context) {
			metrics.WritePrometheus(ctx.Writer, true)
			return
		})
	}
}

func (s *Server) Run() {

	if err := s.router.Run(os.Getenv("PORT")); err != nil {
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
