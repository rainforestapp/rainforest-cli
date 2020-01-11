package main

import "github.com/urfave/cli"

import "fmt"

// NewLocalRun create new local run
func NewLocalRun(c *cli.Context) {
	// parse & create tunnel
	configs := []TunnelConfig{TunnelConfig{siteID: 42, port: 3000, host: "app.rainforest.test"}}
	fmt.Println(configs)
	ngrokURLs := map[int]string{42: "https://foo.ngrok.io"}
	fmt.Println(ngrokURLs)
	// create temp env

	// create run
	// teardown temp env
}
