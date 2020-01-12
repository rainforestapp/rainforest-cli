package main

import (
	"fmt"

	"github.com/rainforestapp/rainforest-cli/rainforest"
	"github.com/urfave/cli"
)

// NewLocalRun create new local run
func NewLocalRun(c *cli.Context, api *rainforest.Client) error {
	// parse & create tunnel
	configs := []TunnelConfig{TunnelConfig{port: "3000", extra: map[string]string{"host": "app.rainforest.test"}}}
	fmt.Println(configs)
	ngrokURLs := &map[int]string{5537: "https://foo.ngrok.io"}

	// create temp env
	env, err := api.CreateTemporaryEnvironment("https://www.example.com")
	fmt.Println(env)
	if err != nil {
		return err
	}

	env, err = api.SetSiteEnvironments(env, ngrokURLs)
	if err == nil {
		// create run
		r := newRunner()
		params, err := r.makeRunParams(c, nil)
		params.EnvironmentID = env.ID
		if err == nil {
			status, err := r.client.CreateRun(params)

			if err == nil {
				err = monitorRunStatus(c, status.ID)
			}
		}
	}

	// teardown temp env
	api.DeleteEnvironment(env.ID)

	return err
}
