package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"subscription-service/internal/app_errors"
	"subscription-service/internal/models"
	"subscription-service/internal/services"
)

type SubscriptionController struct {
	SubscriptionService services.ISubscriptionService
}

type ISubscriptionController interface {
	Subscribe(ctx *gin.Context)
	ListSubscribed(ctx *gin.Context)
}

func NewSubscriptionController(
	subscriptionService services.ISubscriptionService) *SubscriptionController {
	return &SubscriptionController{
		SubscriptionService: subscriptionService,
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

func (c *SubscriptionController) ListSubscribed(ctx *gin.Context) {
	subscriptions, err := c.SubscriptionService.ListSubscribed()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"subscriptions": subscriptions,
	})
}
