package config

import (
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadConfig_Success(t *testing.T) {
	err := LoadConfig("../../.env")
	assert.NoError(t, err)
	assert.NotEqual(t, viper.GetString("PORT"), "")
}

func TestLoadConfig_Failed(t *testing.T) {
	err := LoadConfig(".env")
	assert.Error(t, err)
}
