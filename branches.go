package main

import (
	"fmt"
	"strings"

	"github.com/rainforestapp/rainforest-cli/rainforest"
	"github.com/urfave/cli"
)

type branchAPI interface {
	CreateBranch(*rainforest.Branch) error
}

func newBranch(c cliContext, api branchAPI) error {
	name := c.Args().First()
	name = strings.TrimSpace(name)

	if name == "" {
		return cli.NewExitError("Branch name cannot be blank", 1)
	}

	branch := rainforest.Branch{
		Name: name,
	}

	err := api.CreateBranch(&branch)

	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	fmt.Printf("Created branch %q.\n", name)
	return nil
}
