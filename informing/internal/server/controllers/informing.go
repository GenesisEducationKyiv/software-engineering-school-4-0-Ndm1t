package controllers

import (
	"github.com/gin-gonic/gin"
	"informing-service/internal/services"
	"net/http"
)

type InformingController struct {
	InformingService services.InformingServiceInterface
}

type InformingControllerInterface interface {
	SendInforming(ctx *gin.Context)
}

func NewInformingController(
	informingService services.InformingServiceInterface) *InformingController {
	return &InformingController{
		InformingService: informingService,
	}
}

func (c *InformingController) SendInforming(ctx *gin.Context) {
	c.InformingService.SendEmails()

	ctx.Status(http.StatusOK)
	return
}
