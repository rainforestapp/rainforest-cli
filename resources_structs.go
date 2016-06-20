package main

import "time"

type foldersResp []struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Title     string    `json:"title"`
	Logic     []struct {
		Tag       string `json:"tag"`
		Inclusive bool   `json:"inclusive"`
	} `json:"logic"`
	TestCount int `json:"test_count"`
}

type browsersResp struct {
	AvailableBrowsers []struct {
		Name           string `json:"name"`
		Description    string `json:"description"`
		Category       string `json:"category"`
		BrowserVersion string `json:"browser_version"`
		OsVersion      string `json:"os_version"`
		Default        bool   `json:"default"`
	} `json:"available_browsers"`
}

type clientsResp []struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Default   bool      `json:"default"`
}
