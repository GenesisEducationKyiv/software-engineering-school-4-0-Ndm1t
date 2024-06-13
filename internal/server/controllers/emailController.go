package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gses4_project/internal/apperrors"
	"gses4_project/internal/models"
	"gses4_project/internal/services"
	"net/http"
)

type ISubscriptionService interface {
	Subscribe(email string) error
}

type SubscriptionController struct {
	SubscriptionService ISubscriptionService
}

func NewSubscriptionController() *SubscriptionController {
	return &SubscriptionController{
		SubscriptionService: services.NewSubscriptionService(),
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

	err := c.SubscriptionService.Subscribe(subscriptionData.Email)

	if err != nil {
		var httpErr *apperrors.HttpError
		if errors.As(err, &httpErr) {
			ctx.JSON(httpErr.StatusCode, httpErr.Message)
			return
		}

		ctx.JSON(apperrors.ErrInternalServer.StatusCode, apperrors.ErrInternalServer.Message)
		return
	}
	ctx.Status(http.StatusOK)
}