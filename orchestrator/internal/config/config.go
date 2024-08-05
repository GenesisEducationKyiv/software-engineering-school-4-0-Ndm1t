package config

import (
	"fmt"
	"github.com/spf13/viper"
)

func LoadConfig(path string) error {
	viper.SetConfigFile(path)
	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("error loading config file : %v", err)
	}
	return err
}
