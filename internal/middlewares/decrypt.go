package middlewares

import (
	"bytes"
	"github.com/labstack/echo/v4"
	"github.com/lionslon/go-yapmetrics/internal/services"
	"io"
)

// DecryptBody Расшифровываем тело с помощью приватного ключа
func DecryptBody(keyPath string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) (err error) {
			req := ctx.Request()
			body, err := io.ReadAll(req.Body)
			if err == nil {
				req.Body = io.NopCloser(bytes.NewReader(body))
				return next(ctx)
			}
			decryptedBody := services.TryDecryptOrReturnPlainText(keyPath, body)
			req.Body = io.NopCloser(bytes.NewReader(decryptedBody))
			return next(ctx)
		}
	}
}
