package api

import (
	"time"

	_ "greenShop/api/docs"
	v1 "greenShop/api/handlers/v1"
	"greenShop/api/middleware"

	"greenShop/internal/infrastructure/clientService"
	redisrepo "greenShop/internal/infrastructure/repository/redisdb"

	"greenShop/internal/pkg/config"
	"greenShop/internal/pkg/token"

	"github.com/casbin/casbin/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

type RouteOption struct {
	Config         config.Config
	Logger         *zap.Logger
	ContextTimeout time.Duration
	Cache          redisrepo.Cache
	Enforcer       *casbin.Enforcer
	RefreshToken   token.JWTHandler
	Service        clientService.ServiceClient
	MinIO          *minio.Client
}

// NewRoute
// @Title GreenShop
// @Description Contacs: https://t.me/Abuzada0401
// @securityDefinitions.apikey BearerAuth
// @in 			header
// @name 		Authorization
func NewRoute(option RouteOption) *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	HandlerV1 := v1.New(&v1.HandlerV1Config{
		Config:         option.Config,
		Logger:         option.Logger,
		ContextTimeout: option.ContextTimeout,
		Redis:          option.Cache,
		RefreshToken:   option.RefreshToken,
		Enforcer:       option.Enforcer,
		Service:        option.Service,
		MinIO:          option.MinIO,
	})

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowCredentials = true
	corsConfig.AllowHeaders = []string{"*"}
	corsConfig.AllowBrowserExtensions = true
	corsConfig.AllowMethods = []string{"*"}
	router.Use(cors.New(corsConfig))

	router.Use(middleware.Tracing)
	router.Use(middleware.CheckCasbinPermission(option.Enforcer, option.Config))

	router.Static("/media", "./media")
	router.Static("/swagger", "./swagger")

	apiV1 := router.Group("/v1")

	// address
	apiV1.POST("/address", HandlerV1.CreateAddress)
	apiV1.PUT("/address", HandlerV1.UpdateAddress)
	apiV1.DELETE("/address/:id", HandlerV1.DeleteAddress)
	apiV1.GET("/address/:id", HandlerV1.GetAddress)
	apiV1.GET("/list/address", HandlerV1.GetAllAddress)

	// login
	apiV1.POST("/register", HandlerV1.Register)
	apiV1.POST("/login", HandlerV1.Login)
	apiV1.POST("/forgot/:email", HandlerV1.Forgot)
	apiV1.POST("/verify", HandlerV1.VerifyOTP)
	apiV1.PUT("/reset-password", HandlerV1.ResetPassword)
	apiV1.GET("/token/:refresh", HandlerV1.Token)
	apiV1.POST("/users/verify", HandlerV1.Verify)

	//user
	apiV1.POST("/user", HandlerV1.CreateUser)
	apiV1.PUT("/user", HandlerV1.UpdateUser)
	apiV1.DELETE("/user/:id", HandlerV1.DeleteUser)
	apiV1.GET("/user/:id", HandlerV1.GetUser)
	apiV1.GET("/users", HandlerV1.ListUsers)
	apiV1.PUT("/user/profile", HandlerV1.UpdateProfile)
	apiV1.PUT("/user/password", HandlerV1.UpdatePassword)

	//post
	apiV1.POST("/product", HandlerV1.CreateProduct)
	apiV1.PUT("/product", HandlerV1.UpdateProduct)
	apiV1.DELETE("/product/:id", HandlerV1.DeleteProduct)
	apiV1.GET("/product/:id", HandlerV1.GetProduct)
	apiV1.GET("/products", HandlerV1.ListProducts)

	// category
	apiV1.POST("/category", HandlerV1.CreateCategory)
	apiV1.PUT("/category", HandlerV1.UpdateCategory)
	apiV1.DELETE("/category/:id", HandlerV1.DeleteCategory)
	apiV1.GET("/category/:id", HandlerV1.GetCategory)
	apiV1.GET("/categories", HandlerV1.ListCategory)

	//google
	apiV1.GET("/google/callback", HandlerV1.GoogleCallback)
	apiV1.GET("google/login", HandlerV1.GoogleLogin)

	//comment
	apiV1.POST("/comment", HandlerV1.CreateComment)
	apiV1.PUT("/comment", HandlerV1.UpdateComment)
	apiV1.DELETE("/comment/:id", HandlerV1.DeleteComment)
	apiV1.GET("/comment/:id", HandlerV1.GetComment)
	apiV1.GET("/comments", HandlerV1.ListComment)
	apiV1.GET("/post/comments", HandlerV1.GetAllCommentByPostId)

	// Media
	apiV1.POST("/media/upload-photo", HandlerV1.UploadMedia)
	apiV1.GET("/media/:id", HandlerV1.GetMedia)
	apiV1.DELETE("/media/:id", HandlerV1.DeleteMedia)

	// Wishlist
	apiV1.POST("/like/:id", HandlerV1.LikeProduct)
	apiV1.GET("/wishlist", HandlerV1.UserWishlist)

	// Basket, Stats, Payment ...
	apiV1.POST("/basket", HandlerV1.SaveToBasket)
	apiV1.GET("/basket", HandlerV1.GetBasketProduct)

	// Order 
	apiV1.POST("/order", HandlerV1.CreateOrder)
	apiV1.GET("/order/:id", HandlerV1.GetOrder)
	apiV1.GET("/orders", HandlerV1.ListOrders)
	apiV1.PUT("/order", HandlerV1.UpdateOrder)
	apiV1.DELETE("/order/:id", HandlerV1.DeleteOrder)


	url := ginSwagger.URL("swagger/doc.json")
	apiV1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	return router
}
