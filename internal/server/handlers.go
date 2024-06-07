package server

import (
	"github.com/gin-gonic/gin"
	"gses4_project/internal/database"
	"gses4_project/internal/models"
	"gses4_project/internal/pkg"
	"net/http"
)

func (s *Server) GetRate(c *gin.Context) {
	rate, err := pkg.FetchRate()
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, rate)
}

func (s *Server) Subscribe(c *gin.Context) {
	var subscriptionData models.Email
	if err := c.ShouldBindJSON(&subscriptionData); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	subscription := models.Email{Email: subscriptionData.Email, Status: models.Subscribed}
	result := database.DB.FirstOrCreate(&subscription, models.Email{Email: subscriptionData.Email})
	if result.Error != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if result.RowsAffected == 0 && subscription.Status == models.Subscribed {
		c.Status(http.StatusBadRequest)
		return
	}

	if result.RowsAffected == 0 && subscription.Status == models.Unsubscribed {
		subscription.Status = models.Subscribed
		subscription.DeletedAt = nil
		result = database.DB.Save(&subscription)
	}

	c.Status(http.StatusOK)
}
