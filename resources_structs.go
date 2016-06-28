package main

import (
	"strconv"
	"time"
)

func (f foldersResp) displayTable() [][]string {
	table := make([][]string, 0, len(f))
	for _, folders := range f {
		tableRow := make([]string, 2)
		tableRow[0] = strconv.Itoa(folders.ID)
		tableRow[1] = folders.Title
		table = append(table, tableRow)
	}
	return table
}

func (s sitesResp) displayTable() [][]string {
	table := make([][]string, 0, len(s))
	for _, sites := range s {
		tableRow := make([]string, 2)
		tableRow[0] = strconv.Itoa(sites.ID)
		tableRow[1] = sites.Name
		table = append(table, tableRow)
	}
	return table
}

func (b browsersResp) displayTable() [][]string {
	table := make([][]string, 0, len(b.AvailableBrowsers))
	for _, browsers := range b.AvailableBrowsers {
		tableRow := make([]string, 2)
		tableRow[0] = browsers.Name
		tableRow[1] = browsers.Description
		table = append(table, tableRow)
	}
	return table
}

type foldersResp []folder

type folder struct {
	ID        int          `json:"id"`
	CreatedAt time.Time    `json:"created_at"`
	Title     string       `json:"title"`
	Logic     []logicSlice `json:"logic"`
	TestCount int          `json:"test_count"`
}
type logicSlice struct {
	Tag       string `json:"tag"`
	Inclusive bool   `json:"inclusive"`
}

type browsersResp struct {
	AvailableBrowsers []browser `json:"available_browsers"`
}
type browser struct {
	Name           string `json:"name"`
	Description    string `json:"description"`
	Category       string `json:"category"`
	BrowserVersion string `json:"browser_version"`
	OsVersion      string `json:"os_version"`
	Default        bool   `json:"default"`
}
type sitesResp []sites
type sites struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Default   bool      `json:"default"`
}
