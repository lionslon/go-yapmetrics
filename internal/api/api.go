package api

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/lionslon/go-yapmetrics/internal/config"
	"github.com/lionslon/go-yapmetrics/internal/handlers"
	"github.com/lionslon/go-yapmetrics/internal/middlewares"
	"github.com/lionslon/go-yapmetrics/internal/storage"
	"github.com/lionslon/go-yapmetrics/internal/storing"
	"go.uber.org/zap"
	"log"
)

type APIServer struct {
	echo *echo.Echo
	addr string
	//sugar zap.SugaredLogger
}

func New() *APIServer {
	apiS := &APIServer{}
	apiS.echo = echo.New()
	cfg := config.ServerConfig{}
	st := storage.New()
	handler := handlers.New(st)
	cfg.New()
	apiS.addr = cfg.Addr

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	sugar := *logger.Sugar()

	if cfg.FilePath != "" {
		if cfg.Restore {
			storing.Restore(st, cfg.FilePath)
		}
		if cfg.StoreInterval != 0 {
			go storing.Store(st, cfg.FilePath, cfg.StoreInterval)
		}
	}

	apiS.echo.Use(middlewares.WithLogging(sugar))
	apiS.echo.Use(middlewares.GzipUnpacking())

	apiS.echo.GET("/", handler.AllMetricsValues())
	apiS.echo.POST("/value/", handler.GetValueJSON())
	apiS.echo.GET("/value/:typeM/:nameM", handler.MetricsValue())
	apiS.echo.POST("/update/", handler.UpdateJSON())
	apiS.echo.POST("/update/:typeM/:nameM/:valueM", handler.PostWebhandle())

	return apiS
}

func (a *APIServer) Start() error {
	fmt.Println("Running server on", a.addr)
	err := a.echo.Start(a.addr)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
