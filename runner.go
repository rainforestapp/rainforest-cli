package main

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/rainforestapp/rainforest-cli/rainforest"
	"github.com/urfave/cli"
)

type runnerAPI interface {
	CreateRun(params rainforest.RunParams) (*rainforest.RunStatus, error)
	CreateTemporaryEnvironment(string) (*rainforest.Environment, error)
	rfmlAPI
}

type runner struct {
	client runnerAPI
}

func startRun(c cliContext) error {
	r := newRunner()
	return r.startRun(c)
}

func newRunner() *runner {
	return &runner{client: api}
}

// startRun starts a new Rainforest run & depending on passed flags monitors its execution
func (r *runner) startRun(c cliContext) error {
	// First check if we even want to crate new run or just monitor the existing one.
	if runIDStr := c.String("reattach"); runIDStr != "" {
		runID, err := strconv.Atoi(runIDStr)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		return monitorRunStatus(c, runID)
	}

	var localTests []*rainforest.RFTest
	var err error
	if c.Bool("f") {
		localTests, err = r.prepareLocalRun(c)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
	}

	params, err := r.makeRunParams(c, localTests)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	if c.Bool("git-trigger") {
		var git gitTrigger
		git, err = newGitTrigger()
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

	runStatus, err := r.client.CreateRun(params)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	log.Printf("Run %v has been created.", runStatus.ID)

	// if background flag is enabled we'll skip monitoring run status
	if c.Bool("bg") {
		return nil
	}

	return monitorRunStatus(c, runStatus.ID)
}

func (r *runner) prepareLocalRun(c cliContext) ([]*rainforest.RFTest, error) {
	tags := getTags(c)
	files := c.Args()
	tests, err := readRFMLFiles(files)
	if err != nil {
		return nil, err
	}

	uploads, err := filterUploadTests(tests, tags)
	if err != nil {
		return nil, err
	}
	err = uploadRFMLFiles(uploads, true, r.client)
	if err != nil {
		return nil, err
	}

	forceExecute := map[string]bool{}
	for _, path := range c.StringSlice("force-execute") {
		abs, err := filepath.Abs(path)
		if err != nil {
			log.Printf("%v is not a valid path", path)
			continue
		}
		forceExecute[abs] = true
	}

	forceSkip := map[string]bool{}
	for _, path := range c.StringSlice("exclude") {
		abs, err := filepath.Abs(path)
		if err != nil {
			log.Printf("%v is not a valid path", path)
			continue
		}
		forceSkip[abs] = true
	}
	return filterExecuteTests(tests, tags, forceExecute, forceSkip), nil
}

// filterUploadTests pre-filters tests for upload. The rule is: upload anything
// with the tag *plus* anything that is depended on by a tagged test.
func filterUploadTests(tests []*rainforest.RFTest, tags []string) ([]*rainforest.RFTest, error) {
	testsByID := map[string]*rainforest.RFTest{}
	for _, test := range tests {
		testsByID[test.RFMLID] = test
	}

	// DFS for filtered tests + embeds
	includedTests := make(map[*rainforest.RFTest]bool)
	var q []*rainforest.RFTest

	// Start with tag-filtered tests
	for _, test := range tests {
		if tags == nil || anyMember(tags, test.Tags) {
			q = append(q, test)
		}
	}
	for len(q) > 0 {
		t := q[len(q)-1]
		q = q[:len(q)-1]
		includedTests[t] = true

		for _, step := range t.Steps {
			if embed, ok := step.(rainforest.RFEmbeddedTest); ok {
				embeddedTest, ok := testsByID[embed.RFMLID]
				if !ok {
					return nil, fmt.Errorf("Could not find embedded test %v", embed.RFMLID)
				}
				if _, ok := includedTests[embeddedTest]; !ok {
					q = append(q, embeddedTest)
				}
			}
		}
	}

	result := make([]*rainforest.RFTest, 0, len(includedTests))
	for t := range includedTests {
		result = append(result, t)
	}

	return result, nil
}

// filterExecuteTests filters for tests that should execute. The rules are: it
// should execute if it's tagged properly *and* has Execute set to true (or is
// in forceExecute) *and* isn't in forceSkip.
func filterExecuteTests(tests []*rainforest.RFTest, tags []string, forceExecute, forceSkip map[string]bool) []*rainforest.RFTest {
	var result []*rainforest.RFTest
	for _, test := range tests {
		path, err := filepath.Abs(test.RFMLPath)
		if err != nil {
			path = ""
		}
		if !forceSkip[path] &&
			(test.Execute || forceExecute[path]) &&
			(tags == nil || anyMember(tags, test.Tags)) {

			result = append(result, test)
		}
	}

	return result
}

func monitorRunStatus(c cliContext, runID int) error {
	backoff := 1

	for {
		status, msg, done, err := getRunStatus(c.Bool("fail-fast"), runID)
		log.Print(msg)

		if done {
			if status.FrontendURL != "" {
				log.Printf("The detailed results are available at %v\n", status.FrontendURL)
			}

			postRunJUnitReport(c, runID)

			if status.Result != "passed" {
				return cli.NewExitError("", 1)
			}

			return nil
		}

		// If we've had too many errors, give up
		if backoff >= 5 {
			msg := fmt.Sprintf("Can not get run status after %d attempts, giving up", backoff)
			return cli.NewExitError(msg, 1)
		}

		// If we hit an error, wait longer before retrying
		if err != nil {
			backoff++
		} else {
			// Reset backoff
			backoff = 1
		}

		log.Printf("Waiting for %s before retrying", runStatusPollInterval)
		time.Sleep(runStatusPollInterval)
	}
}

func getRunStatus(failFast bool, runID int) (*rainforest.RunStatus, string, bool, error) {
	newStatus, err := api.CheckRunStatus(runID)
	if err != nil {
		msg := fmt.Sprintf("API error: %v\n", err)
		return newStatus, msg, false, err
	}

	if newStatus.StateDetails.IsFinalState {
		msg := fmt.Sprintf("Run %v is now %v and has %v\n", runID, newStatus.State, newStatus.Result)
		return newStatus, msg, false, nil
	}

	msg := fmt.Sprintf("Run %v is %v and is %v%% complete\n", runID, newStatus.State, newStatus.CurrentProgress.Percent)
	if newStatus.Result == "failed" && failFast {
		return newStatus, msg, true, nil
	}
	return newStatus, msg, false, nil
}

// makeRunParams parses and validates command line arguments + options
// and makes RunParams struct out of them
func (r *runner) makeRunParams(c cliContext, localTests []*rainforest.RFTest) (rainforest.RunParams, error) {
	var err error
	localOnly := localTests != nil

	var smartFolderID int
	if s := c.String("folder"); !localOnly && s != "" {
		smartFolderID, err = strconv.Atoi(c.String("folder"))
		if err != nil {
			return rainforest.RunParams{}, err
		}
	}

	var runGroupID int
	if s := c.String("run-group-id"); !localOnly && s != "" {
		runGroupID, err = strconv.Atoi(c.String("run-group-id"))
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
	if s := c.String("custom-url"); s != "" {
		var customURL *url.URL
		customURL, err = url.Parse(s)
		if err != nil {
			return rainforest.RunParams{}, err
		}

		if (customURL.Scheme != "http") && (customURL.Scheme != "https") {
			return rainforest.RunParams{}, errors.New("custom URL scheme must be http or https")
		}

		var environment *rainforest.Environment
		environment, err = r.client.CreateTemporaryEnvironment(customURL.String())
		if err != nil {
			return rainforest.RunParams{}, err
		}

		log.Printf("Created temporary environment with name %v", environment.Name)
		environmentID = environment.ID
	} else if s := c.String("environment-id"); s != "" {
		environmentID, err = strconv.Atoi(c.String("environment-id"))
		if err != nil {
			return rainforest.RunParams{}, err
		}
	}

	// Figure out test/RFML IDs
	var testIDs interface{}
	var rfmlIDs []string
	testIDsArgs := c.Args()

	if localOnly {
		for _, t := range localTests {
			rfmlIDs = append(rfmlIDs, t.RFMLID)
		}
	} else if testIDsArgs.First() != "all" && testIDsArgs.First() != "" {
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

	tags := getTags(c)

	return rainforest.RunParams{
		Tests:         testIDs,
		RFMLIDs:       rfmlIDs,
		Tags:          tags,
		SmartFolderID: smartFolderID,
		SiteID:        siteID,
		Crowd:         crowd,
		Conflict:      conflict,
		Browsers:      expandedBrowsers,
		Description:   description,
		EnvironmentID: environmentID,
		RunGroupID:    runGroupID,
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

// getTags get tags from a CLI context. It supports expanding comma-separated
// sublists.
func getTags(c cliContext) []string {
	tags := c.StringSlice("tag")
	return expandStringSlice(tags)
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
