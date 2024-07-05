package middlewares

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/labstack/echo/v4"
	"hash"
	"io"
	"net/http"
)

func CheckSignReq(password string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) (err error) {
			req := ctx.Request()
			body, err := io.ReadAll(req.Body)
			//fmt.Println(body)
			if err == nil {
				singPassword := []byte(password)
				bodyHash := GetSign(body, singPassword)
				signR := req.Header.Get("HashSHA256")
				fmt.Println(fmt.Sprintf("body: %s", string(body)))
				fmt.Println(fmt.Sprintf("Пришел: %s", bodyHash))
				fmt.Println(fmt.Sprintf("Рассчитал: %s", signR))
				if signR != bodyHash {
					return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "signature is not valid"})
					//return ctx.String(http.StatusBadRequest, "signature is not valid")
				}
			}
			req.Body = io.NopCloser(bytes.NewReader(body))
			//body1, err := io.ReadAll(req.Body)
			//fmt.Println(body1)
			//req.Body = io.NopCloser(bytes.NewReader(body1))
			return next(ctx)
		}
	}
}

func GetSign(body []byte, pass []byte) string {
	hashValue := hmac.New(sha256.New, pass)
	hashValue.Write(body)
	sum := hashValue.Sum(nil)
	//fmt.Println(sum)
	return hex.EncodeToString(sum)
}

type signResponseWriter struct {
	http.ResponseWriter
	hash hash.Hash
}

func (w signResponseWriter) Write(b []byte) (int, error) {
	w.hash.Write(b)
	w.Header().Set("HashSHA256", hex.EncodeToString(w.hash.Sum(nil)))
	return w.ResponseWriter.Write(b)
}
