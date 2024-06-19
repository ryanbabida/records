package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"
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

func (j *JsonStore) GetAll(searchText string) ([]Album, error) {
	j.mu.Lock()
	defer j.mu.Unlock()

	albums := make([]Album, len(j.albums))
	copy(albums, j.albums)

	s := strings.TrimSpace(searchText)
	s = strings.ToLower(s)

	if s != "" {
		albums = slices.DeleteFunc(albums, func(album Album) bool {
			yearOriginal := strconv.Itoa(album.YearOriginal)
			return !(strings.Contains(strings.ToLower(album.Name), s) ||
				strings.Contains(strings.ToLower(album.Artist.Name), s) ||
				strings.Contains(strings.ToLower(album.LabelOriginal), s) ||
				strings.Contains(strings.ToLower(album.Label), s) ||
				s == yearOriginal)
		})
	}

	sort.Slice(albums, func(i int, j int) bool {
		return albums[i].Name < albums[j].Name
	})

	return albums, nil
}

func (j *JsonStore) GetById(id int) (Album, error) {
	j.mu.Lock()
	defer j.mu.Unlock()

	for _, album := range j.albums {
		if album.Id == id {
			return album, nil
		}
	}

	return Album{}, nil
}
