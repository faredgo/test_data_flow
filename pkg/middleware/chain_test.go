package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"test_data_flow/pkg/middleware"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChain(t *testing.T) {
	baseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Test-1", "value1")
			next.ServeHTTP(w, r)
		})
	}

	middleware2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Test-2", "value2")
			next.ServeHTTP(w, r)
		})
	}

	chain := middleware.Chain(middleware1, middleware2)(baseHandler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	chain.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	assert.Equal(t, "value1", rec.Header().Get("X-Test-1"))
	assert.Equal(t, "value2", rec.Header().Get("X-Test-2"))
}
