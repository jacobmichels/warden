package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func Init() error {
	viper.SetDefault("discord.token", "")
	viper.AddConfigPath(".")
	viper.SetConfigName("warden_config")

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("viper failed to read in config: %w", err)
	}

	return nil
}

func GetStr(s string) string {
	return viper.GetString(s)
}

func GetInt(s string) int {
	return viper.GetInt(s)
}

func GetBool(s string) bool {
	return viper.GetBool(s)
}
