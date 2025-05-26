package config

// App config

type App struct {
	Name               string `mapstructure:"name"`
	Host               string `mapstructure:"host"` // Bind address (0.0.0.0 for listening)
	Port               int    `mapstructure:"port"`
	RegisterAddress    string `mapstructure:"register_address"`     // Address for service registration
	HealthCheckAddress string `mapstructure:"health_check_address"` // Address for Consul health checks
}

// Consul config

type Consul struct {
	Host   string `mapstructure:"host"`
	Port   int    `mapstructure:"port"`
	Scheme string `mapstructure:"scheme"`
}
