package main

import (
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/rainforestapp/rainforest-cli/rainforest"
	"github.com/urfave/cli"
)

func TestShuffleFlags(t *testing.T) {
	var testCases = []struct {
		testArgs []string
		want     []string
	}{
		{
			testArgs: []string{"./rainforest", "--token", "foobar", "run", "--tags", "tag,bag"},
			want:     []string{"./rainforest", "--token", "foobar", "run", "--tags", "tag,bag"},
		},
		{
			testArgs: []string{"./rainforest", "run", "--tags", "tag,bag", "--token", "foobar"},
			want:     []string{"./rainforest", "--token", "foobar", "run", "--tags", "tag,bag"},
		},
		{
			testArgs: []string{"./rainforest", "run", "--tags", "tag,bag", "--token", "foobar", "--site", "123"},
			want:     []string{"./rainforest", "--token", "foobar", "run", "--tags", "tag,bag", "--site", "123"},
		},
		{
			testArgs: []string{"./rainforest", "run", "--tags", "tag,bag", "--token", "foobar", "--site", "123", "--debug"},
			want:     []string{"./rainforest", "--token", "foobar", "--debug", "run", "--tags", "tag,bag", "--site", "123"},
		},
		{
			testArgs: []string{"./rainforest", "--skip-update", "run", "--tags", "tag,bag", "--token", "foobar", "--site", "123"},
			want:     []string{"./rainforest", "--skip-update", "--token", "foobar", "run", "--tags", "tag,bag", "--site", "123"},
		},
	}

	for _, tCase := range testCases {
		got := shuffleFlags(tCase.testArgs)
		if !reflect.DeepEqual(tCase.want, got) {
			t.Errorf("shuffleFlags returned %+v, want %+v", got, tCase.want)
		}
	}
}

func TestDebugFlag(t *testing.T) {

	testCases := []struct {
		mappings map[string]interface{}
		args     []string
		runID    int
		debug    bool
		tag      string
		token    string
		method   string
	}{
		{
			mappings: map[string]interface{}{
				"token":  "testToken123",
				"debug":  true,
				"run-id": 564,
				"tag":    "star",
			},
			args:   []string{"rainforest", "--token", "testToken123", "--debug", "run", "--tag", "star"},
			runID:  564,
			debug:  true,
			tag:    "star",
			token:  "testToken123",
			method: "GET",
		},
		{
			mappings: map[string]interface{}{
				"token":  "testToken123",
				"debug":  false,
				"run-id": 4335,
				"tag":    "star",
			},
			args:   []string{"rainforest", "--token", "testToken123", "run", "--tag", "star"},
			runID:  4335,
			debug:  false,
			tag:    "star",
			token:  "testToken123",
			method: "POST",
		},
	}

	for _, testCase := range testCases {
		c := newFakeContext(testCase.mappings, testCase.args)
		client := rainforest.NewClient(testCase.token, c.Bool("debug"))
		client.BaseURL, _ = url.Parse("https://example.org")

		req, _ := client.NewRequest(testCase.method, "/", nil)
		client.Do(req, nil)

		checkString := strings.Join(testCase.args, " ")
		if out := strings.Contains(checkString, "debug"); out != client.DebugFlag {
			t.Errorf("It is %+v that the --debug flag was in the command line arguments. However, the value was actually %+v.", out, client.DebugFlag)
		}
	}

}

// fakeContext is a helper for testing the cli interfacing functions
type fakeContext struct {
	mappings map[string]interface{}
	args     cli.Args
}

func (f fakeContext) String(s string) string {
	val, ok := f.mappings[s].(string)

	if ok {
		return val
	}
	return ""
}

func (f fakeContext) StringSlice(s string) []string {
	val, ok := f.mappings[s].([]string)

	if ok {
		return val
	}
	return []string{}
}

func (f fakeContext) Bool(s string) bool {
	val, ok := f.mappings[s].(bool)

	if ok {
		return val
	}
	return false
}

func (f fakeContext) Int(s string) int {
	val, ok := f.mappings[s].(int)

	if ok {
		return val
	}
	return 0
}

func (f fakeContext) Args() cli.Args {
	return f.args
}

func newFakeContext(mappings map[string]interface{}, args cli.Args) *fakeContext {
	return &fakeContext{mappings, args}
}
