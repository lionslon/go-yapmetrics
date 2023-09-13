package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"
)

type Config struct {
	pollInterval   int    `env:"POLL_INTERVAL"`
	reportInterval int    `env:"REPORT_INTERVAL"`
	addr           string `env:"ADDRESS"`
}

var cfg Config
var valuesGauge = map[string]float64{}
var pollCount uint64

func main() {

	err := getParameters()
	if err != nil {
		log.Fatal(err)
	}

	go getMetrics()

	time.Sleep(time.Duration(cfg.reportInterval) * time.Second)

	for {
		for k, v := range valuesGauge {
			post("gauge", k, strconv.FormatFloat(v, 'f', -1, 64))
		}
		fmt.Println(pollCount)
		post("counter", "PollCount", strconv.FormatUint(pollCount, 10))
		post("gauge", "RandomValue", strconv.FormatFloat(rand.Float64(), 'f', -1, 64))
		pollCount = 0
		time.Sleep(time.Duration(cfg.reportInterval) * time.Second)
	}
}

func getParameters() error {
	flag.StringVar(&cfg.addr, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&cfg.reportInterval, "r", 10, "report interval in seconds")
	flag.IntVar(&cfg.pollInterval, "p", 2, "poll interval in seconds")
	flag.Parse()

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		cfg.addr = envRunAddr
	}
	if envRunAddr := os.Getenv("REPORT_INTERVAL"); envRunAddr != "" {
		cfg.reportInterval, _ = strconv.Atoi(envRunAddr)
	}
	if envRunAddr := os.Getenv("POLL_INTERVAL"); envRunAddr != "" {
		cfg.pollInterval, _ = strconv.Atoi(envRunAddr)
	}
	return nil
}

func getMetrics() {
	var rtm runtime.MemStats

	pollCount += 1
	runtime.ReadMemStats(&rtm)

	valuesGauge["Alloc"] = float64(rtm.Alloc)
	valuesGauge["BuckHashSys"] = float64(rtm.BuckHashSys)
	valuesGauge["Frees"] = float64(rtm.Frees)
	valuesGauge["GCCPUFraction"] = float64(rtm.GCCPUFraction)
	valuesGauge["HeapAlloc"] = float64(rtm.HeapAlloc)
	valuesGauge["HeapIdle"] = float64(rtm.HeapIdle)
	valuesGauge["HeapInuse"] = float64(rtm.HeapInuse)
	valuesGauge["HeapObjects"] = float64(rtm.HeapObjects)
	valuesGauge["HeapReleased"] = float64(rtm.HeapReleased)
	valuesGauge["HeapSys"] = float64(rtm.HeapSys)
	valuesGauge["LastGC"] = float64(rtm.LastGC)
	valuesGauge["Lookups"] = float64(rtm.Lookups)
	valuesGauge["MCacheInuse"] = float64(rtm.MCacheInuse)
	valuesGauge["MCacheSys"] = float64(rtm.MCacheSys)
	valuesGauge["MSpanInuse"] = float64(rtm.MSpanInuse)
	valuesGauge["MSpanSys"] = float64(rtm.MSpanSys)
	valuesGauge["Mallocs"] = float64(rtm.Mallocs)
	valuesGauge["NextGC"] = float64(rtm.NextGC)
	valuesGauge["NumForcedGC"] = float64(rtm.NumForcedGC)
	valuesGauge["NumGC"] = float64(rtm.NumGC)
	valuesGauge["OtherSys"] = float64(rtm.OtherSys)
	valuesGauge["PauseTotalNs"] = float64(rtm.PauseTotalNs)
	valuesGauge["StackInuse"] = float64(rtm.StackInuse)
	valuesGauge["StackSys"] = float64(rtm.StackSys)
	valuesGauge["Sys"] = float64(rtm.Sys)
	valuesGauge["TotalAlloc"] = float64(rtm.TotalAlloc)
}

func post(mType string, mName string, mValue string) {
	resp, err := http.Post(fmt.Sprintf("http://%s/update/%s/%s/%s", cfg.addr, mType, mName, mValue), "text/plain", bytes.NewReader([]byte{}))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}
