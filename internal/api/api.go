package api

import (
	"github.com/labstack/echo/v4"
	"github.com/lionslon/yapmetrics/internal/handlers"
	"github.com/lionslon/yapmetrics/internal/storage"
	"log"
)

type APIServer struct {
	storage *storage.MemStorage
	echo    *echo.Echo
}

func New() *APIServer {
	apiS := &APIServer{}
	apiS.storage = storage.New()
	apiS.echo = echo.New()

	apiS.echo.GET("/", handlers.AllMetrics(apiS.storage))
	apiS.echo.GET("/value/:typeM/:nameM", handlers.MetricsValue(apiS.storage))
	apiS.echo.POST("/update/:typeM/:nameM/:valueM", handlers.PostWebhandle(apiS.storage))

	return apiS
}

func (a *APIServer) Start() error {
	err := a.echo.Start(":8080")
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
