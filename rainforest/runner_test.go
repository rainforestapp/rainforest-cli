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
	testCases := []struct {
		runParams RunParams
		wantBody  string
	}{
		{
			runParams: RunParams{Tags: []string{"foo", "bar"}, SiteID: 125},
			wantBody:  `{"tags":["foo","bar"],"site_id":125}`,
		},
		{
			runParams: RunParams{FeatureID: 25, Tags: []string{"baz"}},
			wantBody:  `{"tags":["baz"],"feature_id":25}`,
		},
	}

	for _, tc := range testCases {
		setup()
		defer cleanup()

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

		wantStatus := &RunStatus{ID: 123, State: "in_progress"}
		gotStatus, _ := client.CreateRun(tc.runParams)
		if !reflect.DeepEqual(gotStatus, wantStatus) {
			t.Errorf("Response out = %v, want %v", gotStatus, wantStatus)
		}
	}
}

func TestCreateRunFromRerun(t *testing.T) {
	setup()
	defer cleanup()

	const reqMethod = "POST"
	runID := 42
	var completed = false

	runParams := RunParams{
		RunID:    runID,
		Conflict: "abort",
	}
	wantBody := `{"conflict":"abort"}`

	mux.HandleFunc("/runs", func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("/runs endpoint hit, expected /runs/:id/rerun_failed")
	})

	wantPath := fmt.Sprintf("/runs/%v/rerun_failed", runID)
	mux.HandleFunc(wantPath, func(w http.ResponseWriter, r *http.Request) {
		completed = true
		if r.Method != reqMethod {
			t.Errorf("Request method = %v, want %v", r.Method, reqMethod)
		}

		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		s := strings.TrimSpace(buf.String())
		if s != wantBody {
			t.Errorf("Request body = %v, want %v", s, wantBody)
		}
		fmt.Fprint(w, `{"id": 123, "state":"in_progress"}`)
	})
	out, err := client.CreateRun(runParams)
	if err != nil {
		t.Error("Error creating run:", err)
	}

	wantStatus := &RunStatus{ID: 123, State: "in_progress"}
	if !reflect.DeepEqual(out, wantStatus) {
		t.Errorf("Response out = %v, want %v", out, wantStatus)
	}

	if !completed {
		t.Error("Run API endpoint not hit")
	}

	runID = 117
	// Check that you can't combine filters
	badParams := []RunParams{
		{
			RunID: runID,
			Tags:  []string{"foo"},
		},
		{
			RunID:    runID,
			Browsers: []string{"chrome"},
		},
		{
			RunID: runID,
			Tests: "all",
		},
		{
			RunID:   runID,
			RFMLIDs: []string{"rmfl1"},
		},
		{
			RunID: runID,
			Crowd: "automation",
		},
		{
			RunID:       runID,
			Description: "My fancy rerun",
		},
		{
			RunID:   runID,
			Release: "somerandomsha",
		},
		{
			RunID:         runID,
			EnvironmentID: 35,
		},
		{
			RunID:  runID,
			SiteID: 35,
		},
		{
			RunID:     runID,
			FeatureID: 35,
		},
		{
			RunID:         runID,
			SmartFolderID: 123,
		},
		{
			RunID:      runID,
			RunGroupID: 123,
		},
	}
	mux.HandleFunc(fmt.Sprintf("/runs/%v/rerun_failed", runID), func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 123, "state":"in_progress"}`)
	})
	for _, params := range badParams {
		_, err := client.CreateRun(params)
		if err == nil {
			t.Errorf("Expected error for params %v but there was no error", params)
		}
	}
}

func TestCreateRunFromRunGroup(t *testing.T) {
	setup()
	defer cleanup()

	const reqMethod = "POST"
	runGroupID := 25
	var completed = false

	runParams := RunParams{
		EnvironmentID: 23,
		Conflict:      "abort",
		RunGroupID:    runGroupID,
	}
	wantBody := `{"conflict":"abort","environment_id":23}`

	mux.HandleFunc("/runs", func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("/runs endpoint hit, expected /run_groups/:id/runs")
	})

	wantPath := fmt.Sprintf("/run_groups/%v/runs", runGroupID)
	mux.HandleFunc(wantPath, func(w http.ResponseWriter, r *http.Request) {
		completed = true
		if r.Method != reqMethod {
			t.Errorf("Request method = %v, want %v", r.Method, reqMethod)
		}

		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		s := strings.TrimSpace(buf.String())
		if s != wantBody {
			t.Errorf("Request body = %v, want %v", s, wantBody)
		}
		fmt.Fprint(w, `{"id": 123, "state":"in_progress"}`)
	})

	out, err := client.CreateRun(runParams)
	if err != nil {
		t.Error("Error creating run:", err)
	}

	wantStatus := &RunStatus{ID: 123, State: "in_progress"}
	if !reflect.DeepEqual(out, wantStatus) {
		t.Errorf("Response out = %v, want %v", out, wantStatus)
	}

	if !completed {
		t.Error("Run API endpoint not hit")
	}

	runGroupID = 42
	// Check that you can't combine filters
	badParams := []RunParams{
		{
			RunGroupID: runGroupID,
			Tags:       []string{"foo"},
		},
		{
			RunGroupID: runGroupID,
			SiteID:     17,
		},
		{
			RunGroupID: runGroupID,
			Tests:      "all",
		},
		{
			RunGroupID:    runGroupID,
			SmartFolderID: 123,
		},
		{
			RunGroupID: runGroupID,
			Browsers:   []string{"chrome"},
		},
		{
			RunGroupID: runGroupID,
			FeatureID:  35,
		},
	}
	mux.HandleFunc(fmt.Sprintf("/run_groups/%v/runs", runGroupID), func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 123, "state":"in_progress"}`)
	})
	for _, params := range badParams {
		_, err := client.CreateRun(params)
		if err == nil {
			t.Errorf("Expected error for params %v but there was no error", params)
		}
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

func TestLastMatchingRun(t *testing.T) {
	const reqMethod = "GET"

	testCases := []struct {
		runParams          RunParams
		wantBodyByRunGroup string
		wantBodyByRun      string
		wantStatus         RunStatus
	}{
		// When the last identical run passed
		{
			runParams:          RunParams{RunGroupID: 123},
			wantBodyByRunGroup: `[{"id": 111, "result":"passed", "run_group_id": 123}]`,
			wantBodyByRun:      `[]`,
			wantStatus:         RunStatus{ID: 111, Result: "passed"},
		},
		// When the last identical run failed and its rerun passed
		{
			runParams:          RunParams{RunGroupID: 123},
			wantBodyByRunGroup: `[{"id": 111, "result":"failed", "run_group_id": 123}]`,
			wantBodyByRun:      `[{"id": 222, "run_id": 111, "result":"passed"}]`,
			wantStatus:         RunStatus{ID: 222, Result: "passed"},
		},
		// When the last identical run failed and its rerun also failed
		{
			runParams:          RunParams{RunGroupID: 123},
			wantBodyByRunGroup: `[{"id": 111, "result":"failed", "run_group_id": 123}]`,
			wantBodyByRun:      `[{"id": 222, "run_id": 111, "result":"failed"}]`,
			wantStatus:         RunStatus{ID: 222, Result: "failed"},
		},
	}

	for _, tc := range testCases {
		setup()
		defer cleanup()
		mux.HandleFunc("/runs", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != reqMethod {
				t.Errorf("Request method = %v, want %v", r.Method, reqMethod)
			}
			if len(r.URL.Query()["run_id"]) > 0 {
				fmt.Fprint(w, tc.wantBodyByRun)
			} else {
				fmt.Fprint(w, tc.wantBodyByRunGroup)
			}
		})

		out, err := client.LastMatchingRun(tc.runParams)
		if err != nil {
			t.Fatal(err.Error())
		}

		if !(out.ID == tc.wantStatus.ID && out.Result == tc.wantStatus.Result) {
			t.Errorf("Response out = %v, want %v", out, tc.wantStatus)
		}
	}
}
