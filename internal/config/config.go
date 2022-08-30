package config

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
)

func Init() error {
	viper.SetDefault("discord.token", "")
	viper.SetDefault("rcon.address", "")
	viper.SetDefault("rcon.password", "minecraft") // mc's default rcon password
	viper.AddConfigPath(".")
	viper.SetConfigName("warden_config")
	viper.SetEnvPrefix("WARDEN")
	viper.AutomaticEnv()

	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	if err := viper.ReadInConfig(); err != nil {
		if ok := errors.As(err, &viper.ConfigFileNotFoundError{}); !ok {
			return fmt.Errorf("viper failed to read in config: %w", err)
		}

		log.Println("Config file not found, using defaults and env")
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
