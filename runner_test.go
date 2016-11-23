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

type fakeContext struct {
	mappings map[string]interface{}
}

func (f fakeContext) String(s string) string {
	val, ok := f.mappings[s].(string)

	if ok {
		return val
	}
	return ""
}

func (f fakeContext) StringSlice(s string) []string {
	return []string{}
}

func (f fakeContext) Args() cli.Args {
	return []string{}
}

func TestMakeRunParams(t *testing.T) {
	c := fakeContext{}

	res, err := makeRunParams(c)

	expected := rainforest.RunParams{}

	if err != nil {
		t.Errorf("Error trying to create params: %v", err)
	} else if !reflect.DeepEqual(res, expected) {
		t.Errorf("Incorrect value for conflict.\nActual: %#v\nExpected: %#v", res, expected)
	}
}
