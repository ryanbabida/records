package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
)

type JsonStore struct {
	albums []Album
	mu     sync.Mutex
}

func NewJsonStore(filepath string) (*JsonStore, error) {
	recordsFile, err := os.Open("records.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read records json file: %w", err)
	}

	b, _ := io.ReadAll(recordsFile)

	albums := []Album{}
	err = json.Unmarshal(b, &albums)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("failed to parse records json file: %w", err)
	}

	return &JsonStore{albums: albums}, nil
}

func (j *JsonStore) GetAll() ([]Album, error) {
	j.mu.Lock()
	defer j.mu.Unlock()

	return j.albums, nil
}

func (j *JsonStore) GetById(id int) (*Album, error) {
	j.mu.Lock()
	defer j.mu.Unlock()

	for _, album := range j.albums {
		if album.Id == id {
			return &album, nil
		}
	}

	return nil, nil
}
