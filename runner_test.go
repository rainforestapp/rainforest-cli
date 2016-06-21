package main

import (
	"fmt"
	"testing"
)

type fakeFlagParser struct{}

func (f fakeFlagParser) StringSlice(s string) []string {
	switch s {
	case "tags":
		return []string{"foo", "bar"}
	default:
		panic(fmt.Sprintf("fakeFlagParser does expect argument: %s", s))
	}
}

func TestMakeParams(t *testing.T) {
	flagReturner := fakeFlagParser{}
	params := makeParams(&flagReturner)

	expectedTags := []string{"foo", "bar"}
	actualTags := params.Tags

	if expectedLen, actualLen := len(expectedTags), len(actualTags); expectedLen != actualLen {
		t.Errorf("Wrong amount of tags. Expected %d, got %d", expectedLen, actualLen)
	}

	for i, actualTag := range actualTags {
		if expectedTags[i] != actualTag {
			t.Errorf("Unexpected tag. Expected %s, got %s", expectedTags[1], actualTag)
		}
	}
}
