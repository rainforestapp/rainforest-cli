package main

import (
	"fmt"
	"strings"

	"github.com/rainforestapp/rainforest-cli/rainforest"
	"github.com/urfave/cli"
)

type branchAPI interface {
	GetBranches(...string) ([]rainforest.Branch, error)
	CreateBranch(*rainforest.Branch) error
	DeleteBranch(int) error
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

func deleteBranch(c cliContext, api branchAPI) error {
	name := c.Args().First()
	name = strings.TrimSpace(name)

	if name == "" {
		return cli.NewExitError("Branch name cannot be blank", 1)
	}

	branches, err := api.GetBranches(name)

	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	if len(branches) == 0 {
		return cli.NewExitError("Cannot find branch", 1)
	}

	branch := branches[0]

	err = api.DeleteBranch(branch.ID)

	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	fmt.Printf("Deleted branch %q.\n", name)
	return nil
}
