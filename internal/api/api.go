package api

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lionslon/go-yapmetrics/internal/config"
	"github.com/lionslon/go-yapmetrics/internal/handlers"
	"github.com/lionslon/go-yapmetrics/internal/middlewares"
	"github.com/lionslon/go-yapmetrics/internal/storage"
	"github.com/lionslon/go-yapmetrics/pkg/utils/profile"
	"go.uber.org/zap"
	"log"
)

type APIServer struct {
	cfg  *config.ServerConfig
	echo *echo.Echo
	st   *storage.MemStorage
}

func New() *APIServer {
	apiS := &APIServer{}
	cfg := config.NewServer()
	apiS.cfg = cfg
	apiS.echo = echo.New()
	apiS.st = storage.NewMemoryStorage()

	if cfg.EnableProfiling {
		profile.StartProfilingServer()
	}

	handler := handlers.New(apiS.st)
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
	defer logger.Sync()

	var storageProvider storage.StorageWorker
	var err error
	switch cfg.GetProvider() {
	case storage.FileProvider:
		storageProvider = storage.NewFileProvider(cfg.FilePath, cfg.StoreInterval, apiS.st)
	case storage.DBProvider:
		storageProvider, err = storage.NewDBProvider(cfg.DatabaseDSN, cfg.StoreInterval, apiS.st)
	}
	if err != nil {
		zap.S().Error(err)
	}
	if cfg.Restore {
		err := storageProvider.Restore()
		if err != nil {
			zap.S().Error(err)
		}
	}

	if cfg.StoreIntervalNotZero() {
		go storageProvider.IntervalDump()
	}

	apiS.echo.Use(middlewares.WithLogging())
	//apiS.echo.Use(middlewares.GzipUnpacking())
	apiS.echo.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))
	if cfg.SignPass != "" {
		apiS.echo.Use(middlewares.CheckSignReq(cfg.SignPass))
	}

	apiS.echo.GET("/", handler.AllMetricsValues())
	apiS.echo.POST("/value/", handler.GetValueJSON())
	apiS.echo.GET("/value/:typeM/:nameM", handler.MetricsValue())
	apiS.echo.POST("/update/", handler.UpdateJSON())
	apiS.echo.POST("/update/:typeM/:nameM/:valueM", handler.UpdateMetrics())
	apiS.echo.POST("/updates/", handler.UpdatesJSON())
	apiS.echo.GET("/ping", handler.PingDB(storageProvider))

	return apiS
}

func (a *APIServer) Start() error {
	err := a.echo.Start(a.cfg.Addr)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
