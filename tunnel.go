package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/rainforestapp/gonnel"
)

// TunnelConfig contains the configuration to create a tunnel
type TunnelConfig struct {
	siteID int
	port   int
	host   string
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
		LocalAddress: strconv.Itoa(config.port),
		Name:         "adequate",
	})

	client.ConnectAll()

	fmt.Print("Press any to disconnect")
	reader := bufio.NewReader(os.Stdin)
	reader.ReadRune()

	client.DisconnectAll()
}
