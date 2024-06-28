package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestLoadConfig_Success(t *testing.T) {
	err := LoadConfig("../../.env")
	assert.NoError(t, err)
	assert.NotEqual(t, os.Getenv("PORT"), "")
}

func TestLoadConfig_Failed(t *testing.T) {
	err := LoadConfig()
	assert.Error(t, err)
}
