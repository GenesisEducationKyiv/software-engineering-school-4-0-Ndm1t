package database

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gses4_project/internal/models"
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

	dao := NewSubcscriptionDao(db)
	t.Run("create new subscription record success", func(t *testing.T) {
		email := "test@example.com"
		createdSubscription, err := dao.Create(email)
		require.NoError(t, err)
		assert.Equal(t, email, createdSubscription.Email)
		assert.Equal(t, models.Subscribed, createdSubscription.Status)
	})
	t.Run("find subscription success", func(t *testing.T) {
		email := "test1@example.com"
		_, err := dao.Create(email)
		require.NoError(t, err)

		subscription, err := dao.Find(email)
		require.NoError(t, err)
		assert.Equal(t, email, subscription.Email)
		assert.Equal(t, models.Subscribed, subscription.Status)
	})
	t.Run("list subscribed success", func(t *testing.T) {
		emails := []string{"test2@example.com", "test3@example.com"}
		for _, email := range emails {
			_, err := dao.Create(email)
			require.NoError(t, err)
		}

		subscriptions, err := dao.ListSubscribed()
		require.NoError(t, err)
		assert.Len(t, subscriptions, 4)
	})
	t.Run("update subscription success", func(t *testing.T) {
		email := "test@example.com"
		foundSubscription, err := dao.Find(email)
		require.NoError(t, err)

		foundSubscription.Status = models.Unsubscribed
		updatedSubscription, err := dao.Update(*foundSubscription)
		require.NoError(t, err)
		assert.Equal(t, models.Unsubscribed, updatedSubscription.Status)
		foundSubscription, err = dao.Find(email)
		require.NoError(t, err)
		assert.Equal(t, updatedSubscription.Status, foundSubscription.Status)
		assert.Equal(t, updatedSubscription.Email, foundSubscription.Email)
		assert.Equal(t, updatedSubscription.CreatedAt, foundSubscription.CreatedAt)
		assert.Equal(t, updatedSubscription.DeletedAt, foundSubscription.DeletedAt)

	})

	_ = db.Migrator().DropTable(&models.Email{})
}
