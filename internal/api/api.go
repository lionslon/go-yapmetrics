package api

import (
	"flag"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/lionslon/go-yapmetrics/internal/handlers"
	"github.com/lionslon/go-yapmetrics/internal/storage"
	"log"
	"os"
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
	handler := handlers.New()
	var address string
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		address = envRunAddr
	} else {
		flag.StringVar(&address, "a", "localhost:8080", "address and port to run server")
		flag.Parse()
	}
	apiS.addr = address

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
