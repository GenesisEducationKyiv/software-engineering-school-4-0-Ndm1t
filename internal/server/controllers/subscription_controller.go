package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gses4_project/internal/apperrors"
	"gses4_project/internal/container"
	"gses4_project/internal/models"
	"gses4_project/internal/services"
	"net/http"
)

type SubscriptionController struct {
	SubscriptionService services.ISubscriptionService
	container           container.IContainer
}

type ISubscriptionController interface {
	Subscribe(ctx *gin.Context)
}

func NewSubscriptionController(container container.IContainer,
	subscriptionService services.ISubscriptionService) *SubscriptionController {
	return &SubscriptionController{
		SubscriptionService: subscriptionService,
		container:           container,
	}
}

func (c *SubscriptionController) Subscribe(ctx *gin.Context) {
	var subscriptionData models.Email
	if err := ctx.ShouldBindJSON(&subscriptionData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Couldn't bind request",
		})
		return
	}

	subscription, err := c.SubscriptionService.Subscribe(subscriptionData.Email)

	if err != nil {
		var httpErr *apperrors.HttpError
		if errors.As(err, &httpErr) {
			ctx.JSON(httpErr.StatusCode, httpErr.JSONResponse)
			return
		}

		ctx.JSON(apperrors.ErrInternalServer.StatusCode, apperrors.ErrInternalServer.JSONResponse)
		return
	}
	ctx.JSON(http.StatusOK, &subscription)
}
