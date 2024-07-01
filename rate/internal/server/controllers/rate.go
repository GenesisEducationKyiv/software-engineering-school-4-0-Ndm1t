package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"rate-service/internal/app_errors"
	"rate-service/internal/services"
)

type RateController struct {
	RateService services.IRateService
}

type IRateController interface {
	Get(ctx *gin.Context)
}

func NewRateController(rateService services.IRateService) *RateController {
	return &RateController{
		RateService: rateService,
	}
}

func (c *RateController) Get(ctx *gin.Context) {
	rate, err := c.RateService.Get()
	if err != nil {
		var httpErr *apperrors.HttpError
		if errors.As(err, &httpErr) {
			ctx.JSON(httpErr.StatusCode, httpErr.JSONResponse)
			return
		}
		ctx.JSON(apperrors.ErrInternalServer.StatusCode, apperrors.ErrInternalServer.JSONResponse)
		return
	}
	ctx.JSON(http.StatusOK, rate)
}
