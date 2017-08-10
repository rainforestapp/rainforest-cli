package main

import (
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

func (r *fakeRunnerClient) GetRFMLIDs() (rainforest.TestIDMappings, error) {
	var mappings rainforest.TestIDMappings
	for _, test := range r.createdTests {
		mappings = append(mappings, rainforest.TestIDMap{ID: test.TestID, RFMLID: test.RFMLID})
	}
	return mappings, nil
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
				"environment-id": "1337",
				"tag":            []string{"tag", "tag2,tag3"},
				"run-group-id":   "14",
			},
			args: cli.Args{"12", "34", "56, 78"},
			expected: rainforest.RunParams{
				SmartFolderID: 123,
				SiteID:        456,
				Crowd:         "on_premise_crowd",
				Conflict:      "abort",
				Browsers:      []string{"chrome", "firefox", "safari"},
				Description:   "my awesome description",
				EnvironmentID: 1337,
				Tags:          []string{"tag", "tag2", "tag3"},
				Tests:         []int{12, 34, 56, 78},
				RunGroupID:    14,
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

func TestStartLocalRun(t *testing.T) {
	rfmlDir := setupTestRFMLDir()
	defer os.RemoveAll(rfmlDir)

	testCases := []struct {
		mappings    map[string]interface{}
		args        cli.Args
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
				"f":            true,
				"bg":           true,
				"execute":      []string{filepath.Join(rfmlDir, "b/b/b3.rfml")},
				"dont-execute": []string{filepath.Join(rfmlDir, "a/a2.rfml")},
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
		if err != nil {
			t.Error("Error starting run:", err)
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

// func TestStartLocalRunExecuteOverride(t *testing.T) {
// 	rfmlDir := setupTestRFMLDir()
// 	defer os.RemoveAll(rfmlDir)

// 	c := newFakeContext(mappings, args)
// 	r := newRunner()
// 	fakeEnv := rainforest.Environment{ID: 123, Name: "the foo environment"}
// 	client := &fakeRunnerClient{environment: fakeEnv}
// 	r.client = client

// 	err := r.startRun(c)
// 	if err != nil {
// 		t.Error("Error starting run:", err)
// 	}
// }
