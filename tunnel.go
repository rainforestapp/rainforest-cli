package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/rainforestapp/gonnel"
)

func newTunnel(c cliContext) {
	requestDetails := c.Args().First()
	splitDetails := strings.Split(requestDetails, ",")
	path := splitDetails[0]

	extraOpts := map[string]string{}
	opts := splitDetails[1:]

	for _, o := range opts {
		opt := strings.Split(o, "=")
		extraOpts[opt[0]] = opt[1]
	}

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
	fmt.Printf("I AM HERE %v\n", path)

	client.AddTunnel(&gonnel.Tunnel{
		Proto:        gonnel.HTTP,
		LocalAddress: path,
		Name:         "adequate",
		ExtraOpts:    extraOpts,
	})

	client.ConnectAll()

	fmt.Print("Press any to disconnect")
	reader := bufio.NewReader(os.Stdin)
	reader.ReadRune()

	client.DisconnectAll()
}
