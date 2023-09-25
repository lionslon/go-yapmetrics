package api

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/lionslon/go-yapmetrics/internal/config"
	"github.com/lionslon/go-yapmetrics/internal/handlers"
	"github.com/lionslon/go-yapmetrics/internal/storage"
	"log"
)

type APIServer struct {
	echo *echo.Echo
	addr string
}

func New() *APIServer {
	apiS := &APIServer{}
	apiS.echo = echo.New()
	cfg := config.ServerConfig{}
	st := storage.New()
	handler := handlers.New(st)
	cfg.New()
	apiS.addr = cfg.Addr

	apiS.echo.GET("/", handler.AllMetricsValues())
	apiS.echo.GET("/value/:typeM/:nameM", handler.MetricsValue())
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
