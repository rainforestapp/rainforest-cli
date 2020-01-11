package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/rainforestapp/gonnel"
)

func newTunnel(c cliContext) {
	pathName := c.Args().First()
	path := pathName

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
	})

	client.ConnectAll()

	fmt.Print("Press any to disconnect")
	reader := bufio.NewReader(os.Stdin)
	reader.ReadRune()

	client.DisconnectAll()
}
