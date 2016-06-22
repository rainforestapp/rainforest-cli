package main

import (
	"encoding/json"
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
			json:          fakeFolderByte,
			expectedTable: [][]string{{"707", "The Foo Folder"}, {"708", "The Baz Folder"}},
		},
		{
			resource:      new(sitesResp),
			json:          fakeSitesByte,
			expectedTable: [][]string{{"1337", "Dyer"}, {"42", "Situation"}},
		},
		{
			resource:      new(browsersResp),
			json:          fakeClientsByte,
			expectedTable: [][]string{{"firefox", "Mozilla Firefox"}, {"ie11", "Microsoft Internet Explorer 11"}},
		},
	}
	for _, tcase := range testCases {
		err := json.Unmarshal(tcase.json, tcase.resource)
		if err != nil {
			panic(err)
		}
		matrixTestHelper(t, tcase.resource, tcase.expectedTable)
	}

}

func matrixTestHelper(t *testing.T, resBody returnTable, expectedTable [][]string) {
	actualTable := resBody.TableSlice()
	if expectedLen, actualLen := len(expectedTable), len(actualTable); expectedLen != actualLen {
		t.Errorf("Wrong number of matrix rows. Expected %d, got %d", expectedLen, actualLen)
	}

	for i, actualrow := range actualTable {
		for j, actualColumn := range actualrow {
			if expectedTable[i][j] != actualColumn {
				t.Errorf("Unexpected matrix entry. Expected %s, got %s", expectedTable[i][j], actualColumn)
			}
		}
	}
}

var fakeFolderByte = []byte(`[
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

var fakeSitesByte = []byte(`[
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

var fakeClientsByte = []byte(`{
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
