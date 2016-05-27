package main

import (
	"os"

	"gopkg.in/urfave/cli.v2"
)

func main() {
	app := cli.NewApp()
	app.Name = "Rainforest CLI"
	app.Usage = "Command line utility for Rainforest QA"
	app.Commands = []cli.Command{
		{
			Name:  "run",
			Usage: "Run your tests on Rainforest",
			Action: func(c *cli.Context) error {
				createRun()
				return nil
			},
		},
	}
	app.Run(os.Args)
}
