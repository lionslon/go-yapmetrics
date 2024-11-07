package profile

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
)

func ProfileIfEnabled() {
	fmt.Println("started")
	pprofPort := os.Getenv("PPROF_PORT")
	if pprofPort == "" {
		pprofPort = "9090"
	}

	go func() {
		fmt.Println(http.ListenAndServe(fmt.Sprintf(":%s", pprofPort), nil))
	}()

}
