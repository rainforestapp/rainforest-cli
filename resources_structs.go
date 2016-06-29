package main

import "strconv"

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
	ID    int    `json:"id"`
	Title string `json:"title"`
}

type browsersResp struct {
	AvailableBrowsers []browser `json:"available_browsers"`
}

type browser struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type sitesResp []sites

type sites struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
