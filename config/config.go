package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	RedisAddr        string
	RedisPass        string
	LimitPerIP       int
	LimitPerToken    int
	BlockTimeSeconds int
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	limitPerIP, _ := strconv.Atoi(os.Getenv("LIMIT_PER_IP"))
	limitPerToken, _ := strconv.Atoi(os.Getenv("LIMIT_PER_TOKEN"))
	blockTimeSeconds, _ := strconv.Atoi(os.Getenv("BLOCK_TIME_SECONDS"))

	return &Config{
		RedisAddr:        os.Getenv("REDIS_ADDR"),
		RedisPass:        os.Getenv("REDIS_PASS"),
		LimitPerIP:       limitPerIP,
		LimitPerToken:    limitPerToken,
		BlockTimeSeconds: blockTimeSeconds,
	}
}
