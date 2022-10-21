package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/rainforestapp/rainforest-cli/rainforest"
	"github.com/urfave/cli"
)

type branchAPI interface {
	GetBranches(...string) ([]rainforest.Branch, error)
	CreateBranch(*rainforest.Branch) error
	MergeBranch(int) error
	DeleteBranch(int) error
}

func newBranch(c cliContext, api branchAPI) error {
	name, err := getBranchName(c)

	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	branch := rainforest.Branch{
		Name: name,
	}

	err = api.CreateBranch(&branch)

	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	fmt.Printf("Created branch %q.\n", name)
	return nil
}

func mergeBranch(c cliContext, api branchAPI) error {
	name, err := getBranchName(c)

	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	branchID, err := getBranchID(name, api)

	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	err = api.MergeBranch(branchID)

	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	fmt.Printf("Merged branch %q into main.\n", name)
	return nil
}

func deleteBranch(c cliContext, api branchAPI) error {
	name, err := getBranchName(c)

	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	branchID, err := getBranchID(name, api)

	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	err = api.DeleteBranch(branchID)

	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	fmt.Printf("Deleted branch %q.\n", name)
	return nil
}

// getBranchID gets branchID by using the branch name to query the API
func getBranchID(name string, api branchAPI) (int, error) {
	branches, err := api.GetBranches(name)

	if err != nil {
		return 0, err
	}

	if len(branches) == 0 {
		return 0, errors.New("Cannot find branch")
	}

	branch := branches[0]

	return branch.ID, nil
}

// getBranchName gets branchName from the cli
func getBranchName(c cliContext) (string, error) {
	name := c.Args().First()
	name = strings.TrimSpace(name)

	if name == "" {
		return "", errors.New("Branch name cannot be blank")
	}

	return name, nil
}
