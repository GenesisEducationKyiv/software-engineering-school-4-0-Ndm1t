package database

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"subscription-service/internal/models"
	"testing"
)

func setupSQLiteDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&models.Email{})
	require.NoError(t, err)

	return db
}

func TestSubscriptionDao(t *testing.T) {
	db := setupSQLiteDB(t)

	repository := NewSubscriptionRepository(db)
	t.Run("create new subscription record success", func(t *testing.T) {
		email := "test@example.com"
		createdSubscription, err := repository.Create(email)
		require.NoError(t, err)
		assert.Equal(t, email, createdSubscription.Email)
		assert.Equal(t, models.Subscribed, createdSubscription.Status)
	})
	t.Run("find subscription success", func(t *testing.T) {
		email := "test1@example.com"
		_, err := repository.Create(email)
		require.NoError(t, err)

		subscription, err := repository.Find(email)
		require.NoError(t, err)
		assert.Equal(t, email, subscription.Email)
		assert.Equal(t, models.Subscribed, subscription.Status)
	})
	t.Run("list subscribed success", func(t *testing.T) {
		emails := []string{"test2@example.com", "test3@example.com"}
		for _, email := range emails {
			_, err := repository.Create(email)
			require.NoError(t, err)
		}

		subscriptions, err := repository.ListSubscribed()
		require.NoError(t, err)
		assert.Len(t, subscriptions, 4)
	})
	t.Run("update subscription success", func(t *testing.T) {
		email := "test@example.com"
		foundSubscription, err := repository.Find(email)
		require.NoError(t, err)

		foundSubscription.Status = models.Unsubscribed
		updatedSubscription, err := repository.Update(*foundSubscription)
		require.NoError(t, err)
		assert.Equal(t, models.Unsubscribed, updatedSubscription.Status)
		foundSubscription, err = repository.Find(email)
		require.NoError(t, err)
		assert.Equal(t, updatedSubscription.Status, foundSubscription.Status)
		assert.Equal(t, updatedSubscription.Email, foundSubscription.Email)
		assert.Equal(t, updatedSubscription.CreatedAt, foundSubscription.CreatedAt)
		assert.Equal(t, updatedSubscription.DeletedAt, foundSubscription.DeletedAt)

	})

	_ = db.Migrator().DropTable(&models.Email{})
}
