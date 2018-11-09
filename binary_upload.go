package main

import (
  "github.com/rainforestapp/rainforest-cli/rainforest"
	"github.com/urfave/cli"
	"log"
	"os"
)

func binaryUpload(c cliContext) error {
	// Get the file path either from the command argument
	fileName := c.Args().First()
	if fileName == "" {
		return cli.NewExitError("File not specified", 1)
	}
	log.Println(fileName)
  log.Println(rainforest.GetUploadEndpoint(fileName))

	file, err := os.OpenFile(fileName, os.O_RDONLY, 0666)
	if err != nil {
		if os.IsPermission(err) {
			log.Println("Unable to read from ", fileName)
			log.Println(err)
			os.Exit(1)
		}
	}
	file.Close()

	return nil
}