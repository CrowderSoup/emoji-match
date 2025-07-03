package config

import (
	"os"

	"go.uber.org/fx"
)

// Config our app config
type Config struct {
	ServiceName     string
	Address         string
	Environment     string
	DBUser          string
	DBPassword      string
	DBHost          string
	AuthSecret      string
	TracingEndpoint string
}

// LoadConfig loads and returns our app config
func LoadConfig() *Config {
	config := Config{
		ServiceName:     getEnv("SERVICENAME", "morphe-api-service"),
		Address:         getEnv("ADDRESS", ":3001"),
		Environment:     getEnv("ENVIRONMENT", "development"),
		DBUser:          getEnv("DBUSER", "sa"),
		DBPassword:      getEnv("DBPASSWORD", "P@ssw0rd"),
		DBHost:          getEnv("DBHOST", "localhost:1433"),
		AuthSecret:      getEnv("AUTHSECRET", "some-secret-key"),
		TracingEndpoint: getEnv("TRACINGENDPOINT", "localhost:4318"),
	}
	return &config
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// Module provided to fx
var Module = fx.Options(
	fx.Provide(
		LoadConfig,
	),
)
