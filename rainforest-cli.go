package main

import (
	"os"

	"gopkg.in/urfave/cli.v2"
)

var apiToken string

type resourceGetter interface {
	getFolders() [][]string
	getSites() [][]string
	getBrowsers() [][]string
}
type webResGetter struct{}

func main() {
	var web webResGetter
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

		{
			Name:  "folders",
			Usage: "Retreive folders on Rainforest",
			Action: func(c *cli.Context) error {
				apiToken = c.String("token")
				fetchResource("Folders", web)
				return nil
			},
		},

		{
			Name:  "sites",
			Usage: "Retreive sites on Rainforest",
			Action: func(c *cli.Context) error {
				apiToken = c.String("token")
				fetchResource("Sites", web)
				return nil
			},
		},

		{
			Name:  "browsers",
			Usage: "Retreive sites on Rainforest",
			Action: func(c *cli.Context) error {
				apiToken = c.String("token")
				fetchResource("Browsers", web)
				return nil
			},
		},
	}
	app.Run(os.Args)
}
