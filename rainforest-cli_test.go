package main

import (
	"reflect"
	"testing"

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
			testArgs: []string{"./rainforest", "run", "--tags", "tag,bag", "--token", "foobar", "--site", "123", "--run-group-id", "255"},
			want:     []string{"./rainforest", "--token", "foobar", "run", "--tags", "tag,bag", "--site", "123", "--run-group-id", "255"},
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
