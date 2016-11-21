package main

import (
	"reflect"
	"testing"
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
