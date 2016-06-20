package main

import (
	"strconv"
	"time"
)

type returnTable interface {
	TableSlice() ([]string, []string)
}

func (f foldersResp) TableSlice() (idArray []string, titleArray []string) {
	for _, folderSlice := range f {
		idArray = append(idArray, strconv.Itoa(folderSlice.ID))
		titleArray = append(titleArray, folderSlice.Title)
	}
	return idArray, titleArray
}

func (s sitesResp) TableSlice() (idArray []string, titleArray []string) {
	for _, sitesSlice := range s {
		idArray = append(idArray, strconv.Itoa(sitesSlice.ID))
		titleArray = append(titleArray, sitesSlice.Name)
	}
	return idArray, titleArray
}

func (b browsersResp) TableSlice() (idArray []string, titleArray []string) {
	for _, browserSlice := range b.AvailableBrowsers {
		idArray = append(idArray, browserSlice.Name)
		titleArray = append(titleArray, browserSlice.Description)
	}
	return idArray, titleArray
}

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

type sitesResp []struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Default   bool      `json:"default"`
}
