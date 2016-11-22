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

	runStatus, err := api.CreateRun(params)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	log.Printf("Run %v has been created...\n", runStatus.ID)

	// if foreground flag is enabled we'll monitor run status
	if c.Bool("fg") {
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

	return nil
}

// statusWithError is a helper type used to send RunStatus and error using single channel
type statusWithError struct {
	status rainforest.RunStatus
	err    error
}

func updateRunStatus(c *cli.Context, runID int, t *time.Ticker, resChan chan statusWithError) {
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
func makeRunParams(c *cli.Context) (rainforest.RunParams, error) {
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
	var testIDs []int
	testIDsArg := c.Args().Get(0)
	if testIDsArg != "all" {
		testIDs, err = stringToIntSlice(testIDsArg)
		if err != nil {
			return rainforest.RunParams{}, err
		}
	} else {
		// TODO: Figure out how to do 'all' tests as it's not an integer
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
	var slicedInt []int
	if s == "" {
		return slicedInt, nil
	}
	splitString := strings.Split(s, ",")
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
