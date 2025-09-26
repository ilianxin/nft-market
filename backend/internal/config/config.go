package config

import (
	"os"
)

// Config 应用配置结构体
type Config struct {
	Port            string
	DatabaseURL     string
	EthereumRPC     string
	ContractAddress string
	PrivateKey      string
	Environment     string
	JWTSecret       string
}

// Load 加载配置
func Load() *Config {
	return &Config{
		Port:            getEnv("PORT", "8080"),
		DatabaseURL:     getEnv("DATABASE_URL", "postgres://username:password@localhost/nft_market?sslmode=disable"),
		EthereumRPC:     getEnv("ETHEREUM_RPC", "http://localhost:8545"),
		ContractAddress: getEnv("CONTRACT_ADDRESS", ""),
		PrivateKey:      getEnv("PRIVATE_KEY", ""),
		Environment:     getEnv("ENVIRONMENT", "development"),
		JWTSecret:       getEnv("JWT_SECRET", "your-secret-key"),
	}
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
