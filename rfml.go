package main

import "os"
import "github.com/urfave/cli"
import "github.com/rainforestapp/rainforest-cli/rainforest"
import "fmt"

func validateRFML(c cliContext) error {
	if path := c.Args().First(); path != "" {
		err := validateSingleRFMLFile(path)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		return nil
	}
	return nil
}

func validateSingleRFMLFile(filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	rfmlReader := rainforest.NewRFMLReader(f)
	test, err := rfmlReader.ReadAll()
	fmt.Printf("Parsed test: %#v", test)
	return err
}
