package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type key int

const (
	requestID key = iota
)

type datastore interface {
	GetAll(searchText string) ([]Album, error)
	GetById(id int) (Album, error)
}

type handlers struct {
	config    *Config
	datastore datastore
	logger    *log.Logger
}

func NewHandlers(config *Config, datastore datastore, logger *log.Logger) *handlers {
	return &handlers{config: config, datastore: datastore, logger: logger}
}

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

func WithMiddleware(h http.HandlerFunc, m ...middleware) http.HandlerFunc {
	if len(m) < 1 {
		return h
	}

	wrapped := h
	for i := len(m) - 1; i >= 0; i-- {
		wrapped = m[i](wrapped)
	}

	return wrapped
}

func (h *handlers) GetAll(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	searchText := query.Get("searchText")
	albums, err := h.datastore.GetAll(searchText)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(albums)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *handlers) GetById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.PathValue("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	album, err := h.datastore.GetById(idInt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	empty := Album{}
	if album == empty {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = json.NewEncoder(w).Encode(album)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
