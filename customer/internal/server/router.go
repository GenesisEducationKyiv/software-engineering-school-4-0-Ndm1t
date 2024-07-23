package server

import (
	"github.com/VictoriaMetrics/metrics"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
)

type Server struct {
	router *gin.Engine
}

func NewServer() *Server {
	server := &Server{
		router: gin.Default(),
	}

	server.routes()
	return server
}

func (s *Server) routes() {
	api := s.router.Group("/api")
	{
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
