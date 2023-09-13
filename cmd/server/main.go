package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type gauge float64
type counter int64

type MemStorage struct {
	gaugeData   map[string]gauge
	counterData map[string]counter
}

var storage = MemStorage{
	gaugeData:   make(map[string]gauge),
	counterData: make(map[string]counter),
}

// функция main вызывается автоматически при запуске приложения
func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	return http.ListenAndServe(`:8080`, http.HandlerFunc(webhandle))
}

func webhandle(w http.ResponseWriter, r *http.Request) {
	// разрешаем только POST-запросы
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// разбиваем запрос и проверяем соответствии формату
	reqURL := strings.Split(r.URL.Path, "/")
	if len(reqURL) != 5 || reqURL[1] != "update" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// обработка запроса
	metricsType := reqURL[2]
	metricsName := reqURL[3]
	metricsValue := reqURL[4]
	if metricsType == "counter" {
		if value, err := strconv.ParseInt(metricsValue, 10, 64); err == nil {
			storage.counterData[metricsName] += counter(value)
			print(fmt.Sprintf("%v: %v (%v)\r\n", metricsName, storage.counterData[metricsName], metricsType))
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	} else if metricsType == "gauge" {
		if value, err := strconv.ParseFloat(metricsValue, 64); err == nil {
			storage.gaugeData[metricsName] = gauge(value)
			print(fmt.Sprintf("%v: %v (%v)\r\n", metricsName, storage.gaugeData[metricsName], metricsType))
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

}
