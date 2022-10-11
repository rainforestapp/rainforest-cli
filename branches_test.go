package main

import (
	"testing"

	"github.com/rainforestapp/rainforest-cli/rainforest"
)

type testBranchAPI struct {
}

func (t *testBranchAPI) CreateBranch(branch *rainforest.Branch) error {
	return nil
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
