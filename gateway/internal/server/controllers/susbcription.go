package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type (
	SubscribeReq struct {
		Email string `json:"email"`
	}

	SubscriptionProviderInterface interface {
		Subscribe(req SubscribeReq) (*int, *string, []byte, error)
		Unsubscribe(req SubscribeReq) (*int, *string, []byte, error)
	}

	SubscriptionControllerInterface interface {
		Subscribe(ctx *gin.Context)
		Unsubscribe(ctx *gin.Context)
	}

	SubscriptionController struct {
		subscriptionProvider SubscriptionProviderInterface
	}
)

func NewSubscriptionController(subscriptionProvider SubscriptionProviderInterface) *SubscriptionController {
	return &SubscriptionController{
		subscriptionProvider: subscriptionProvider,
	}
}

func (c *SubscriptionController) Subscribe(ctx *gin.Context) {
	var req SubscribeReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to serialize json",
			"cause": err,
		})
		return
	}
	status, contentType, body, err := c.subscriptionProvider.Subscribe(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to access service",
			"cause": err,
		})
		return
	}
	ctx.Data(*status, *contentType, body)
	return

}

func (c *SubscriptionController) Unsubscribe(ctx *gin.Context) {
	var req SubscribeReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to serialize json",
			"cause": err,
		})
		return
	}
	status, contentType, body, err := c.subscriptionProvider.Unsubscribe(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to access service",
			"cause": err,
		})
		return
	}
	ctx.Data(*status, *contentType, body)
	return

}
