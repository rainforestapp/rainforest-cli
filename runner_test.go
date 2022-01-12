package main

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"sync"
	"testing"

	"github.com/rainforestapp/rainforest-cli/rainforest"
	"github.com/urfave/cli"
)

func TestStringToIntSlice(t *testing.T) {
	var testCases = []struct {
		commaSepList string
		want         []int
	}{
		{
			commaSepList: "123,456,789",
			want:         []int{123, 456, 789},
		},
		{
			commaSepList: "123, 456, 789",
			want:         []int{123, 456, 789},
		},
		{
			commaSepList: "123",
			want:         []int{123},
		},
	}

	for _, tCase := range testCases {
		got, _ := stringToIntSlice(tCase.commaSepList)
		if !reflect.DeepEqual(tCase.want, got) {
			t.Errorf("stringToIntSlice returned %+v, want %+v", got, tCase.want)
		}
	}
}

func TestExpandStringSlice(t *testing.T) {
	var testCases = []struct {
		stringSlice []string
		want        []string
	}{
		{
			stringSlice: []string{"foo,bar,baz"},
			want:        []string{"foo", "bar", "baz"},
		},
		{
			stringSlice: []string{"foo", "bar,baz"},
			want:        []string{"foo", "bar", "baz"},
		},
		{
			stringSlice: []string{"foo, bar", "baz"},
			want:        []string{"foo", "bar", "baz"},
		},
		{
			stringSlice: []string{"foo"},
			want:        []string{"foo"},
		},
	}

	for _, tCase := range testCases {
		got := expandStringSlice(tCase.stringSlice)
		if !reflect.DeepEqual(tCase.want, got) {
			t.Errorf("expandStringSlice returned %+v, want %+v", got, tCase.want)
		}
	}
}

type fakeRunnerClient struct {
	runStatuses []rainforest.RunStatus
	environment rainforest.Environment
	// runParams captures whatever params were sent
	runParams rainforest.RunParams
	// createdTests captures which tests were created
	createdTests []*rainforest.RFTest
	// got some potential race conditions!
	mu sync.Mutex
	// "inherit" from RFML API
	testRfmlAPI
}

func (r *fakeRunnerClient) CheckRunStatus(runID int) (*rainforest.RunStatus, error) {
	for _, runStatus := range r.runStatuses {
		if runStatus.ID == runID {
			return &runStatus, nil
		}
	}

	return nil, fmt.Errorf("Unable to find run status for run ID %v", runID)
}

func (r *fakeRunnerClient) GetTestIDs() ([]rainforest.TestIDPair, error) {
	pairs := make([]rainforest.TestIDPair, len(r.createdTests))
	for idx, test := range r.createdTests {
		pairs[idx] = rainforest.TestIDPair{ID: test.TestID, RFMLID: test.RFMLID}
	}

	return pairs, nil
}

func (r *fakeRunnerClient) CreateTest(t *rainforest.RFTest) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	t.TestID = len(r.createdTests)
	r.createdTests = append(r.createdTests, t)
	return nil
}

func (r *fakeRunnerClient) UpdateTest(t *rainforest.RFTest) error {
	// meh
	return nil
}

func (r *fakeRunnerClient) CreateTemporaryEnvironment(s string) (*rainforest.Environment, error) {
	return &r.environment, nil
}

func (r *fakeRunnerClient) CreateRun(p rainforest.RunParams) (*rainforest.RunStatus, error) {
	r.runParams = p
	return &rainforest.RunStatus{}, nil
}

func TestGetRunStatus(t *testing.T) {
	runID := 123
	client := new(fakeRunnerClient)

	// Run is in a final state
	client.runStatuses = []rainforest.RunStatus{
		{
			ID: runID,
			StateDetails: struct {
				Name         string `json:"name"`
				IsFinalState bool   `json:"is_final_state"`
			}{"", true},
		},
	}

	_, _, done, err := getRunStatus(false, runID, client)
	if err != nil {
		t.Error(err.Error())
	}
	if !done {
		t.Errorf("Expected \"done\" to be true, got %v", done)
	}

	// Run is failed and failfast is true
	client.runStatuses = []rainforest.RunStatus{
		{
			ID:     runID,
			Result: "failed",
			StateDetails: struct {
				Name         string `json:"name"`
				IsFinalState bool   `json:"is_final_state"`
			}{"", false},
		},
	}

	_, _, done, err = getRunStatus(true, runID, client)
	if err != nil {
		t.Error(err.Error())
	}
	if !done {
		t.Errorf("Expected \"done\" to be true, got %v", done)
	}
}

func TestMakeRunParams(t *testing.T) {
	fakeEnvID := 445566

	var testCases = []struct {
		mappings map[string]interface{}
		args     cli.Args
		expected rainforest.RunParams
	}{
		{
			mappings: make(map[string]interface{}),
			args:     cli.Args{},
			expected: rainforest.RunParams{},
		},
		{
			mappings: map[string]interface{}{
				"folder":         "123",
				"site":           "456",
				"crowd":          "on_premise_crowd",
				"conflict":       "abort",
				"browser":        []string{"chrome", "firefox,safari"},
				"description":    "my awesome description",
				"release":        "1a2b3c",
				"environment-id": "1337",
				"tag":            []string{"tag", "tag2,tag3"},
			},
			args: cli.Args{"12", "34", "56, 78"},
			expected: rainforest.RunParams{
				SmartFolderID: 123,
				SiteID:        456,
				Crowd:         "on_premise_crowd",
				Conflict:      "abort",
				Browsers:      []string{"chrome", "firefox", "safari"},
				Description:   "my awesome description",
				Release:       "1a2b3c",
				EnvironmentID: 1337,
				Tags:          []string{"tag", "tag2", "tag3"},
				Tests:         []int{12, 34, 56, 78},
			},
		},
		{
			mappings: map[string]interface{}{
				"custom-url": "https://www.rainforestqa.com",
			},
			args: cli.Args{},
			expected: rainforest.RunParams{
				EnvironmentID: fakeEnvID,
			},
		},
		{
			mappings: map[string]interface{}{
				"feature": 123,
			},
			args: cli.Args{},
			expected: rainforest.RunParams{
				FeatureID: 123,
			},
		},
		{
			mappings: map[string]interface{}{
				"run-group": 75,
			},
			expected: rainforest.RunParams{
				RunGroupID: 75,
			},
		},
		{
			mappings: map[string]interface{}{
				"automation-max-retries": 2,
			},
			expected: rainforest.RunParams{
				AutomationMaxRetries: 2,
			},
		},
		{
			mappings: map[string]interface{}{
				"crowd": "automation",
			},
			expected: rainforest.RunParams{
				Crowd: "automation",
			},
		},
		{
			mappings: map[string]interface{}{
				"crowd": "automation_and_crowd",
			},
			expected: rainforest.RunParams{
				Crowd: "automation_and_crowd",
			},
		},
	}

	for _, testCase := range testCases {
		c := newFakeContext(testCase.mappings, testCase.args)
		r := newRunner()
		fakeEnv := rainforest.Environment{ID: fakeEnvID, Name: "the foo environment"}
		r.client = &fakeRunnerClient{environment: fakeEnv}
		res, err := r.makeRunParams(c, nil)

		if err != nil {
			t.Errorf("Error trying to create params: %v", err)
		} else if !reflect.DeepEqual(res, testCase.expected) {
			t.Errorf("Incorrect resulting run params.\nActual: %#v\nExpected: %#v", res, testCase.expected)
		}
	}
}

func TestMakeRerunParams(t *testing.T) {
	envRunID, isSet := os.LookupEnv("RAINFOREST_RUN_ID")
	if isSet {
		defer os.Setenv("RAINFOREST_RUN_ID", envRunID)
	} else {
		defer os.Unsetenv("RAINFOREST_RUN_ID")
	}
	os.Setenv("RAINFOREST_RUN_ID", "117")

	var testCases = []struct {
		mappings map[string]interface{}
		args     cli.Args
		expected rainforest.RunParams
	}{
		{
			mappings: make(map[string]interface{}),
			args:     cli.Args{},
			expected: rainforest.RunParams{
				RunID: 117,
			},
		},
		{
			mappings: make(map[string]interface{}),
			args:     cli.Args{"41"},
			expected: rainforest.RunParams{
				RunID: 41,
			},
		},
		{
			mappings: map[string]interface{}{
				"conflict": "abort",
			},
			args: cli.Args{"82"},
			expected: rainforest.RunParams{
				RunID:    82,
				Conflict: "abort",
			},
		},
	}

	for _, testCase := range testCases {
		c := newFakeContext(testCase.mappings, testCase.args)
		r := newRunner()
		fakeEnv := rainforest.Environment{ID: 2401, Name: "the foo environment"}
		r.client = &fakeRunnerClient{environment: fakeEnv}
		res, err := r.makeRerunParams(c)

		if err != nil {
			t.Errorf("Error trying to create params: %v", err)
		} else if !reflect.DeepEqual(res, testCase.expected) {
			t.Errorf("Incorrect resulting run params.\nActual: %#v\nExpected: %#v", res, testCase.expected)
		}
	}
}

func TestStartLocalRun(t *testing.T) {
	rfmlDir := setupTestRFMLDir()
	defer os.RemoveAll(rfmlDir)

	testCases := []struct {
		mappings    map[string]interface{}
		args        cli.Args
		wantError   bool
		wantUpload  []string
		wantExecute []string
	}{
		{
			mappings: map[string]interface{}{
				"f":   true,
				"tag": []string{"foo", "bar"},
				// There's less to stub with bg
				"bg": true,
			},
			args: cli.Args{filepath.Join(rfmlDir, "a"), filepath.Join(rfmlDir, "b/b")},
			// a1 depends on b4 and b4 depends on b5, so b4 and b5 are uploaded
			// even though they're not tagged properly.
			wantUpload: []string{"a1", "a3", "b3", "b4", "b5"},
			// b3 is execute: false so we shouldn't run it
			wantExecute: []string{"a1", "a3"},
		},
		{
			mappings: map[string]interface{}{
				"f":  true,
				"bg": true,
			},
			args:      cli.Args{filepath.Join(rfmlDir, "a")},
			wantError: true,
		},
		{
			mappings: map[string]interface{}{
				"f":       true,
				"feature": 42,
				"bg":      true,
			},
			args:      cli.Args{filepath.Join(rfmlDir, "a"), filepath.Join(rfmlDir, "b/b")},
			wantError: true,
		},
		{
			mappings: map[string]interface{}{
				"f":         true,
				"run-group": 42,
				"bg":        true,
			},
			args:      cli.Args{filepath.Join(rfmlDir, "a"), filepath.Join(rfmlDir, "b/b")},
			wantError: true,
		},
		{
			mappings: map[string]interface{}{
				"f":    true,
				"site": "42",
				"bg":   true,
			},
			args:      cli.Args{filepath.Join(rfmlDir, "a"), filepath.Join(rfmlDir, "b/b")},
			wantError: true,
		},
		{

			mappings: map[string]interface{}{
				"f":             true,
				"bg":            true,
				"force-execute": []string{filepath.Join(rfmlDir, "b/b/b3.rfml")},
				"exclude":       []string{filepath.Join(rfmlDir, "a/a2.rfml")},
			},
			args: cli.Args{
				filepath.Join(rfmlDir, "a/a2.rfml"),
				filepath.Join(rfmlDir, "a/a3.rfml"),
				filepath.Join(rfmlDir, "b/b/b3.rfml"),
			},
			wantUpload: []string{"a2", "a3", "b3"},
			// We don't want to execute a2 but we override execute:false on b3
			wantExecute: []string{"a3", "b3"},
		},
	}

	for _, testCase := range testCases {
		c := newFakeContext(testCase.mappings, testCase.args)
		r := newRunner()
		fakeEnv := rainforest.Environment{ID: 123, Name: "the foo environment"}
		client := &fakeRunnerClient{environment: fakeEnv}
		r.client = client

		err := r.startRun(c)
		if !testCase.wantError && err != nil {
			t.Error("Error starting run:", err)
		} else if testCase.wantError && err == nil {
			t.Errorf("Expected test to fail with args %v but it didn't", testCase.args)
		}

		var got []string
		for _, t := range client.createdTests {
			got = append(got, t.RFMLID)
		}
		sort.Strings(got)
		want := testCase.wantUpload
		if !reflect.DeepEqual(want, got) {
			t.Errorf("Tests were not uploaded correctly, wanted %v, got %v", want, got)
		}

		// Check that the right tests were requested for the run
		got = client.runParams.RFMLIDs
		sort.Strings(got)
		want = testCase.wantExecute
		if !reflect.DeepEqual(want, got) {
			t.Errorf("Incorrect tests were requested when starting run, wanted %v, got %v", want, got)
		}
	}
}

func TestBuildRerunArgs(t *testing.T) {
	testCases := []struct {
		Mappings map[string]interface{}
		Args     cli.Args
		RunID    int
		WantArgs []string
	}{
		{
			Mappings: map[string]interface{}{
				"junit-file":  "result.xml",
				"max-reruns":  uint(2),
				"skip-update": true,
			},
			Args:  cli.Args{},
			RunID: 123,
			WantArgs: []string{
				"rainforest-cli",
				"rerun",
				"123",
				"--max-reruns", "2",
				"--rerun-attempt", "1",
				"--skip-update",
				"--junit-file", "result.xml",
			},
		},
		{
			Mappings: map[string]interface{}{
				"max-reruns": uint(1),
				"token":      "deadbeef",
			},
			Args:  cli.Args{},
			RunID: 123,
			WantArgs: []string{
				"rainforest-cli",
				"rerun",
				"123",
				"--max-reruns", "1",
				"--rerun-attempt", "1",
				"--skip-update",
				"--token", "deadbeef",
			},
		},
	}
	for _, testCase := range testCases {
		c := newFakeContext(testCase.Mappings, testCase.Args)
		gotArgs, _ := buildRerunArgs(c, testCase.RunID)
		if !reflect.DeepEqual(gotArgs, testCase.WantArgs) {
			t.Errorf("\nWanted %v\n   got %v", testCase.WantArgs, gotArgs)
		}
	}
}
