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
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to fetch currency data",
		})
		return
	}
	c.JSON(http.StatusOK, rate)
}

func (s *Server) Subscribe(c *gin.Context) {
	var subscriptionData models.Email
	if err := c.ShouldBindJSON(&subscriptionData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "incorrect request format",
		})
		return
	}

	subscription := models.Email{Email: subscriptionData.Email, Status: models.Subscribed}
	result := database.DB.FirstOrCreate(&subscription, models.Email{Email: subscriptionData.Email})

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "database error",
		})
		return
	}

	if result.RowsAffected == 0 && subscription.Status == models.Subscribed {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "already subscribed",
		})
		return
	}

	if result.RowsAffected == 0 && subscription.Status == models.Unsubscribed {
		subscription.Status = models.Subscribed
		subscription.DeletedAt = nil
		database.DB.Save(&subscription)
	}

	c.Status(http.StatusOK)
}
