package main

import (
	"log"
	"net/http"
)

func main() {
	cfg, err := NewConfig("settings.json")
	if err != nil {
		log.Fatalln(err)
	}

	if cfg == nil || cfg.DataFilePath == nil || cfg.Port == nil {
		log.Fatalln("config was not configured")
	}

	jsonStore, err := NewJsonStore(*cfg.DataFilePath)
	if err != nil {
		log.Fatalln(err)
	}

	h := NewHandlers(cfg, jsonStore)
	m := http.NewServeMux()

	m.Handle("GET /{$}", http.HandlerFunc(h.GetAll))
	m.Handle("GET /{id}", http.HandlerFunc(h.GetById))

	log.Printf("Listening on port :%s\n", *cfg.Port)
	if err = http.ListenAndServe(":"+*cfg.Port, m); err != nil {
		log.Fatalln(err)
	}
}
