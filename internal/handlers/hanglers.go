package handlers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type StorageUpdater interface {
	UpdateCounter(string, int64)
	UpdateGauge(string, float64)
	GetValue(string, string) (string, int)
	AllMetrics() string
}

type Handler struct {
	stor StorageUpdater
}

func PostWebhandle(s StorageUpdater) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		metricsType := ctx.Param("typeM")
		metricsName := ctx.Param("nameM")
		metricsValue := ctx.Param("valueM")

		if metricsType == "counter" {
			if value, err := strconv.ParseInt(metricsValue, 10, 64); err == nil {
				s.UpdateCounter(metricsName, value)
			} else {
				return ctx.String(http.StatusBadRequest, fmt.Sprintf("%s cannot be converted to an integer", metricsValue))
			}
		} else if metricsType == "gauge" {
			if value, err := strconv.ParseFloat(metricsValue, 64); err == nil {
				s.UpdateGauge(metricsName, value)
			} else {
				return ctx.String(http.StatusBadRequest, fmt.Sprintf("%s cannot be converted to a float", metricsValue))
			}
		} else {
			return ctx.String(http.StatusBadRequest, "Invalid metric type. Can only be 'gauge' or 'counter'")
		}

		ctx.Response().Header().Set("Content-Type", "text/plain; charset=utf-8")
		return ctx.String(http.StatusOK, "")
	}
}

func MetricsValue(s StorageUpdater) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		typeM := ctx.Param("typeM")
		nameM := ctx.Param("nameM")

		val, status := s.GetValue(typeM, nameM)
		err := ctx.String(status, val)
		if err != nil {
			return err
		}

		return nil
	}
}

func AllMetrics(s StorageUpdater) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		err := ctx.String(http.StatusOK, s.AllMetrics())
		if err != nil {
			return err
		}

		return nil
	}
}
