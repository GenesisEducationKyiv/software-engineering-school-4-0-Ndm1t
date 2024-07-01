package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type (
	RateClientInterface interface {
		FetchRate() (*int, *string, []byte, error)
	}

	RateControllerInterface interface {
		Get(ctx *gin.Context)
	}

	RateController struct {
		rateClient RateClientInterface
	}
)

func NewRateController(rateProvider RateClientInterface) *RateController {
	return &RateController{
		rateClient: rateProvider,
	}
}

func (c *RateController) Get(ctx *gin.Context) {
	status, contentType, body, err := c.rateClient.FetchRate()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to access rate service",
			"cause": err,
		})
		return
	}

	ctx.Data(*status, *contentType, body)
	return
}
