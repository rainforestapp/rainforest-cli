package rainforest

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

func TestGetFolders(t *testing.T) {
	setup()
	defer cleanup()

	const reqMethod = "GET"
	const pages = 3
	var lastPage int

	mux.HandleFunc("/folders", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != reqMethod {
			t.Errorf("Request method = %v, want %v", r.Method, reqMethod)
		}
		w.Header().Add("X-Total-Pages", strconv.Itoa(pages))
		fmt.Fprint(w, `[{"id": 707, "title": "Foo"}, {"id": 777, "title": "Bar"}]`)
	})

	mux.HandleFunc("/folders?page_size=100&page=", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != reqMethod {
			t.Errorf("Request method = %v, want %v", r.Method, reqMethod)
		}

		reg, _ := regexp.Compile(`page_size=100&page=(\d+)`)
		matches := reg.FindStringSubmatch(r.URL.String())
		page, err := strconv.Atoi(matches[1])

		if err != nil {
			t.Fatal(err.Error())
		}

		if page > pages {
			t.Fatalf("Unexpected page argument: %v", page)
		}

		expectedPage := lastPage + 1
		if page != expectedPage {
			t.Errorf("Unexpected page argument. Want %v, Got %v", expectedPage, page)
		}

		lastPage = page

		w.Header().Add("X-Total-Pages", strconv.Itoa(pages))
		fmt.Fprint(w, `[{"id": 707, "title": "Foo"}, {"id": 777, "title": "Bar"}]`)
	})

	out, err := client.GetFolders()
	if err != nil {
		t.Fatal(err.Error())
	}

	want := []Folder{
		{ID: 707, Title: "Foo"}, {ID: 777, Title: "Bar"},
		{ID: 707, Title: "Foo"}, {ID: 777, Title: "Bar"},
		{ID: 707, Title: "Foo"}, {ID: 777, Title: "Bar"},
	}

	if !reflect.DeepEqual(out, want) {
		t.Errorf("Response out = %v, want %v", out, want)
	}
}

func TestGetBrowsers(t *testing.T) {
	setup()
	defer cleanup()

	const reqMethod = "GET"

	mux.HandleFunc("/clients", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != reqMethod {
			t.Errorf("Request method = %v, want %v", r.Method, reqMethod)
		}
		fmt.Fprint(w, `{"available_browsers": [{"name": "firefox", "description": "Mozilla Firefox"}]}`)
	})

	out, _ := client.GetBrowsers()

	want := []Browser{{Name: "firefox", Description: "Mozilla Firefox"}}

	if !reflect.DeepEqual(out, want) {
		t.Errorf("Response out = %v, want %v", out, want)
	}
}

func TestGetSites(t *testing.T) {
	setup()
	defer cleanup()

	const reqMethod = "GET"

	mux.HandleFunc("/sites", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != reqMethod {
			t.Errorf("Request method = %v, want %v", r.Method, reqMethod)
		}
		fmt.Fprint(w, `[{"id": 1337, "name": "Dyer"}, {"id": 31337, "name": "Situation"}]`)
	})

	out, _ := client.GetSites()

	want := []Site{{ID: 1337, Name: "Dyer"}, {ID: 31337, Name: "Situation"}}

	if !reflect.DeepEqual(out, want) {
		t.Errorf("Response out = %v, want %v", out, want)
	}
}

func TestGetRunGroups(t *testing.T) {
	setup()
	defer cleanup()

	const reqMethod = "GET"

	mux.HandleFunc("/run_groups", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != reqMethod {
			t.Errorf("Request method = %v, want %v", r.Method, reqMethod)
		}
		fmt.Fprint(w, `[{"id": 3457, "title": "Star"}, {"id": 289, "title": "Trek"}]`)
	})

	out, _ := client.GetRunGroups()

	want := []RunGroup{{ID: 3457, Title: "Star"}, {ID: 289, Title: "Trek"}}

	if !reflect.DeepEqual(out, want) {
		t.Errorf("Response out = %v, want %v", out, want)
	}
}

func TestRunGroupDetailsPrint(t *testing.T) {
	rgd := RunGroupDetails{
		ID:    6678,
		Title: "Main run group",
		Environment: struct {
			Name string `json:"name"`
		}{
			Name: "staging",
		},
		Crowd:      "default",
		RerouteGeo: "usa",
	}

	output, _ := captureStdout(func() error {
		rgd.Print()
		return nil
	})
	expectedNameStr := "Name: Main run group"
	if !strings.Contains(output, expectedNameStr) {
		t.Errorf("Run group name was not printed properly.\nExpected: %v\nto be included in: %v", expectedNameStr, output)
	}

	expectedEnvStr := "Environment: staging"
	if !strings.Contains(output, expectedEnvStr) {
		t.Errorf("Run group environment was not printed properly.\nExpected: %v\nto be included in: %v", expectedEnvStr, output)
	}

	expectecCrowdStr := "Tester Crowd: default"
	if !strings.Contains(output, expectecCrowdStr) {
		t.Errorf("Run group tester crowd was not printed properly.\nExpected: %v\nto be included in: %v", expectecCrowdStr, output)
	}

	expectedLocationStr := "Location: usa"
	if !strings.Contains(output, expectedLocationStr) {
		t.Errorf("Run group location was not printed properly.\nExpected: %v\nto be included in: %v", expectedLocationStr, output)
	}

	if strings.Contains(output, "Schedule:") {
		t.Error("Run group schedule found in output when schedule does not exist")
	}

	rgd.Schedule = struct {
		RepeatRules []struct {
			Day  string `json:"day"`
			Time string `json:"time"`
		} `json:"repeat_rules"`
	}{
		RepeatRules: []struct {
			Day  string `json:"day"`
			Time string `json:"time"`
		}{
			{
				Day:  "tuesday",
				Time: "12:00",
			},
			{
				Day:  "friday",
				Time: "12:00",
			},
		},
	}

	output, _ = captureStdout(func() error {
		rgd.Print()
		return nil
	})
	expectedScheduleStr := "Schedule: tuesday, friday @ 12:00"
	if !strings.Contains(output, expectedScheduleStr) {
		t.Errorf("Run group schedule was not printed properly.\nExpected: %v\nto be included in: %v", expectedScheduleStr, output)
	}
}

func captureStdout(fn func() error) (string, error) {
	r, w, err := os.Pipe()
	if err != nil {
		return "", err
	}

	originalStdOut := os.Stdout
	os.Stdout = w
	defer func() {
		os.Stdout = originalStdOut
	}()

	err = fn()
	if err != nil {
		return "", err
	}

	w.Close()

	buf := bytes.Buffer{}
	_, err = io.Copy(&buf, r)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
