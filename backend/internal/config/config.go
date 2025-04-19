package config

import (
	"time"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	DatabaseURL      string
	JWTSecret        string
	JWTExpiry        time.Duration
	ServerPort       string
	AllowedOrigins   []string
	AllowCredentials bool
	GoogleClientID   string
    GoogleClientSecret string
    OAuthCallbackURL string
	FrontendURL string
	SessionSecret string
}

func LoadConfig() *Config {
	// Load .env silently
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	_ = viper.ReadInConfig() // No log, no panic

	// Set default values if not provided
	viper.SetDefault("DATABASE_URL", "")
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("FRONTEND_URL", "http://localhost:3000")
	viper.SetDefault("GOOGLE_CLIENT_ID", "")
    viper.SetDefault("GOOGLE_CLIENT_SECRET", "")
    viper.SetDefault("OAUTH_CALLBACK_URL", "http://localhost:8080/auth/google/callback")
	viper.SetDefault("SESSION_SECRET","")


	originsStr := viper.GetString("ALLOWED_ORIGINS")
    origins := []string{}
    if originsStr != "" {
        // Split by comma and trim spaces
        for _, origin := range strings.Split(originsStr, ",") {
            origins = append(origins, strings.TrimSpace(origin))
        }
    }

	if len(origins) == 0 {
        origins = []string{viper.GetString("FRONTEND_URL")}
    }

	return &Config{
		DatabaseURL:      viper.GetString("DATABASE_URL"),
		JWTSecret:        viper.GetString("JWT_SECRET"),
		JWTExpiry:        24 * time.Hour,
		ServerPort:       viper.GetString("PORT"),
		AllowedOrigins:   []string{viper.GetString("FRONTEND_URL")},
		AllowCredentials: true,
		GoogleClientID:     viper.GetString("GOOGLE_CLIENT_ID"),
        GoogleClientSecret: viper.GetString("GOOGLE_CLIENT_SECRET"),
        OAuthCallbackURL:   viper.GetString("OAUTH_CALLBACK_URL"),
		FrontendURL: viper.GetString("FRONTEND_URL"),
		SessionSecret: viper.GetString("SESSION_SECRET"),
	}
}
