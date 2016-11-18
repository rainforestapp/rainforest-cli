package main

import (
	"bytes"
	"os"
	"rainforest-cli/rainforest"
	"reflect"
	"regexp"
	"testing"
)

func TestFormatAsTable(t *testing.T) {
	var testCases = []struct {
		testSlice []rainforest.Resource
		want      [][]string
	}{
		{
			testSlice: []rainforest.Resource{
				rainforest.Site{
					ID:   1337,
					Name: "Dyer",
				},
				rainforest.Site{
					ID:   42,
					Name: "Situation",
				},
			},
			want: [][]string{{"1337", "Dyer"}, {"42", "Situation"}},
		},
		{
			testSlice: []rainforest.Resource{
				rainforest.Browser{
					Name:        "firefox",
					Description: "Mozilla Firefox",
				},
				rainforest.Browser{
					Name:        "ie11",
					Description: "Microsoft Internet Explorer 11",
				},
			},
			want: [][]string{{"firefox", "Mozilla Firefox"}, {"ie11", "Microsoft Internet Explorer 11"}},
		},
		{
			testSlice: []rainforest.Resource{
				rainforest.Folder{
					ID:    707,
					Title: "The Foo Folder",
				},
				rainforest.Folder{
					ID:    708,
					Title: "The Baz Folder",
				},
			},
			want: [][]string{{"707", "The Foo Folder"}, {"708", "The Baz Folder"}},
		},
	}

	for _, tCase := range testCases {
		got := formatAsTable(tCase.testSlice)
		if !reflect.DeepEqual(tCase.want, got) {
			t.Errorf("formatAsTable returned %+v, want %+v", got, tCase.want)
		}
	}
}

func TestPrintResourceTable(t *testing.T) {
	out = &bytes.Buffer{}
	defer func() {
		out = os.Stdout
	}()

	testBody := [][]string{{"1337", "Dyer"}, {"42", "Situation"}}
	printResourceTable("TEST", testBody)
	regexMatchOut(`\| +TEST ID +\| +TEST DESCRIPTION +\|`, t)
	regexMatchOut(`\| +1337 +\| +Dyer +\|`, t)
}

func regexMatchOut(pattern string, t *testing.T) {
	matched, err := regexp.Match(pattern, out.(*bytes.Buffer).Bytes())
	if err != nil {
		t.Error("Error with pattern match:", err)
	}
	if !matched {
		t.Errorf("Printed out %v, want %v", out, pattern)
	}
}
