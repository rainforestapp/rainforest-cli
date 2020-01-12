package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/rainforestapp/gonnel"
)

// TunnelConfig contains the configuration to create a tunnel
type TunnelConfig struct {
	port  string
	extra map[string]string
}

func splitTunnelArgs(requestDetails string) TunnelConfig {
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

func newTunnel(config TunnelConfig) {
	client, err := gonnel.NewClient(gonnel.Options{
		BinaryPath: "ngrok",
	})
	client.LogApi = true
	if err != nil {
		fmt.Println(err)
	}
	defer client.Close()

	done := make(chan bool)
	fmt.Println("FOOOO")
	go client.StartServer(done)
	<-done
	fmt.Printf("I AM HERE %v\n", config.port)

	client.AddTunnel(&gonnel.Tunnel{
		Proto:        gonnel.HTTP,
		LocalAddress: config.port,
		Name:         "adequate",
		ExtraOpts:    config.extra,
	})

	client.ConnectAll()

	fmt.Print("Press any to disconnect")
	reader := bufio.NewReader(os.Stdin)
	reader.ReadRune()

	client.DisconnectAll()
}
