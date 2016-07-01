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
		panic(fmt.Sprintf("fakeFlagParser does not expect argument: %s", s))
	}
}

func (f fakeFlagParser) Int(s string) int {
	switch s {
	case "smart-folder-id":
		return 123
	default:
		panic(fmt.Sprintf("fakeFlagParser does not expect argument: %s", s))
	}
}

func TestMakeParams(t *testing.T) {
	// flagReturner := fakeFlagParser{}
	// params := makeParams(&flagReturner)
	//
	// expectedTags := []string{"foo", "bar"}
	// actualTags := params.Tags
	//
	// if expectedLen, actualLen := len(expectedTags), len(actualTags); expectedLen != actualLen {
	// 	t.Errorf("Wrong amount of tags. Expected %d, got %d", expectedLen, actualLen)
	// }
	//
	// for i, actualTag := range actualTags {
	// 	if expectedTags[i] != actualTag {
	// 		t.Errorf("Unexpected tags parameter. Expected %s, got %s", expectedTags[1], actualTag)
	// 	}
	// }
	//
	// if params.SmartFolderID != 123 {
	// 	t.Errorf("Unexpected smart_folder_id parameter. Expected %d, got %d", 123, params.SmartFolderID)
	// }
}
