package api

import (
	"github.com/labstack/echo/v4"
	"github.com/lionslon/go-yapmetrics/internal/config"
	"github.com/lionslon/go-yapmetrics/internal/database"
	"github.com/lionslon/go-yapmetrics/internal/handlers"
	"github.com/lionslon/go-yapmetrics/internal/middlewares"
	"github.com/lionslon/go-yapmetrics/internal/storage"
	"github.com/lionslon/go-yapmetrics/internal/storing"
	"go.uber.org/zap"
	"log"
)

type APIServer struct {
	cfg  config.ServerConfig
	echo *echo.Echo
	st   *storage.MemStorage
	db   *database.DBConnection
}

func New() *APIServer {
	apiS := &APIServer{}
	cfg := config.ServerConfig{}
	apiS.cfg = cfg.New()
	apiS.echo = echo.New()
	apiS.st = storage.New()
	handler := handlers.New(apiS.st)
	apiS.db = database.New(apiS.cfg.DatabaseDSN)

	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)

	defer logger.Sync()

	apiS.echo.Use(middlewares.WithLogging())
	apiS.echo.Use(middlewares.GzipUnpacking())

	apiS.echo.GET("/", handler.AllMetricsValues())
	apiS.echo.POST("/value/", handler.GetValueJSON())
	apiS.echo.GET("/value/:typeM/:nameM", handler.MetricsValue())
	apiS.echo.POST("/update/", handler.UpdateJSON())
	apiS.echo.POST("/update/:typeM/:nameM/:valueM", handler.UpdateMetrics())
	apiS.echo.GET("/ping", handler.PingDB(apiS.db))

	return apiS
}

func (a *APIServer) Start() error {
	if a.cfg.FilePath != "" {
		if a.cfg.Restore {
			storing.Restore(a.st, a.cfg.FilePath)
		}
		if a.cfg.StoreInterval != 0 {
			go storing.IntervalDump(a.st, a.cfg.FilePath, a.cfg.StoreInterval)
		}
	}

	err := a.echo.Start(a.cfg.Addr)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
