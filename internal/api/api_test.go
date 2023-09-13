package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWebHandle(t *testing.T) {
	// описываем набор данных: метод запроса, запрос, ожидаемый код ответа
	testCases := []struct {
		method       string
		expectedCode int
		request      string
	}{
		{method: http.MethodGet, request: "/update/gauge/NumGC/1", expectedCode: http.StatusMethodNotAllowed},
		{method: http.MethodPut, request: "/update/gauge/NumGC/1", expectedCode: http.StatusMethodNotAllowed},
		{method: http.MethodDelete, request: "/update/gauge/NumGC/1", expectedCode: http.StatusMethodNotAllowed},
		{method: http.MethodConnect, request: "/update/gauge/NumGC/1", expectedCode: http.StatusMethodNotAllowed},
		{method: http.MethodHead, request: "/update/gauge/NumGC/1", expectedCode: http.StatusMethodNotAllowed},
		{method: http.MethodOptions, request: "/update/gauge/NumGC/1", expectedCode: http.StatusMethodNotAllowed},
		{method: http.MethodPatch, request: "/update/gauge/NumGC/1", expectedCode: http.StatusMethodNotAllowed},
		{method: http.MethodTrace, request: "/update/gauge/NumGC/1", expectedCode: http.StatusMethodNotAllowed},
		{method: http.MethodPost, request: "/update/gauge/NumGC/1", expectedCode: http.StatusOK},
		{method: http.MethodPost, request: "/read/gauge/NumGC/1", expectedCode: http.StatusNotFound},
		{method: http.MethodPost, request: "/update/result/NumGC/1", expectedCode: http.StatusBadRequest},
	}

	for _, tc := range testCases {
		t.Run(tc.method, func(t *testing.T) {
			request := httptest.NewRequest(tc.method, tc.request, nil)
			w := httptest.NewRecorder()
			s := New()
			// вызовем хендлер как обычную функцию, без запуска самого сервера
			h := http.HandlerFunc(s.webhandle())
			h(w, request)
			result := w.Result()
			assert.Equal(t, tc.expectedCode, result.StatusCode, "Код ответа не совпадает с ожидаемым")
			defer result.Body.Close()
		})
	}
}
