package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mark3labs/mcp-go/server"
	"github.com/rainforestapp/rainforest-cli/rainforest"
)

const version = "0.1.0"

func main() {
	token := os.Getenv("RAINFOREST_API_TOKEN")
	if token == "" {
		log.Fatal("RAINFOREST_API_TOKEN environment variable is required")
	}

	client := rainforest.NewClient(token, false)

	s := server.NewMCPServer(
		"rainforest",
		version,
		server.WithToolCapabilities(false),
	)

	registerTools(s, client)

	if err := server.ServeStdio(s); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}
