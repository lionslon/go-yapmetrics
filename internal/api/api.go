package api

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lionslon/go-yapmetrics/internal/config"
	"github.com/lionslon/go-yapmetrics/internal/handlers"
	"github.com/lionslon/go-yapmetrics/internal/middlewares"
	"github.com/lionslon/go-yapmetrics/internal/storage"
	"github.com/lionslon/go-yapmetrics/pkg/utils/profile"
	"go.uber.org/zap"
	"os/signal"
	"syscall"
)

type APIServer struct {
	cfg             *config.ServerConfig
	echo            *echo.Echo
	st              *storage.MemStorage
	storageProvider storage.StorageWorker
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

	var err error
	switch cfg.GetProvider() {
	case storage.FileProvider:
		apiS.storageProvider = storage.NewFileProvider(cfg.FilePath, cfg.StoreInterval, apiS.st)
	case storage.DBProvider:
		apiS.storageProvider, err = storage.NewDBProvider(cfg.DatabaseDSN, cfg.StoreInterval, apiS.st)
	}
	if err != nil {
		zap.S().Error(err)
	}
	if cfg.Restore {
		err := apiS.storageProvider.Restore()
		if err != nil {
			zap.S().Error(err)
		}
	}

	if cfg.StoreIntervalNotZero() {
		go apiS.storageProvider.IntervalDump()
	}

	apiS.echo.Use(middlewares.WithLogging())
	//apiS.echo.Use(middlewares.GzipUnpacking())
	apiS.echo.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))
	if cfg.SignPass != "" {
		apiS.echo.Use(middlewares.CheckSignReq(cfg.SignPass))
	}
	if cfg.CryptoKey != "" {
		apiS.echo.Use(middlewares.DecryptBody(cfg.SignPass))
	}

	apiS.echo.GET("/", handler.AllMetricsValues())
	apiS.echo.POST("/value/", handler.GetValueJSON())
	apiS.echo.GET("/value/:typeM/:nameM", handler.MetricsValue())
	apiS.echo.POST("/update/", handler.UpdateJSON())
	apiS.echo.POST("/update/:typeM/:nameM/:valueM", handler.UpdateMetrics())
	apiS.echo.POST("/updates/", handler.UpdatesJSON())
	apiS.echo.GET("/ping", handler.PingDB(apiS.storageProvider))

	return apiS
}

func (a *APIServer) Start() error {
	rootContext, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	err := a.echo.Start(a.cfg.Addr)
	if err != nil {
		zap.S().Error(err)
	}
	go a.gracefulShutdown(rootContext)

	return nil
}

// gracefulShutdown - Запускается в получении любого из сигнала (syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
func (a *APIServer) gracefulShutdown(ctx context.Context) {
	<-ctx.Done()

	zap.S().Info("shutting down gracefully")
	err := a.storageProvider.Close()
	if err != nil {
		zap.S().Error("error during storage closing")
	}

	err = a.echo.Shutdown(context.Background())
	if err != nil {
		zap.S().Error("error during server shutdown")
	}
}
