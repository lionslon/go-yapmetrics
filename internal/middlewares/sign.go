package middlewares

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/labstack/echo/v4"
	"hash"
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
				bodyHash := GetSign(body, singPassword)
				if signR != bodyHash {
					return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "signature is not valid"})
				}
			}
			req.Body = io.NopCloser(bytes.NewReader(body))
			return next(ctx)
		}
	}
}

// GetSign формирует хеш-подпись из фразы-пароля
func GetSign(body []byte, pass []byte) string {
	hashValue := hmac.New(sha256.New, pass)
	hashValue.Write(body)
	sum := hashValue.Sum(nil)
	return hex.EncodeToString(sum)
}

// signResponseWriter реализует интерфейс http.ResponseWriter и позволяет прозрачно для сервера
// работать с хешом запросов
type signResponseWriter struct {
	http.ResponseWriter
	hash hash.Hash
}

// Write вставляет хеш в заголовок
func (w signResponseWriter) Write(b []byte) (int, error) {
	w.hash.Write(b)
	w.Header().Set("HashSHA256", hex.EncodeToString(w.hash.Sum(nil)))
	return w.ResponseWriter.Write(b)
}
