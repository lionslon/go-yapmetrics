package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/lionslon/go-yapmetrics/internal/config"
	"github.com/lionslon/go-yapmetrics/internal/middlewares"
	"github.com/lionslon/go-yapmetrics/internal/models"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
	"io"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

var (
	valuesGauge = map[string]float64{}
	pollCount   uint64
	mu          sync.Mutex
)

func main() {

	config.PrintBuildInfo()
	cfg := config.NewClient()
	var wg sync.WaitGroup

	pollTicker := time.NewTicker(time.Duration(cfg.PollInterval) * time.Second)
	defer pollTicker.Stop()
	reportTicker := time.NewTicker(time.Duration(cfg.ReportInterval) * time.Second)
	defer reportTicker.Stop()

	limitChan := make(chan struct{}, cfg.RateLimit)
	wg.Add(1)

	go func() {
		defer wg.Done()
		for range pollTicker.C {
			var metricsWg sync.WaitGroup
			metricsWg.Add(2)
			go func() {
				defer metricsWg.Done()
				getMetrics()
			}()
			go func() {
				defer metricsWg.Done()
				getExtraMetrics()
			}()
			metricsWg.Wait()
		}
	}()

	go func() {
		defer wg.Done()
		for range reportTicker.C {
			limitChan <- struct{}{}
			go func() {
				postQueries(cfg)
				<-limitChan
			}()
		}
	}()
	wg.Wait()
}

func getMetrics() {
	mu.Lock()
	defer mu.Unlock()

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

func getExtraMetrics() {
	mu.Lock()
	defer mu.Unlock()

	vmm, _ := mem.VirtualMemory()
	cpm, _ := cpu.Percent(0, true)

	valuesGauge["TotalMemory"] = float64(vmm.Total)
	valuesGauge["FreeMemory"] = float64(vmm.Free)
	number, _ := cpu.Counts(true)
	for i := 0; i < number; i++ {
		cpuNumber := fmt.Sprintf("CPUutilization%d", i+1)
		valuesGauge[cpuNumber] = cpm[i]
	}

}

func postQueries(cfg *config.ClientConfig) {
	mu.Lock()
	defer mu.Unlock()

	url := fmt.Sprintf("http://%s/update/", cfg.Addr)
	urlBatch := fmt.Sprintf("http://%s/updates/", cfg.Addr)
	client := retryablehttp.NewClient()
	client.RetryMax = 3
	client.RetryWaitMin = time.Second * 1
	client.RetryWaitMax = time.Second * 5

	var payload []models.Metrics

	for k, v := range valuesGauge {
		payload = append(payload, models.Metrics{ID: k, MType: "gauge", Value: &v})
	}
	postJSONBatch(client, urlBatch, payload, cfg.SignPass)
	pc := int64(pollCount)
	err := postJSON(client, url, models.Metrics{ID: "PollCount", MType: "counter", Delta: &pc}, cfg.SignPass)
	if err != nil {
		pollCount = 0
	}
	r := rand.Float64()
	postJSON(client, url, models.Metrics{ID: "RandomValue", MType: "gauge", Value: &r}, cfg.SignPass)
}

func postJSON(c *retryablehttp.Client, url string, m models.Metrics, password string) error {
	js, err := json.Marshal(m)
	if err != nil {
		return err
	}

	gz, err := compress(js)
	if err != nil {
		return err
	}

	req, err := retryablehttp.NewRequest("POST", url, gz)
	if err != nil {
		return err
	}

	singPassword := []byte(password)
	if singPassword != nil {
		req.Header.Add("HashSHA256", middlewares.GetSign(js, singPassword))
	}

	req.Header.Add("content-type", "application/json")
	req.Header.Add("content-encoding", "gzip")
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	req.Body = io.NopCloser(bytes.NewReader(body))

	defer resp.Body.Close()
	return nil
}

func postJSONBatch(c *retryablehttp.Client, url string, m []models.Metrics, password string) error {
	js, err := json.Marshal(m)
	if err != nil {
		return err
	}

	gz, err := compress(js)
	if err != nil {
		return err
	}

	req, err := retryablehttp.NewRequest("POST", url, gz)
	if err != nil {
		return err
	}

	singPassword := []byte(password)
	if singPassword != nil {
		req.Header.Add("HashSHA256", middlewares.GetSign(js, singPassword))
	}

	req.Header.Add("content-type", "application/json")
	req.Header.Add("content-encoding", "gzip")
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	req.Body = io.NopCloser(bytes.NewReader(body))

	defer resp.Body.Close()
	return nil
}

func compress(b []byte) ([]byte, error) {
	var bf bytes.Buffer
	gz, err := gzip.NewWriterLevel(&bf, gzip.BestSpeed)
	if err != nil {
		return nil, err
	}
	_, err = gz.Write(b)
	if err != nil {
		return nil, err
	}
	gz.Close()
	return bf.Bytes(), nil
}
