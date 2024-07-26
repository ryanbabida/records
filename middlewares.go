package main

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (h *handlers) requestMiddleware(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), requestID, uuid.New())
		f(w, r.WithContext(ctx))
	}
}

func (h *handlers) jsonMiddleware(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		f(w, r)
	}
}

func (h *handlers) perfMiddleware(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestId := r.Context().Value(requestID)
		t := time.Now()
		f(w, r)
		h.logger.Printf("%s - %d ms\n", requestId, time.Since(t).Milliseconds())
	}
}

type middleware func(f http.HandlerFunc) http.HandlerFunc

func withMiddleware(h http.HandlerFunc, m ...middleware) http.HandlerFunc {
	if len(m) < 1 {
		return h
	}

	wrapped := h
	for i := len(m) - 1; i >= 0; i-- {
		wrapped = m[i](wrapped)
	}

	return wrapped
}
