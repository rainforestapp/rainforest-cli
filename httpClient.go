package main

import "net/http"

func addAuthHeaders(req *http.Request) {
	req.Header.Add("CLIENT_TOKEN", apiToken)
	req.Header.Add("Content-Type", "application/json")
}
