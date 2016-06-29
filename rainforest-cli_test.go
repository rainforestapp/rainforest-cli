package main

import (
	"os"
	"reflect"
	"testing"
)

func TestParseCommands(t *testing.T) {
	var testCases = []struct {
		input []string
		want  []string
	}{
		{
			input: []string{"foo", "bar", "-words=baz"},
			want:  []string{"bar"},
		},
		{
			input: []string{"foo", "-words=baz", "bar"},
			want:  []string{"bar"},
		},
		{
			input: []string{"foo", "-numbers=321", "bar", "-words=baz"},
			want:  []string{"bar"},
		},
		{
			input: []string{"foo", "-numbers=321", "-words=baz"},
			want:  []string{},
		},
		{
			input: []string{"foo"},
			want:  []string{},
		},
		{
			input: []string{"-words=baz"},
			want:  []string{},
		},
		{
			input: []string{"foo", "bar", "wow"},
			want:  []string{"wow", "bar"},
		},
	}
	tempOsArgs := os.Args
	for _, test := range testCases {
		os.Args = test.input
		got := parseCommands()
		if !reflect.DeepEqual(test.want, got) {
			t.Errorf("Expected %v, got %v", test.want, got)
		}
	}
	os.Args = tempOsArgs
}
