package main

import (
	"errors"
	"fmt"
	"testing"

	"github.com/rainforestapp/rainforest-cli/rainforest"
)

type testBranchAPI struct {
	handleGetBranches  func(params ...string) ([]rainforest.Branch, error)
	handleDeleteBranch func(branchID int) error
	handleMergeBranch  func(branchID int) error
}

func (t *testBranchAPI) GetBranches(params ...string) ([]rainforest.Branch, error) {
	branches, err := t.handleGetBranches(params...)

	return branches, err
}

func (t *testBranchAPI) CreateBranch(branch *rainforest.Branch) error {
	return nil
}

func (t *testBranchAPI) MergeBranch(branchID int) error {
	err := t.handleMergeBranch(branchID)

	return err
}

func (t *testBranchAPI) DeleteBranch(branchID int) error {
	err := t.handleDeleteBranch(branchID)

	return err
}

func TestNewBranch(t *testing.T) {
	context := new(fakeContext)
	testAPI := new(testBranchAPI)

	testCases := []struct {
		branchName    string
		errorExpected bool
		errorMessage  string
	}{
		{"new-branch", false, ""},
		{"", true, "Branch name cannot be blank"},
		{" \n\t ", true, "Branch name cannot be blank"},
	}

	for _, testCase := range testCases {
		context.args = []string{testCase.branchName}

		err := newBranch(context, testAPI)
		errorExpected := testCase.errorExpected
		expectedErrorMessage := testCase.errorMessage

		if errorExpected {
			if err == nil {
				t.Fatal("Expected error, but none occured.")
			}

			errorMessage := err.Error()
			if errorMessage != expectedErrorMessage {
				t.Errorf("newBranch returned error %+v, want %+v", errorMessage, expectedErrorMessage)
			}
		}

		if err != nil && !errorExpected {
			t.Fatal(err.Error())
		}
	}
}

func TestDeleteBranch(t *testing.T) {
	context := new(fakeContext)
	testAPI := new(testBranchAPI)

	testAPI.handleGetBranches = func(params ...string) ([]rainforest.Branch, error) {
		branches := []rainforest.Branch{}
		name := params[0]

		if name != "non-existing-branch" {
			otherBranch := rainforest.Branch{
				ID:   1,
				Name: "also matched " + name,
			}
			branch := rainforest.Branch{
				ID:   2,
				Name: name,
			}

			branches = append(branches, otherBranch, branch)
		}

		return branches, nil
	}

	testAPI.handleDeleteBranch = func(branchID int) (err error) {
		err = nil

		if branchID != 2 {
			err = errors.New(fmt.Sprintf("deleteBranch deleted wrong branch %d, want 2", branchID))
		}

		return err
	}

	testCases := []struct {
		branchName    string
		errorExpected bool
		errorMessage  string
	}{
		{"existing-branch", false, ""},
		{"non-existing-branch", true, "Cannot find branch"},
		{"", true, "Branch name cannot be blank"},
		{" \n\t ", true, "Branch name cannot be blank"},
	}

	for _, testCase := range testCases {
		context.args = []string{testCase.branchName}

		err := deleteBranch(context, testAPI)
		errorExpected := testCase.errorExpected
		expectedErrorMessage := testCase.errorMessage

		if errorExpected {
			if err == nil {
				t.Fatal("Expected error, but none occured.")
			}

			errorMessage := err.Error()
			if errorMessage != expectedErrorMessage {
				t.Errorf("deleteBranch returned error %+v, want %+v", errorMessage, expectedErrorMessage)
			}
		}

		if err != nil && !errorExpected {
			t.Fatal(err.Error())
		}
	}
}

func TestMergeBranch(t *testing.T) {
	context := new(fakeContext)
	testAPI := new(testBranchAPI)

	testAPI.handleGetBranches = func(params ...string) ([]rainforest.Branch, error) {
		branches := []rainforest.Branch{}
		name := params[0]

		if name != "non-existing-branch" {
			branch := rainforest.Branch{
				ID:   1,
				Name: name,
			}
			otherBranch := rainforest.Branch{
				ID:   2,
				Name: "also matched " + name,
			}

			branches = append(branches, branch, otherBranch)
		}

		return branches, nil
	}

	testAPI.handleMergeBranch = func(branchID int) (err error) {
		err = nil

		if branchID != 1 {
			err = errors.New(fmt.Sprintf("mergeBranch merged wrong branch %d, want 1", branchID))
		}

		return err
	}

	testCases := []struct {
		branchName    string
		errorExpected bool
		errorMessage  string
	}{
		{"existing-branch", false, ""},
		{"non-existing-branch", true, "Cannot find branch"},
		{"", true, "Branch name cannot be blank"},
		{" \n\t ", true, "Branch name cannot be blank"},
	}

	for _, testCase := range testCases {
		context.args = []string{testCase.branchName}

		err := mergeBranch(context, testAPI)
		errorExpected := testCase.errorExpected
		expectedErrorMessage := testCase.errorMessage

		if errorExpected {
			if err == nil {
				t.Fatal("Expected error, but none occured.")
			}

			errorMessage := err.Error()
			if errorMessage != expectedErrorMessage {
				t.Errorf("mergeBranch returned error %+v, want %+v", errorMessage, expectedErrorMessage)
			}
		}

		if err != nil && !errorExpected {
			t.Fatal(err.Error())
		}
	}
}
