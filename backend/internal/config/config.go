package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DatabaseURL      string
	JWTSecret        string
	JWTExpiry        time.Duration
	ServerPort       string
	UploadDir        string
	AllowedOrigins   []string
	AllowCredentials bool
}

func LoadConfig() *Config {
	// Load .env silently
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	_ = viper.ReadInConfig() // No log, no panic

	// Set default values if not provided
	viper.SetDefault("DATABASE_URL", "")
	viper.SetDefault("JWT_SECRET", "")
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("UPLOAD_DIR", "./uploads")
	viper.SetDefault("FRONTEND_URL", "http://localhost:3000")

	return &Config{
		DatabaseURL:      viper.GetString("DATABASE_URL"),
		JWTSecret:        viper.GetString("JWT_SECRET"),
		JWTExpiry:        24 * time.Hour,
		ServerPort:       viper.GetString("PORT"),
		UploadDir:        viper.GetString("UPLOAD_DIR"),
		AllowedOrigins:   []string{viper.GetString("FRONTEND_URL")},
		AllowCredentials: true,
	}
}
