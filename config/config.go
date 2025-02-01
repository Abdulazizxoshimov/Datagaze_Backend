package config

import (
	"os"
	"time"
)

type Config struct {
	App         string
	Environment string
	LogLevel    string

	WeatherAPIKey string

	Server struct {
		Host         string
		Port         string
		ReadTimeout  string
		WriteTimeout string
		IdleTimeout  string
	}
	DB struct {
		Host     string
		Port     string
		Name     string
		User     string
		Password string
		SslMode  string
	}

	Context struct {
		Timeout time.Duration
	}

	Token struct {
		Secret     string
		AccessTTL  time.Duration
		RefreshTTL time.Duration
		SignInKey  string
	}

	SMTP struct {
		Email         string
		EmailPassword string
		SMTPPort      string
		SMTPHost      string
	}

	Redis struct {
		Host     string
		Port     string
		Password string
		Name     string
	}
}

func NewConfig() (*Config, error) {
	var config Config

	// general configuration
	config.App = getEnv("APP", "app")
	config.Environment = getEnv("ENVIRONMENT", "develop")
	config.LogLevel = getEnv("LOG_LEVEL", "debug")

	// server configuration
	config.Server.Host = getEnv("SERVER_HOST", "app") //app
	config.Server.Port = getEnv("SERVER_PORT", ":7777")
	config.Server.ReadTimeout = getEnv("SERVER_READ_TIMEOUT", "10s")
	config.Server.WriteTimeout = getEnv("SERVER_WRITE_TIMEOUT", "10s")
	config.Server.IdleTimeout = getEnv("SERVER_IDLE_TIMEOUT", "120s")

	// db configuration
	config.DB.Host = getEnv("POSTGRES_HOST", "postgres") // postgres
	config.DB.Port = getEnv("POSTGRES_PORT", "5432")
	config.DB.User = getEnv("POSTGRES_USER", "postgres")
	config.DB.Password = getEnv("POSTGRES_PASSWORD", "4444")
	config.DB.SslMode = getEnv("POSTGRES_SSLMODE", "disable")
	config.DB.Name = getEnv("POSTGRES_DATABASE", "datagaze_backend")

	// redis configuration
	config.Redis.Host = getEnv("REDIS_HOST", "redisdb") //redisdb
	config.Redis.Port = getEnv("REDIS_PORT", "6379")
	config.Redis.Password = getEnv("REDIS_PASSWORD", "")
	config.Redis.Name = getEnv("REDIS_DATABASE", "0")

	//smtp confifuration
	config.SMTP.Email = getEnv("SMTP_EMAIL", "theuniver77@gmail.com")
	config.SMTP.EmailPassword = getEnv("SMTP_EMAIL_PASSWORD", "fywqgrsyhvybjyxa")
	config.SMTP.SMTPPort = getEnv("SMTP_PORT", "587")
	config.SMTP.SMTPHost = getEnv("SMTP_HOST", "smtp.gmail.com")

	//weather api key configuration
	config.WeatherAPIKey = getEnv("WEATHER_API_KEY", "27b3b6dd5b3b4639939195251253101")
	//context configuration
	ContexTimeout, err := time.ParseDuration(getEnv("CONTEXT_TIMEOUT", "30s"))
	if err != nil {
		return nil, err
	}

	config.Context.Timeout = ContexTimeout

	// access ttl parse
	accessTTl, err := time.ParseDuration(getEnv("TOKEN_ACCESS_TTL", "3h"))
	if err != nil {
		return nil, err
	}
	// refresh ttl parse
	refreshTTL, err := time.ParseDuration(getEnv("TOKEN_REFRESH_TTL", "24h"))
	if err != nil {
		return nil, err
	}
	config.Token.AccessTTL = accessTTl
	config.Token.RefreshTTL = refreshTTL
	config.Token.SignInKey = getEnv("TOKEN_SIGNIN_KEY", "debug")

	return &config, nil
}

func getEnv(key string, defaultVaule string) string {
	value, exists := os.LookupEnv(key)
	if exists {
		return value
	}
	return defaultVaule
}
