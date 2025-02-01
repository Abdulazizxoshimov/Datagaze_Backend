package api

import (
	"time"

	_ "github.com/Abdulazizxoshimov/Datagaze_Backend/api/docs"
	"github.com/Abdulazizxoshimov/Datagaze_Backend/api/handler"
	"github.com/Abdulazizxoshimov/Datagaze_Backend/api/middleware"
	"github.com/Abdulazizxoshimov/Datagaze_Backend/config"
	"github.com/Abdulazizxoshimov/Datagaze_Backend/internal/pkg/logger"
	"github.com/Abdulazizxoshimov/Datagaze_Backend/internal/pkg/token"
	"github.com/Abdulazizxoshimov/Datagaze_Backend/internal/repo"
	"github.com/Abdulazizxoshimov/Datagaze_Backend/internal/repo/interfaces"
	"github.com/Abdulazizxoshimov/Datagaze_Backend/internal/service"
	"github.com/casbin/casbin/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type RouteOption struct {
	Config         config.Config
	Logger         logger.Logger
	ContextTimeout time.Duration
	Cache          interfaces.Redis
	Enforcer       *casbin.Enforcer
	RefreshToken   token.JWTHandler
	Service        repo.StorageI
	WeatherService    *service.WeatherService // Weather cron servisi

}

// NewRoute
// @Title Datagaze_Backend
// @Description Contacs: https://t.me/Abuzada0401
// @securityDefinitions.apikey BearerAuth
// @in 			header
// @name 		Authorization
func NewRoute(option RouteOption) *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	HandlerV1 := handler.New(&handler.HandlerV1Config{
		Config:         option.Config,
		Logger:         option.Logger,
		ContextTimeout: option.ContextTimeout,
		Redis:          option.Cache,
		RefreshToken:   option.RefreshToken,
		Enforcer:       option.Enforcer,
		Service:        option.Service,
		WeatherService:    option.WeatherService,
	})

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowCredentials = true
	corsConfig.AllowHeaders = []string{"*"}
	corsConfig.AllowBrowserExtensions = true
	corsConfig.AllowMethods = []string{"*"}
	router.Use(cors.New(corsConfig))

	router.Use(middleware.CheckCasbinPermission(option.Enforcer, option.Config))
	router.Static("/media", "./file")

	// login
	router.POST("/register", HandlerV1.Register)
	router.POST("/login", HandlerV1.Login)
	router.POST("/forgot/:email", HandlerV1.Forgot)
	router.POST("/verify", HandlerV1.VerifyOTP)
	router.PUT("/reset-password", HandlerV1.ResetPassword)
	router.GET("/token/:refresh", HandlerV1.Token)
	router.POST("/users/verify", HandlerV1.Verify)


	//user
	router.POST("/user", HandlerV1.CreateUser)
	router.PUT("/user/:id", HandlerV1.UpdateUser)
	router.DELETE("/user/:id", HandlerV1.DeleteUser)
	router.GET("/user/:id", HandlerV1.GetUser)
	router.GET("/users", HandlerV1.GetAllUsers)

	//weather
	router.GET("/weather", HandlerV1.GetWeather)

	url := ginSwagger.URL("swagger/doc.json")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	return router

}
