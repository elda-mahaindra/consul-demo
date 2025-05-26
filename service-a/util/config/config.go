package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	App    App    `mapstructure:"app"`
	Consul Consul `mapstructure:"consul"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("json")

	// Enable automatic environment variable reading
	viper.AutomaticEnv()

	// Set up environment variable mappings for nested config
	viper.BindEnv("consul.host", "CONSUL_HOST")
	viper.BindEnv("consul.port", "CONSUL_PORT")
	viper.BindEnv("consul.scheme", "CONSUL_SCHEME")

	err = viper.ReadInConfig()
	if err != nil {
		return config, fmt.Errorf("failed to read configuration file: %s", err)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return config, fmt.Errorf("failed to unmarshal configuration: %s", err)
	}

	// Log the final Consul configuration for debugging
	consulHost := config.Consul.Host
	if envHost := os.Getenv("CONSUL_HOST"); envHost != "" {
		consulHost = envHost
		config.Consul.Host = envHost
	}

	fmt.Printf("🔧 Consul configuration: %s://%s:%d\n", config.Consul.Scheme, consulHost, config.Consul.Port)

	return
}
