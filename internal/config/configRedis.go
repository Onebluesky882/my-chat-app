package config

import (
	"github.com/Onebluesky882/my-chat-app/libs"
)

type Config struct {
	RedisAddr     string
	RedisPassword string
	RedisDB       int
	ScyllaHosts   []string
}

func LoadConfig() *Config {

	return &Config{
		RedisAddr:     libs.GetEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: libs.GetEnv("REDIS_PASSWORD", ""),
		RedisDB:       0,
		ScyllaHosts:   []string{libs.GetEnv("SCYLLA_HOST", "127.0.0.1")},
	}
}
