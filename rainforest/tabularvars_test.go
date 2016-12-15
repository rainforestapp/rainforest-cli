package rainforest

import (
	"bytes"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestGetGenerators(t *testing.T) {
	setup()
	defer cleanup()

	const reqMethod = "GET"

	mux.HandleFunc("/generators", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != reqMethod {
			t.Errorf("Request method = %v, want %v", r.Method, reqMethod)
		}
		fmt.Fprint(w, `[{"id":1337,"created_at":"2016-11-24T14:56:19Z","name":"foo","description":`+
			`"bar","generator_type":"tabular","single_use":false,"columns":[{"id":30225,"created_at":"2016-11-24T14:56:19Z","name":"username"}`+
			`,{"id":30226,"created_at":"2016-11-24T14:56:19Z","name":"password"}],"row_count":42}]`)
	})

	out, _ := client.GetGenerators()

	want := []Generator{{
		ID:           1337,
		Name:         "foo",
		CreationDate: time.Date(2016, time.November, 24, 14, 56, 19, 0, time.UTC),
		Description:  "bar",
		Type:         "tabular",
		SingleUse:    false,
		RowCount:     42,
		Columns: []GeneratorColumn{
			{
				ID:           30225,
				Name:         "username",
				CreationDate: time.Date(2016, time.November, 24, 14, 56, 19, 0, time.UTC),
			},
			{
				ID:           30226,
				Name:         "password",
				CreationDate: time.Date(2016, time.November, 24, 14, 56, 19, 0, time.UTC),
			},
		},
	}}

	if !reflect.DeepEqual(out, want) {
		t.Errorf("Response out = %v, want %v", out, want)
	}
}

func TestDeleteGenerator(t *testing.T) {
	setup()
	defer cleanup()

	const reqMethod = "DELETE"
	const genID = 123

	mux.HandleFunc("/generators/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != reqMethod {
			t.Errorf("Request method = %v, want %v", r.Method, reqMethod)
		}
		if !strings.Contains(r.URL.String(), strconv.Itoa(genID)) {
			t.Errorf("Expected genID %v in URL, got %v", genID, r.URL.String())
		}

		// we don't care about reply
		fmt.Fprint(w, `{"foo": "bar"}`)
	})

	err := client.DeleteGenerator(genID)
	if err != nil {
		t.Errorf("Got error: %v", err.Error())
	}
}

func TestCreateTabularVar(t *testing.T) {
	setup()
	defer cleanup()

	const reqMethod = "POST"

	mux.HandleFunc("/generators", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != reqMethod {
			t.Errorf("Request method = %v, want %v", r.Method, reqMethod)
		}

		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		s := strings.TrimSpace(buf.String())
		if wantedBody := `{"name":"foo","description":"bar","columns":["baz","wut"]}`; s != wantedBody {
			t.Errorf("Request body = %v, want %v", s, wantedBody)
		}
		fmt.Fprint(w, `{"id":1337,"created_at":"2016-11-24T14:56:19Z","name":"foo","description":`+
			`"bar","generator_type":"tabular","single_use":false,"columns":[{"id":30225,"created_at":"2016-11-24T14:56:19Z","name":"baz"}`+
			`,{"id":30226,"created_at":"2016-11-24T14:56:19Z","name":"wut"}],"row_count":0}`)
	})

	newGenerator, _ := client.CreateTabularVar("foo", "bar", []string{"baz", "wut"}, false)
	want := &Generator{
		ID:           1337,
		Name:         "foo",
		CreationDate: time.Date(2016, time.November, 24, 14, 56, 19, 0, time.UTC),
		Description:  "bar",
		Type:         "tabular",
		SingleUse:    false,
		RowCount:     0,
		Columns: []GeneratorColumn{
			{
				ID:           30225,
				Name:         "baz",
				CreationDate: time.Date(2016, time.November, 24, 14, 56, 19, 0, time.UTC),
			},
			{
				ID:           30226,
				Name:         "wut",
				CreationDate: time.Date(2016, time.November, 24, 14, 56, 19, 0, time.UTC),
			},
		},
	}

	if !reflect.DeepEqual(newGenerator, want) {
		t.Errorf("Response out = %v, want %v", newGenerator, want)
	}
}

func TestAddGeneratorRows(t *testing.T) {
	setup()
	defer cleanup()

	const reqMethod = "POST"
	const genID = 999

	mux.HandleFunc("/generators/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != reqMethod {
			t.Errorf("Request method = %v, want %v", r.Method, reqMethod)
		}
		if !strings.Contains(r.URL.String(), strconv.Itoa(genID)) {
			t.Errorf("Expected genID %v in URL, got %v", genID, r.URL.String())
		}

		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		s := strings.TrimSpace(buf.String())
		if wantedBody := `{"data":[{"123":"foo","456":"baz"},{"123":"bar","456":"wut"}]}`; s != wantedBody {
			t.Errorf("Request body = %v, want %v", s, wantedBody)
		}

		// response is ignored
		fmt.Fprint(w, `{"foo":"bar"}`)
	})
	fakeGen := Generator{
		ID: genID,
		Columns: []GeneratorColumn{
			{
				ID:           123,
				Name:         "qwe",
				CreationDate: time.Date(2016, time.November, 24, 14, 56, 19, 0, time.UTC),
			},
			{
				ID:           456,
				Name:         "asd",
				CreationDate: time.Date(2016, time.November, 24, 14, 56, 19, 0, time.UTC),
			},
		},
	}
	rowData := []map[int]string{{123: "foo", 456: "baz"}, {123: "bar", 456: "wut"}}
	err := client.AddGeneratorRows(&fakeGen, rowData)
	if err != nil {
		t.Errorf("Got error: %v", err.Error())
	}
}

func TestAddGeneratorRowsFromTable(t *testing.T) {
	setup()
	defer cleanup()

	const reqMethod = "POST"
	const genID = 999

	mux.HandleFunc("/generators/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != reqMethod {
			t.Errorf("Request method = %v, want %v", r.Method, reqMethod)
		}
		if !strings.Contains(r.URL.String(), strconv.Itoa(genID)) {
			t.Errorf("Expected genID %v in URL, got %v", genID, r.URL.String())
		}

		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		s := strings.TrimSpace(buf.String())
		if wantedBody := `{"data":[{"123":"foo","456":"baz"},{"123":"bar","456":"wut"}]}`; s != wantedBody {
			t.Errorf("Request body = %v, want %v", s, wantedBody)
		}

		// response is ignored
		fmt.Fprint(w, `{"foo":"bar"}`)
	})
	fakeGen := Generator{
		ID: genID,
		Columns: []GeneratorColumn{
			{
				ID:           123,
				Name:         "qwe",
				CreationDate: time.Date(2016, time.November, 24, 14, 56, 19, 0, time.UTC),
			},
			{
				ID:           456,
				Name:         "asd",
				CreationDate: time.Date(2016, time.November, 24, 14, 56, 19, 0, time.UTC),
			},
		},
	}
	rowData := [][]string{{"foo", "baz"}, {"bar", "wut"}}
	err := client.AddGeneratorRowsFromTable(&fakeGen, []string{"qwe", "asd"}, rowData)
	if err != nil {
		t.Errorf("Got error: %v", err.Error())
	}
}
