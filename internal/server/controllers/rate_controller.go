package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gses4_project/internal/apperrors"
	"gses4_project/internal/container"
	"gses4_project/internal/services"
	"net/http"
)

type RateController struct {
	RateService services.IRateService
	container   container.IContainer
}

type IRateController interface {
	Get(ctx *gin.Context)
}

func NewRateController(container container.IContainer, rateService services.IRateService) *RateController {
	return &RateController{
		RateService: rateService,
		container:   container,
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
