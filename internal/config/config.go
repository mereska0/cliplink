package config

import "os"

type Config struct {
	HTTPAddr       string
	GRPCAddr       string
	GRPCClientAddr string
	DatabaseURL    string
}

func Load() Config {
	return Config{
		HTTPAddr:       getEnv("HTTP_ADDR", ":8080"),
		GRPCAddr:       getEnv("GRPC_ADDR", ":50051"),
		GRPCClientAddr: getEnv("GRPC_CLIENT_ADDR", "localhost:50051"),
		DatabaseURL:    getEnv("DATABASE_URL", "postgres://cliplink:cliplink@localhost:5433/cliplink?sslmode=disable"),
	}
}

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	return value
}
