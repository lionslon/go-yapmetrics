package api

import (
	"flag"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/lionslon/go-yapmetrics/internal/handlers"
	"github.com/lionslon/go-yapmetrics/internal/storage"
	"log"
)

type APIServer struct {
	storage *storage.MemStorage
	echo    *echo.Echo
	addr    string
}

func New() *APIServer {
	apiS := &APIServer{}
	apiS.storage = storage.New()
	apiS.echo = echo.New()
	address := flag.String("a", "localhost:8080", "address and port to run server")
	flag.Parse()
	apiS.addr = *address

	apiS.echo.GET("/", handlers.AllMetrics(apiS.storage))
	apiS.echo.GET("/value/:typeM/:nameM", handlers.MetricsValue(apiS.storage))
	apiS.echo.POST("/update/:typeM/:nameM/:valueM", handlers.PostWebhandle(apiS.storage))

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
