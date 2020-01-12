package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/rainforestapp/gonnel"
	"github.com/rainforestapp/rainforest-cli/rainforest"
	"github.com/urfave/cli"
)

type localRunTunnel struct {
	config TunnelConfig
	tunnel *gonnel.Tunnel
}

// NewLocalRun create new local run
func NewLocalRun(c *cli.Context, api *rainforest.Client) error {
	// parse & create tunnel
	// configs := []TunnelConfig{TunnelConfig{port: "3000", extra: map[string]string{"host": "app.rainforest.test"}}}
	// fmt.Println(configs)
	// ngrokURLs := &map[int]string{5537: "https://foo.ngrok.io"}
	tunnelArgs := c.StringSlice("tunnel")
	tunnels := make(map[int]*localRunTunnel)
	for _, conf := range tunnelArgs {
		s := strings.SplitN(conf, "=", 2)
		if len(s) != 2 {
			return fmt.Errorf("Invalid arg for tunnel: %v", conf)
		}

		siteID, err := strconv.Atoi(s[0])
		if err != nil {
			return err
		}

		tunnels[siteID] = &localRunTunnel{config: parseTunnelArgs(s[1])}
	}

	client, err := newTunnelClient()
	if err != nil {
		fmt.Println("ERR", err)
		return err
	}
	for siteID, t := range tunnels {
		tunnels[siteID].tunnel = newTunnel(t.config, client)
	}
	client.ConnectAll()
	defer client.DisconnectAll()

	tunnelURLs := make(map[int]string)
	for siteID, t := range tunnels {
		tunnelURLs[siteID] = t.tunnel.RemoteAddress
	}

	// create temp env
	env, err := api.CreateTemporaryEnvironment("https://www.example.com")
	if err != nil {
		return err
	}

	env, err = api.SetSiteEnvironments(env, &tunnelURLs)
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
	fmt.Println("Tore down environment")

	return err
}
