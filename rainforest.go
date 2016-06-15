package main

import (
	"net/http"
	"os"

	"gopkg.in/urfave/cli.v2"
)

var apiToken string
var client = new(http.Client)

func main() {
	app := cli.NewApp()
	app.Name = "Rainforest CLI"
	app.Usage = "Command line utility for Rainforest QA"

	app.Flags = []cli.Flag{
		cli.StringSliceFlag{
			Name:  "tags",
			Usage: "Filter tests by tag",
		},
		cli.StringFlag{
			Name:  "token",
			Usage: "Rainforest API token",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "run",
			Usage: "Run your tests on Rainforest",
			Action: func(c *cli.Context) error {
				apiToken = c.String("token")
				createRun(c)
				return nil
			},
		},
	}
	app.Run(os.Args)
}
