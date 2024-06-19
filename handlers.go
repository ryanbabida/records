package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type datastore interface {
	GetAll(searchText string) ([]Album, error)
	GetById(id int) (Album, error)
}

type handlers struct {
	config    *Config
	datastore datastore
}

func NewHandlers(config *Config, datastore datastore) *handlers {
	return &handlers{config: config, datastore: datastore}
}

func (h *handlers) GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
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
