package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gses4_project/internal/apperrors"
	"gses4_project/internal/services"
	"net/http"
)

type IRateService interface {
	Get() (*float64, error)
}

type RateController struct {
	RateService IRateService
}

func NewRateController() *RateController {
	return &RateController{
		RateService: services.NewRateService(),
	}
}

func (c *RateController) Get(ctx *gin.Context) {
	rate, err := c.RateService.Get()
	if err != nil {
		var httpErr *apperrors.HttpError
		if errors.As(err, &httpErr) {
			ctx.JSON(httpErr.StatusCode, httpErr.Message)
			return
		}
		ctx.JSON(apperrors.ErrInternalServer.StatusCode, apperrors.ErrInternalServer.Message)
		return
	}
	ctx.JSON(http.StatusOK, rate)
}
