package main

type Artist struct {
	Name string `json:"name"`
}

type Album struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	YearOriginal  int    `json:"yearOriginal"`
	LabelOriginal string `json:"labelOriginal"`
	Year          int    `json:"year"`
	Label         string `json:"label"`
	ArtworkPath   string `json:"artworkPath"`
	AudioPath     string `json:"audioPath"`
	Artist        Artist `json:"artist"`
}
