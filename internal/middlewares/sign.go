package middlewares

import (
	"bytes"
	"github.com/labstack/echo/v4"
	"github.com/lionslon/go-yapmetrics/internal/services"
	"io"
	"net/http"
)

// CheckSignReq проверяет хеш из заголовков
func CheckSignReq(password string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) (err error) {
			req := ctx.Request()
			signR := req.Header.Get("HashSHA256")
			if signR == "" {
				return next(ctx)
			}
			body, err := io.ReadAll(req.Body)
			if err == nil {
				singPassword := []byte(password)
				bodyHash := services.GetSign(body, singPassword)
				if signR != bodyHash {
					return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "signature is not valid"})
				}
			}
			req.Body = io.NopCloser(bytes.NewReader(body))
			return next(ctx)
		}
	}
}
