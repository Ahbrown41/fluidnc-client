package config

import (
	"fmt"
	"strings"

	"fluidnc-client/internal/fluidnc"
	"github.com/spf13/viper"
)

// LoadConfig loads configuration from file, environment variables, and defaults
func LoadConfig() (*fluidnc.Config, error) {
	viper.SetConfigName("fluidnc-cli")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.fluidnc-cli")
	viper.AddConfigPath("/etc/fluidnc-cli")

	// Set defaults
	viper.SetDefault("host", "192.168.1.100")
	viper.SetDefault("port", 80)
	viper.SetDefault("websocket_port", 81)
	viper.SetDefault("timeout", "30s")
	viper.SetDefault("retry_attempts", 3)
	viper.SetDefault("retry_delay", "1s")
	viper.SetDefault("output_format", "text")
	viper.SetDefault("verbose", false)
	viper.SetDefault("status_interval", "1s")
	viper.SetDefault("command_delay", "100ms")

	viper.SetEnvPrefix("FLUIDNC")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	var config fluidnc.Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}
