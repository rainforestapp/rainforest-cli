package main

import (
	"flag"
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

func TestMakeRunParams(t *testing.T) {
	var testCases = []struct {
		flags          [][]string
		expectedResult rainforest.RunParams
	}{
		{
			flags:          [][]string{},
			expectedResult: rainforest.RunParams{},
		},
	}

	for _, tCase := range testCases {
		fakeApp := cli.NewApp()
		fakeFlagSet := flag.NewFlagSet("fakeFlagSet", 1)
		fakeContext := cli.NewContext(fakeApp, fakeFlagSet, nil)

		for _, flag := range tCase.flags {
			fakeFlagSet.Set(flag[0], flag[1])
		}

		result, err := makeRunParams(fakeContext)

		if err != nil {
			t.Errorf("Received error from makeRunParams: %v", err)
		}

		if !reflect.DeepEqual(result, tCase.expectedResult) {
			t.Errorf("Unexpected results from makeRunParams.\nExpected: %#v\nActual: %#v", tCase.expectedResult, result)
		}
	}

}
