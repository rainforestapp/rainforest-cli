package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
)

var client = new(http.Client)

func addAuthHeaders(req *http.Request) {
	req.Header.Add("CLIENT_TOKEN", apiToken)
	req.Header.Add("Content-Type", "application/json")
}

func postRequest(url string, body []byte) []byte {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		log.Fatalf("Error when making POST request: %v\n", err)
	}
	addAuthHeaders(req)
	res, _ := client.Do(req)
	data, _ := ioutil.ReadAll(res.Body)
	return data
}

func getRequest(url string) []byte {
	req, err := http.NewRequest("GET", baseURL+"/"+url, nil)
	if err != nil {
		log.Fatalf("Error when making GET request: %v\n", err)
	}
	addAuthHeaders(req)
	res, _ := client.Do(req)
	data, _ := ioutil.ReadAll(res.Body)
	return data
}
