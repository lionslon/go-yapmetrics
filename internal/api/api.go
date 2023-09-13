package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/lionslon/go-yapmetrics/internal/storage"
)

type APIServer struct {
	storage *storage.MemStorage
}

func New() *APIServer {
	return &APIServer{storage.New()}
}

func (s *APIServer) Start() error {
	return http.ListenAndServe(`:8080`, http.HandlerFunc(s.webhandle()))
}

func (s *APIServer) webhandle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// разрешаем только POST-запросы
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		//println(r.URL.Path)
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
				s.storage.UpdateCounter(metricsName, value)
				w.WriteHeader(http.StatusOK)
			} else {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		} else if metricsType == "gauge" {
			if value, err := strconv.ParseFloat(metricsValue, 64); err == nil {
				s.storage.UpdateGauge(metricsName, value)
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
}
