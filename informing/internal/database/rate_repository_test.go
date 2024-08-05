package database

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"informing-service/internal/models"
	"testing"
	"time"
)

func TestRateRepository(t *testing.T) {
	db := SetupSQLiteDB(t)

	repository := NewRateRepository(db)
	t.Run("create new rate record success", func(t *testing.T) {
		createdAt := time.Now()
		rate := models.Rate{
			Rate:      42.00,
			CreatedAt: createdAt,
		}
		createdRate, err := repository.Create(rate)
		require.NoError(t, err)
		assert.Equal(t, 42.00, createdRate.Rate)
		assert.NotNil(t, createdRate.ID)
		assert.Equal(t, createdAt, createdRate.CreatedAt)
	})
	t.Run("get latest rate success", func(t *testing.T) {
		createdRate, err := repository.GetLatest()
		require.NoError(t, err)
		assert.Equal(t, 42.00, createdRate.Rate)
		assert.NotNil(t, createdRate.ID)
		assert.NotNil(t, createdRate.CreatedAt)
	})

	_ = db.Migrator().DropTable(&models.Subscription{})
}
