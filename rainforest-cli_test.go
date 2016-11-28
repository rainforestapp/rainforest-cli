package main

import "github.com/urfave/cli"

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
