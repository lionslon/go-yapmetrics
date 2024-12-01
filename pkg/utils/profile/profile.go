package profile

import (
	"fmt"
	"go.uber.org/zap"
	"net/http"
	_ "net/http/pprof"
	"os"
)

func StartProfilingServer() {
	pprofPort := os.Getenv("PPROF_PORT")
	if pprofPort == "" {
		pprofPort = "9090"
	}

	go func() {
		zap.S().Info(http.ListenAndServe(fmt.Sprintf(":%s", pprofPort), nil))
	}()

}
