package config

import "os"

type Config struct {
	RedisAddr     string
	RedisStream   string
	ServerAddress string
}

func Load() *Config {
	return &Config{
		RedisAddr:     getEnv("REDIS_ADDR", ":6379"),
		RedisStream:   getEnv("REDIS_QUEUE", "local-q"),
		ServerAddress: getEnv("SERVER_ADDR", ":3000"),
	}
}

func getEnv(key string, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
