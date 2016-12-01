package main

import (
	"errors"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/rainforestapp/rainforest-cli/rainforest"
	"github.com/urfave/cli"
)

// startRun starts a new Rainforest run & depending on passed flags monitors its execution
func startRun(c *cli.Context) error {
	params, err := makeRunParams(c)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	if c.Bool("git-trigger") {
		git, err := newGitTrigger()
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		if !git.checkTrigger() {
			log.Printf("Git trigger enabled, but %v was not found in latest commit. Exiting...", git.Trigger)
			return nil
		}
		if tags := git.getTags(); len(tags) > 0 {
			if len(params.Tags) == 0 {
				log.Print("Found tag list in the commit message, overwriting argument.")
			} else {
				log.Print("Found tag list in the commit message.")
			}
			params.Tags = tags
		}
	}

	err = preRunCSVUpload(c, api)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	runStatus, err := api.CreateRun(params)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	log.Printf("Run %v has been created...\n", runStatus.ID)

	// if background flag is enabled we'll skip monitoring run status
	if c.Bool("bg") {
		return nil
	}

	// Create two channels to communicate with the polling goroutine
	// One that will tick when it's time to poll and the other to gather final state
	t := time.NewTicker(runStatusPollInterval)
	statusChan := make(chan statusWithError, 0)
	go updateRunStatus(c, runStatus.ID, t, statusChan)

	// This channel readout will block until updateRunStatus pushed final result to it
	finalState := <-statusChan
	if finalState.err != nil {
		return cli.NewExitError(finalState.err.Error(), 1)
	}

	if finalState.status.FrontendURL != "" {
		log.Printf("The detailed results are available at %v\n", finalState.status.FrontendURL)
	}

	if finalState.status.Result != "passed" {
		return cli.NewExitError("", 1)
	}

	return nil
}

// statusWithError is a helper type used to send RunStatus and error using single channel
type statusWithError struct {
	status *rainforest.RunStatus
	err    error
}

func updateRunStatus(c cliContext, runID int, t *time.Ticker, resChan chan statusWithError) {
	for {
		// Wait for tick
		<-t.C
		newStatus, err := api.CheckRunStatus(runID)
		if err != nil {
			resChan <- statusWithError{status: newStatus, err: err}
		}

		isFinalState := newStatus.StateDetails.IsFinalState
		state := newStatus.State
		currentPercent := newStatus.CurrentProgress.Percent

		if !isFinalState {
			log.Printf("Run %v is %v and is %v%% complete\n", runID, state, currentPercent)
			if newStatus.Result == "failed" && c.Bool("fail-fast") {
				resChan <- statusWithError{status: newStatus, err: nil}
			}
		} else {
			log.Printf("Run %v is now %v and has %v\n", runID, state, newStatus.Result)
			resChan <- statusWithError{status: newStatus, err: nil}
		}
	}
}

// makeRunParams parses and validates command line arguments + options
// and makes RunParams struct out of them
func makeRunParams(c cliContext) (rainforest.RunParams, error) {
	var err error

	var smartFolderID int
	if s := c.String("folder"); s != "" {
		smartFolderID, err = strconv.Atoi(c.String("folder"))
		if err != nil {
			return rainforest.RunParams{}, err
		}
	}

	var siteID int
	if s := c.String("site"); s != "" {
		siteID, err = strconv.Atoi(c.String("site"))
		if err != nil {
			return rainforest.RunParams{}, err
		}
	}

	var crowd string
	if crowd = c.String("crowd"); crowd != "" && crowd != "default" && crowd != "on_premise_crowd" {
		return rainforest.RunParams{}, errors.New("Invalid crowd option specified")
	}

	var conflict string
	if conflict = c.String("conflict"); conflict != "" && conflict != "abort" && conflict != "abort-all" {
		return rainforest.RunParams{}, errors.New("Invalid conflict option specified")
	}

	browsers := c.StringSlice("browser")
	expandedBrowsers := expandStringSlice(browsers)

	description := c.String("description")

	var environmentID int
	if s := c.String("environment-id"); s != "" {
		environmentID, err = strconv.Atoi(c.String("environment-id"))
		if err != nil {
			return rainforest.RunParams{}, err
		}
	}

	// Parse command argument as a list of test IDs
	var testIDs interface{}
	testIDsArgs := c.Args()
	if testIDsArgs.First() != "all" && testIDsArgs.First() != "" {
		testIDs = []int{}
		for _, arg := range testIDsArgs {
			nextTestIDs, err := stringToIntSlice(arg)
			if err != nil {
				return rainforest.RunParams{}, err
			}
			testIDs = append(testIDs.([]int), nextTestIDs...)
		}
	} else if testIDsArgs.First() == "all" {
		testIDs = "all"
	}

	// We get tags slice from arguments and then expand comma separated lists into separate entries
	tags := c.StringSlice("tag")
	expandedTags := expandStringSlice(tags)

	return rainforest.RunParams{
		Tests:         testIDs,
		Tags:          expandedTags,
		SmartFolderID: smartFolderID,
		SiteID:        siteID,
		Crowd:         crowd,
		Conflict:      conflict,
		Browsers:      expandedBrowsers,
		Description:   description,
		EnvironmentID: environmentID,
	}, nil
}

// stringToIntSlice takes a string of comma separated integers and returns a slice of them
func stringToIntSlice(s string) ([]int, error) {
	if s == "" {
		return nil, nil
	}
	splitString := strings.Split(s, ",")
	var slicedInt []int
	for _, slice := range splitString {
		newInt, err := strconv.Atoi(strings.TrimSpace(slice))
		if err != nil {
			return slicedInt, err
		}
		slicedInt = append(slicedInt, newInt)
	}
	return slicedInt, nil
}

// expandStringSlice takes a slice of strings and expands any comma separated sublists
// into one slice. This allows us to accept args like: -tag abc -tag qwe,xyz
func expandStringSlice(slice []string) []string {
	var result []string
	for _, element := range slice {
		splitElement := strings.Split(element, ",")
		for _, singleElement := range splitElement {
			result = append(result, strings.TrimSpace(singleElement))
		}
	}
	return result
}
