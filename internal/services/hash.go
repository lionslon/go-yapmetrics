package services

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"net/http"
)

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
