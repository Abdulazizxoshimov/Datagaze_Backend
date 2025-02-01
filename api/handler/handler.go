package handler

import (
	"time"

	"github.com/Abdulazizxoshimov/Datagaze_Backend/config"
	"github.com/Abdulazizxoshimov/Datagaze_Backend/internal/pkg/logger"
	"github.com/Abdulazizxoshimov/Datagaze_Backend/internal/pkg/token"
	"github.com/Abdulazizxoshimov/Datagaze_Backend/internal/repo"
	redis "github.com/Abdulazizxoshimov/Datagaze_Backend/internal/repo/interfaces"
	"github.com/Abdulazizxoshimov/Datagaze_Backend/internal/service"
	"github.com/casbin/casbin/v2"
)

type HandlerV1 struct {
	Config         config.Config
	Logger         logger.Logger
	ContextTimeout time.Duration
	RedisStorage   redis.Redis
	RefreshToken   token.JWTHandler
	Enforcer       *casbin.Enforcer
	Service        repo.StorageI
	WeatherService    *service.WeatherService // Weather cron servisi

}

// HandlerV1Config ...
type HandlerV1Config struct {
	Config         config.Config
	Logger         logger.Logger
	ContextTimeout time.Duration
	Redis          redis.Redis
	RefreshToken   token.JWTHandler
	Enforcer       *casbin.Enforcer
	Service        repo.StorageI
	WeatherService    *service.WeatherService // Weather cron servisi

}

// New ...
func New(c *HandlerV1Config) *HandlerV1 {
	return &HandlerV1{
		Config:         c.Config,
		Logger:         c.Logger,
		ContextTimeout: c.ContextTimeout,
		RedisStorage:   c.Redis,
		Enforcer:       c.Enforcer,
		RefreshToken:   c.RefreshToken,
		Service:        c.Service,
		WeatherService:    c.WeatherService,
	}
}
