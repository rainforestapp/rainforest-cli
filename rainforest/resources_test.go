package rainforest

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestGetFolders(t *testing.T) {
	setup()
	defer cleanup()

	const reqMethod = "GET"

	mux.HandleFunc("/folders", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != reqMethod {
			t.Errorf("Request method = %v, want %v", r.Method, reqMethod)
		}
		fmt.Fprint(w, `[{"id": 707, "title": "Foo"}, {"id": 777, "title": "Bar"}]`)
	})

	out, _ := client.GetFolders()

	want := []Folder{{ID: 707, Title: "Foo"}, {ID: 777, Title: "Bar"}}

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
