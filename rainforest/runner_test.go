package rainforest

import (
	"bytes"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func TestCreateRun(t *testing.T) {
	var testCases = []struct {
		runParams RunParams
		wantBody  string
	}{
		{
			runParams: RunParams{Tags: []string{"foo", "bar"}, SiteID: 125, RunGroupID: 14},
			wantBody:  `{"tags":["foo","bar"],"site_id":125,"run_group_id":14}`,
		},
		{
			runParams: RunParams{RFMLIDs: []string{"login", "sign-up"}, SiteID: 125},
			wantBody:  `{"rfml_ids":["login","sign-up"],"site_id":125}`,
		},
	}
	for _, tc := range testCases {
		setup()

		const reqMethod = "POST"
		mux.HandleFunc("/runs", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != reqMethod {
				t.Errorf("Request method = %v, want %v", r.Method, reqMethod)
			}

			buf := new(bytes.Buffer)
			buf.ReadFrom(r.Body)
			s := strings.TrimSpace(buf.String())
			if s != tc.wantBody {
				t.Errorf("Request body = %v, want %v", s, tc.wantBody)
			}
			fmt.Fprint(w, `{"id": 123, "state":"in_progress"}`)
		})

		gotStatus, _ := client.CreateRun(tc.runParams)

		wantStatus := &RunStatus{ID: 123, State: "in_progress"}
		if !reflect.DeepEqual(gotStatus, wantStatus) {
			t.Errorf("Response out = %v, want %v", gotStatus, wantStatus)
		}
		cleanup()
	}
}

func TestCheckRunStatus(t *testing.T) {
	setup()
	defer cleanup()

	const reqMethod = "GET"

	mux.HandleFunc("/runs/123", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != reqMethod {
			t.Errorf("Request method = %v, want %v", r.Method, reqMethod)
		}

		fmt.Fprint(w, `{"id": 123, "state":"in_progress"}`)
	})

	out, _ := client.CheckRunStatus(123)

	want := &RunStatus{ID: 123, State: "in_progress"}

	if !reflect.DeepEqual(out, want) {
		t.Errorf("Response out = %v, want %v", out, want)
	}
}
