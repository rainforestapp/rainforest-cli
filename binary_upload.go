package main

import (
	"os"

	"github.com/urfave/cli"
)

func binaryUpload(c cliContext) error {
	// Get the file path from the command argument
	filePath := c.Args().First()
	if filePath == "" {
		return cli.NewExitError("File not specified", 1)
	}

	err := verifyBinary(filePath)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	err = api.UploadBinary(filePath)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	return nil
}

func verifyBinary(filePath string) error {
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0666)
	file.Close()
	return err
}
