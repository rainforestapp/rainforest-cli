package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

var client = new(http.Client)

type getter interface {
	getRequest(url string) []byte
}

type getResponse struct{}

func addAuthHeaders(req *http.Request) {
	req.Header.Add("CLIENT_TOKEN", apiToken)
	req.Header.Add("Content-Type", "application/json")
}

func postRequest(url string, body []byte) []byte {
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	addAuthHeaders(req)
	res, _ := client.Do(req)
	data, _ := ioutil.ReadAll(res.Body)
	return data
}

func (r getResponse) getRequest(url string) []byte {
	req, _ := http.NewRequest("GET", url, nil)
	addAuthHeaders(req)
	res, _ := client.Do(req)
	data, _ := ioutil.ReadAll(res.Body)
	return data
}
