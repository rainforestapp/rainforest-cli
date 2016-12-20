package rainforest

import (
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
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
