package app

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Abdulazizxoshimov/Datagaze_Backend/api"
	"github.com/Abdulazizxoshimov/Datagaze_Backend/api/server"
	"github.com/Abdulazizxoshimov/Datagaze_Backend/config"
	"github.com/Abdulazizxoshimov/Datagaze_Backend/internal/pkg/logger"
	"github.com/Abdulazizxoshimov/Datagaze_Backend/internal/pkg/storage"
	"github.com/Abdulazizxoshimov/Datagaze_Backend/internal/repo"
	"github.com/Abdulazizxoshimov/Datagaze_Backend/internal/repo/postgres"
	"github.com/Abdulazizxoshimov/Datagaze_Backend/internal/repo/redis"
	"github.com/Abdulazizxoshimov/Datagaze_Backend/internal/service"
	"github.com/casbin/casbin/v2"
	defaultrolemanager "github.com/casbin/casbin/v2/rbac/default-role-manager"
	"github.com/casbin/casbin/v2/util"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type App struct {
	Config      config.Config
	Logger      logger.Logger
	server      *http.Server
	DB          *sqlx.DB
	Enforcer    *casbin.Enforcer
	RedisDB     *storage.RedisDB
	StorageI    repo.StorageI
	WeatherCron *service.WeatherService // Weather cron servisi
	stopChan    chan struct{}        // Cronni to‘xtatish uchun
}

func NewApp(cfg config.Config) (*App, error) {
	// init logger
	logger, err := logger.New(cfg.LogLevel, cfg.Environment, cfg.App+".log")
	if err != nil {
		return nil, err
	}

	//init redis
	redisdb, err := storage.NewRedis(&cfg)
	if err != nil {
		return nil, err
	}

	//init casbin enforcer
	enforcer, err := casbin.NewEnforcer("./config/auth.conf", "./config/auth.csv")
	if err != nil {
		return nil, err
	}

	// init db
	db, err := storage.NewSqlDatabase(&cfg)
	if err != nil {
		return nil, err
	}

	storageI := repo.NewStoragePG(db, logger)

	// Weather service uchun repo yaratish
	weatherRepo := postgres.NewWeatherRepo(db)

	// Weather cron service yaratish
	weatherCron := service.NewWeatherService(weatherRepo, cfg)
	

	return &App{
		Config:      cfg,
		Logger:      logger,
		DB:          db,
		RedisDB:     redisdb,
		Enforcer:    enforcer,
		StorageI:    storageI,
		WeatherCron: weatherCron,
		stopChan:    make(chan struct{}),
	}, nil
}

func (a *App) Run() error {
	// initialize cache
	cache := redis.NewRedis(a.RedisDB)

	
	// api init
	handler := api.NewRoute(api.RouteOption{
		Config:         a.Config,
		Logger:         a.Logger,
		ContextTimeout: a.Config.Context.Timeout,
		Cache:          cache,
		Enforcer:       a.Enforcer,
		Service:        a.StorageI,
		WeatherService: a.WeatherCron,
	})

	// Casbin policy load
	err := a.Enforcer.LoadPolicy()
	if err != nil {
		return err
	}
	roleManager := a.Enforcer.GetRoleManager().(*defaultrolemanager.RoleManagerImpl)
	roleManager.AddMatchingFunc("keyMatch", util.KeyMatch)
	roleManager.AddMatchingFunc("keyMatch3", util.KeyMatch3)

	// **Cron jobni ishga tushirish**
	go a.startWeatherCron()

	// server init
	a.server, err = server.NewServer(&a.Config, handler)
	if err != nil {
		return fmt.Errorf("error while initializing server: %v", err)
	}

	return a.server.ListenAndServe()
}

func (a *App) startWeatherCron() {
	ticker := time.NewTicker(10 * time.Minute) // Har 10 daqiqada yangilash
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			cities := []string{"Tashkent", "New York", "Tokyo"} // Shaharlar ro‘yxati
			for _, city := range cities {
				err := a.WeatherCron.FetchAndStoreWeather(city)
				if err != nil {
					a.Logger.Error("Weather cron job failed for %s: %v", zap.Error(err))
				}
			}
		case <-a.stopChan:
			a.Logger.Info("Weather cron job stopped")
			return
		}
	}
}

func (a *App) Stop() {
	// Cron jobni to‘xtatish
	close(a.stopChan)

	// database connection
	a.DB.Close()

	// zap logger sync
	a.Logger.Sync()
}
