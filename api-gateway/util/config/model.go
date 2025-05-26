package config

// App config

type App struct {
	Name string `mapstructure:"name"`
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

// Consul config
type Consul struct {
	Host   string `mapstructure:"host"`
	Port   int    `mapstructure:"port"`
	Scheme string `mapstructure:"scheme"`
}
