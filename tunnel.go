package main

import (
	"os"
	"os/signal"
	"strings"

	"github.com/rainforestapp/gonnel"
	"github.com/urfave/cli"
)

// TunnelConfig contains the configuration to create a tunnel
type TunnelConfig struct {
	port  string
	extra map[string]string
}

func parseTunnelArgs(requestDetails string) TunnelConfig {
	splitDetails := strings.Split(requestDetails, ",")
	port := splitDetails[0]

	extraOpts := map[string]string{}
	opts := splitDetails[1:]

	for _, o := range opts {
		opt := strings.Split(o, "=")
		extraOpts[opt[0]] = opt[1]
	}

	return TunnelConfig{port: port, extra: extraOpts}
}

func startTunnel(c *cli.Context) error {
	config := parseTunnelArgs(c.Args().First())
	client, err := newTunnelClient()
	if err != nil {
		return err
	}
	defer client.Close()
	newTunnel(config, client)
	client.ConnectAll()

	done := make(chan os.Signal, 1)
	signal.Notify(done)
	<-done
	client.DisconnectAll()

	return nil
}

func newTunnelClient() (*gonnel.Client, error) {
	return gonnel.NewClient(gonnel.Options{
		BinaryPath: "ngrok",
	})
}

func newTunnel(config TunnelConfig, client *gonnel.Client) *gonnel.Tunnel {
	done := make(chan bool)
	go client.StartServer(done)
	<-done

	t := &gonnel.Tunnel{
		Proto:        gonnel.HTTP,
		LocalAddress: config.port,
		Name:         "adequate",
		ExtraOpts:    config.extra,
	}
	client.AddTunnel(t)

	return t
}
