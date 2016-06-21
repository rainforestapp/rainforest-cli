package main

import (
	"strconv"
	"time"
)

type returnTable interface {
	TableSlice() [][]string
}

func (f foldersResp) TableSlice() (idTitleTable [][]string) {
	for _, folderSlice := range f {
		tableRow := make([]string, 2)
		tableRow[0] = strconv.Itoa(folderSlice.ID)
		tableRow[1] = folderSlice.Title
		idTitleTable = append(idTitleTable, tableRow)
	}
	return idTitleTable
}

func (s sitesResp) TableSlice() (idTitleTable [][]string) {
	for _, sitesSlice := range s {
		tableRow := make([]string, 2)
		tableRow[0] = strconv.Itoa(sitesSlice.ID)
		tableRow[1] = sitesSlice.Name
		idTitleTable = append(idTitleTable, tableRow)
	}
	return idTitleTable
}

func (b browsersResp) TableSlice() (idTitleTable [][]string) {
	for _, browserSlice := range b.AvailableBrowsers {
		tableRow := make([]string, 2)
		tableRow[0] = browserSlice.Name
		tableRow[1] = browserSlice.Description
		idTitleTable = append(idTitleTable, tableRow)
	}
	return idTitleTable
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
