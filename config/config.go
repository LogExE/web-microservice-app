package config

import "os"

type BrokerConfig struct {
	BootstrapAddr string
}

type RESTConfig struct {
	Addr string
}

type Config struct {
	Broker BrokerConfig
	REST   RESTConfig
}

func New() *Config {
	return &Config{
		Broker: BrokerConfig{
			BootstrapAddr: getEnv("KFK_BOOTSTRAP_ADDR", "localhost:9092"),
		},
		REST: RESTConfig{
			Addr: getEnv("REST_ADDR", "localhost:8080"),
		},
	}
}

// Simple helper function to read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
