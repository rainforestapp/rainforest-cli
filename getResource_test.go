package main

import (
	"encoding/json"
	"testing"
)

type fakeReturnTable []struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

func (f fakeReturnTable) TableSlice() [][]string {
	returnSlice := []string{"200", "Test OK!"}
	returnArray := [][]string{returnSlice}
	return returnArray
}

func TestFoldersTableSlice(t *testing.T) {
	var resBody *foldersResp
	json.Unmarshal(fakeFolderByte, &resBody)
	slice1 := []string{"707", "The Foo Folder"}
	slice2 := []string{"708", "The Baz Folder"}
	expectedTable := [][]string{slice1, slice2}
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
