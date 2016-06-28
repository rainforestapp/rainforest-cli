package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

var client = new(http.Client)

func checkAPIErr(err error) {
	if err != nil {
		log.Fatalf("API error: %v. Please contact Rainforest support.", err)
	}
}

func addAuthHeaders(req *http.Request) {
	req.Header.Add("CLIENT_TOKEN", apiToken)
	req.Header.Add("Content-Type", "application/json")
}

func postRequest(url string, body []byte) []byte {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	checkAPIErr(err)
	addAuthHeaders(req)
	res, err := client.Do(req)
	checkAPIErr(err)
	data, err := ioutil.ReadAll(res.Body)
	checkAPIErr(err)
	return data
}

func getRequest(url string) []byte {
	req, err := http.NewRequest("GET", baseURL+"/"+url, nil)
	checkAPIErr(err)
	addAuthHeaders(req)
	res, err := client.Do(req)
	checkAPIErr(err)
	data, err := ioutil.ReadAll(res.Body)
	checkAPIErr(err)
	return data
}

func getFolders(url string, resBody *foldersResp) []byte {
	req, err := http.NewRequest("GET", baseURL+"/"+url, nil)
	checkAPIErr(err)
	addAuthHeaders(req)
	res, err := client.Do(req)
	checkAPIErr(err)
	data, err := ioutil.ReadAll(res.Body)
	checkAPIErr(err)
	err = json.Unmarshal(data, &resBody)
	checkAPIErr(err)
	return data
}
