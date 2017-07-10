package main

import (
	"reflect"
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
}

func (r *fakeRunnerClient) CreateTemporaryEnvironment(s string) (*rainforest.Environment, error) {
	return &r.environment, nil
}

func (r *fakeRunnerClient) CreateRun(p rainforest.RunParams) (*rainforest.RunStatus, error) {
	return nil, nil
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
		res, err := r.makeRunParams(c)

		if err != nil {
			t.Errorf("Error trying to create params: %v", err)
		} else if !reflect.DeepEqual(res, testCase.expected) {
			t.Errorf("Incorrect resulting run params.\nActual: %#v\nExpected: %#v", res, testCase.expected)
		}
	}
}
