package main

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/rainforestapp/rainforest-cli/gittrigger"
	"github.com/rainforestapp/rainforest-cli/rainforest"
	"github.com/urfave/cli"
)

type runnerAPI interface {
	CreateRun(params rainforest.RunParams) (*rainforest.RunStatus, error)
	CreateTemporaryEnvironment(string) (*rainforest.Environment, error)
	CheckRunStatus(int) (*rainforest.RunStatus, error)
	rfmlAPI
}

type runner struct {
	client runnerAPI
}

func startRun(c cliContext) error {
	r := newRunner()
	return r.startRun(c)
}

func rerunRun(c cliContext) error {
	r := newRunner()
	return r.rerunRun(c)
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

	// verify --max-reruns is not used with either --fail-fast or --background
	failFast := c.Bool("fail-fast")
	background := c.Bool("background")
	maxReruns := c.Uint("max-reruns")
	if (maxReruns > 0) && (failFast || background) {
		return cli.NewExitError(
			"You can't use --fail-fast or --background when --max-reruns is greater than 0. "+
				"For the CLI to rerun on failure, it has to wait until completion.",
			1,
		)
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
		git, err := gitTrigger.NewGitTrigger()
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		if !git.CheckTrigger() {
			log.Printf("Git trigger enabled, but %v was not found in latest commit. Exiting...", git.Trigger)
			return nil
		}
		if tags := git.GetTags(); len(tags) > 0 {
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
	r.showRunCreated(runStatus)

	// if background flag is enabled we'll skip monitoring run status
	if c.Bool("bg") {
		return nil
	}

	return monitorRunStatus(c, runStatus.ID)
}

// rerunRun reruns failed tests from a previous Rainforest run & depending on passed flags monitors its execution
func (r *runner) rerunRun(c cliContext) error {
	params, err := r.makeRerunParams(c)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	runStatus, err := r.client.CreateRun(params)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	r.showRunCreated(runStatus)

	// if background flag is enabled we'll skip monitoring run status
	if c.Bool("bg") {
		return nil
	}

	return monitorRunStatus(c, runStatus.ID)
}

func (r *runner) showRunCreated(runStatus *rainforest.RunStatus) {
	log.Printf("Run %v has been created. The detailed results are available at %v", runStatus.ID, runStatus.FrontendURL)
}

func (r *runner) prepareLocalRun(c cliContext) ([]*rainforest.RFTest, error) {
	invalidFilters := []string{"folder", "feature", "run-group", "site"}
	for _, filter := range invalidFilters {
		if c.Int(filter) != 0 || (c.String(filter) != "" && c.String(filter) != "0") {
			return nil, fmt.Errorf("%s cannot be specified with run -f", filter)
		}
	}
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
	failedAttempts := 1

	for {
		status, msg, done, err := getRunStatus(c.Bool("fail-fast"), runID, api)
		log.Print(msg)

		if done {
			postRunJUnitReport(c, runID)

			if status.Result != "passed" {
				rerunAttempt := c.Uint("rerun-attempt")
				remainingReruns := c.Uint("max-reruns") - rerunAttempt
				if remainingReruns > 0 {
					cmd, _ := buildRerunArgs(c, runID)
					path, err := os.Executable()
					if err != nil {
						return cli.NewExitError(err.Error(), 1)
					}

					log.Printf("Rerunning %v, attempt %v", runID, rerunAttempt+1)
					exec_err := syscall.Exec(path, cmd, []string{})
					if exec_err != nil {
						return cli.NewExitError(exec_err.Error(), 1)
					}
				} else {
					return cli.NewExitError("", 1)
				}
			}

			if status.FrontendURL != "" {
				log.Printf("The detailed results are available at %v\n", status.FrontendURL)
			}

			return nil
		}

		// If we've had too many errors, give up
		if failedAttempts >= 5 {
			msg := fmt.Sprintf("Can not get run status after %d attempts, giving up", failedAttempts)
			return cli.NewExitError(msg, 1)
		}

		// If we hit an error, record it
		if err != nil {
			failedAttempts++
		} else {
			// Reset attempts
			failedAttempts = 1
		}

		time.Sleep(runStatusPollInterval)
	}
}

func buildRerunArgs(c cliContext, runID int) ([]string, error) {
	maxReruns := c.Uint("max-reruns")
	rerunAttempt := c.Uint("rerun-attempt")

	cmd := []string{
		"rainforest-cli",
		"rerun", strconv.Itoa(runID),
		"--max-reruns", fmt.Sprint(maxReruns),
		"--rerun-attempt", fmt.Sprint(rerunAttempt + 1),
		"--skip-update", // skip auto-updates for reruns
	}

	if token := c.GlobalString("token"); len(token) > 0 {
		cmd = append(cmd, "--token", token)
	}
	if conflict := c.String("conflict"); len(conflict) > 0 {
		cmd = append(cmd, "--conflict", conflict)
	}
	if junitFile := c.String("junit-file"); len(junitFile) > 0 {
		cmd = append(cmd, "--junit-file", junitFile)
	}

	return cmd, nil
}

func getRunStatus(failFast bool, runID int, client runnerAPI) (*rainforest.RunStatus, string, bool, error) {
	newStatus, err := client.CheckRunStatus(runID)
	if err != nil {
		msg := fmt.Sprintf("API error: %v\n", err)
		return newStatus, msg, false, err
	}

	if newStatus.StateDetails.IsFinalState {
		msg := fmt.Sprintf("Run %v is now %v and has %v (%v failed, %v passed)\n", runID, newStatus.State, newStatus.Result, newStatus.CurrentProgress.Failed, newStatus.CurrentProgress.Passed)
		return newStatus, msg, true, nil
	}

	msg := fmt.Sprintf("Run %v is %v\n", runID, newStatus.State)
	if newStatus.State != "queued" && newStatus.State != "validating" {
		msg = fmt.Sprintf("Run %v is %v and is %v%% complete (%v tests in progress, %v failed, %v passed)\n", runID, newStatus.State, newStatus.CurrentProgress.Percent, (newStatus.CurrentProgress.Total - newStatus.CurrentProgress.Complete), newStatus.CurrentProgress.Failed, newStatus.CurrentProgress.Passed)
	}

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

	var siteID int
	if s := c.String("site"); s != "" {
		siteID, err = strconv.Atoi(c.String("site"))
		if err != nil {
			return rainforest.RunParams{}, err
		}
	}

	var crowd string
	if crowd = c.String("crowd"); crowd != "" && crowd != "default" && crowd != "on_premise_crowd" && crowd != "automation" && crowd != "automation_and_crowd" {
		return rainforest.RunParams{}, errors.New("Invalid crowd option specified")
	}

	var conflict string
	if conflict, err = getConflict(c); err != nil {
		return rainforest.RunParams{}, err
	}

	featureID := c.Int("feature")
	runGroupID := c.Int("run-group")

	browsers := c.StringSlice("browser")
	expandedBrowsers := expandStringSlice(browsers)

	description := c.String("description")
	release := c.String("release")

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
		Release:       release,
		EnvironmentID: environmentID,
		FeatureID:     featureID,
		RunGroupID:    runGroupID,
	}, nil
}

func (r *runner) makeRerunParams(c cliContext) (rainforest.RunParams, error) {
	var err error

	var runID int
	runIDString := c.Args().First()
	if runIDString == "" {
		runIDString = os.Getenv("RAINFOREST_RUN_ID")
	}
	if runIDString == "" {
		return rainforest.RunParams{}, errors.New("Missing run ID")
	}
	runID, err = strconv.Atoi(runIDString)
	if err != nil {
		return rainforest.RunParams{}, errors.New("Invalid run ID specified")
	}

	var conflict string
	if conflict, err = getConflict(c); err != nil {
		return rainforest.RunParams{}, err
	}

	return rainforest.RunParams{
		RunID:    runID,
		Conflict: conflict,
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

// getConflict gets conflict from a CLI context. It returns an error if value isn't allowed
func getConflict(c cliContext) (string, error) {
	var conflict string
	if conflict = c.String("conflict"); conflict != "" && conflict != "abort" && conflict != "abort-all" {
		return "", errors.New("Invalid conflict option specified")
	}

	return conflict, nil
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
