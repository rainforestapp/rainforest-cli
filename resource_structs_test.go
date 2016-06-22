package main

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestResourcesTableSlice(t *testing.T) {
	var testCases = []struct {
		resource      returnTable
		json          []byte
		expectedTable [][]string
	}{
		{
			resource:      new(foldersResp),
			json:          fakeFolderResp,
			expectedTable: [][]string{{"707", "The Foo Folder"}, {"708", "The Baz Folder"}},
		},
		{
			resource:      new(sitesResp),
			json:          fakeSitesResp,
			expectedTable: [][]string{{"1337", "Dyer"}, {"42", "Situation"}},
		},
		{
			resource:      new(browsersResp),
			json:          fakeClientsResp,
			expectedTable: [][]string{{"firefox", "Mozilla Firefox"}, {"ie11", "Microsoft Internet Explorer 11"}},
		},
	}
	for _, tcase := range testCases {
		err := json.Unmarshal(tcase.json, tcase.resource)
		if err != nil {
			panic(err)
		}
		checkTablesEqual(t, tcase.resource, tcase.expectedTable)
	}

}

func checkTablesEqual(t *testing.T, resBody returnTable, want [][]string) {
	got := resBody.TableSlice()
	if expectedLen, actualLen := len(want), len(got); expectedLen != actualLen {
		t.Errorf("Wrong number of matrix rows. Expected %d, got %d", expectedLen, actualLen)
	}
	if !reflect.DeepEqual(want, got) {
		wantB, _ := json.Marshal(want)
		gotB, _ := json.Marshal(got)
		t.Log("want:")
		t.Logf("\t%+v", string(wantB))
		t.Log("got =")
		t.Errorf("\t%+v", string(gotB))
	}
}

var fakeFolderResp = []byte(`[
  {
    "id": 707,
    "created_at": "2016-04-18T18:09:42Z",
    "title": "The Foo Folder",
    "logic": [
      {
        "tag": "foo",
        "inclusive": true
      }
    ],
    "test_count": 0
  },
  {
    "id": 708,
    "created_at": "2016-04-18T18:09:51Z",
    "title": "The Baz Folder",
    "logic": [
      {
        "tag": "baz",
        "inclusive": true
      }
    ],
    "test_count": 0
  }
]`)

var fakeSitesResp = []byte(`[
    {
      "id": 1337,
      "created_at": "2016-02-23T06:12:38Z",
      "name": "Dyer",
      "default": true
    },
    {
      "id": 42,
      "created_at": "2016-02-23T06:12:38Z",
      "name": "Situation",
      "default": true
    }
  ]`)

var fakeClientsResp = []byte(`{
  "id": 4938,
  "name": "Edward CLI testing",
  "enabled_features": [
    "test_variables_v1"
  ],
  "default_environment_id": 5334,
  "billing_email": "edward@rainforestapp.com",
  "test_count": 44,
  "available_browsers": [
    {
      "name": "firefox",
      "description": "Mozilla Firefox",
      "category": "browser",
      "browser_version": "43.0.3",
      "os_version": "Windows 7 Ultimate N (SP1)",
      "default": true
    },
    {
      "name": "ie11",
      "description": "Microsoft Internet Explorer 11",
      "category": "browser",
      "browser_version": "11.0.9600.17843",
      "os_version": "Windows 7 Ultimate N (SP1)",
      "default": true
    }
  ],
  "owner_email": "edward@rainforestapp.com"
}`)
