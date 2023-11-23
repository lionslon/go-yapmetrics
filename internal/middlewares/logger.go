package middlewares

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"time"
)

func WithLogging(sugar zap.SugaredLogger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) (err error) {
			//var metric models.Metrics
			start := time.Now()

			req := ctx.Request()
			res := ctx.Response()
			if err = next(ctx); err != nil {
				ctx.Error(err)
			}
			duration := time.Since(start)
			// logging
			sugar.Infoln(
				"uri:", req.RequestURI,
				"method:", req.Method,
				"duration:", duration,
				"status:", res.Status,
				"size:", res.Size,
				//"body:", json.NewDecoder(ctx.Request().Body).Decode(&metric),
			)

			return err
		}
	}
}
