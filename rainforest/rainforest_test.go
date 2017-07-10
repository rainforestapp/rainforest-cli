package rainforest

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

var (
	fakeServer *httptest.Server
	// mux on which tests need to registers endpoint handlers to reply
	// with mock data.
	mux *http.ServeMux
	// API Client configured to use fake server
	client *Client
)

// configure testing server and fill globals with it
func setup() {
	mux = http.NewServeMux()
	fakeServer = httptest.NewServer(mux)

	client = NewClient("testToken123", false)
	url, _ := url.Parse(fakeServer.URL)
	client.BaseURL = url
}

// close the fakeServer
func cleanup() {
	fakeServer.Close()
}

func TestNewClient(t *testing.T) {
	token := "testToken123"
	client = NewClient(token, false)
	if out := client.ClientToken; client.ClientToken != token {
		t.Errorf("NewClient didn't set proper token %+v, want %+v", out, token)
	}
}

func TestNewRequest(t *testing.T) {
	token := "testToken123"
	client = NewClient(token, false)
	client.BaseURL, _ = url.Parse("https://example.org")
	req, _ := client.NewRequest("GET", "test", nil)
	if out := req.Header.Get(authTokenHeader); out != token {
		t.Errorf("NewRequest didn't set proper token header %+v, want %+v", out, token)
	}
	if out := req.URL; out.String() != "https://example.org/test" {
		t.Errorf("NewRequest didn't set proper URL %+v, want %+v", out, "https://example.org/test")
	}
	if req.Body != nil {
		t.Fatalf("constructed request contains a non-nil Body")
	}

	// Should not make any HTTP requests without a token
	var err error
	client = NewClient("", false)
	req, err = client.NewRequest("GET", "/", nil)

	if err == nil {
		t.Error("Expected an error")
	} else if !strings.Contains(err.Error(), "Please provide your API Token") {
		t.Errorf("Expected error for missing API token, got \"%v\"", err.Error())
	}
}

func TestCheckResponse(t *testing.T) {
	var testCases = []struct {
		httpResp  *http.Response
		wantError bool
	}{
		{
			httpResp:  &http.Response{StatusCode: 200},
			wantError: false,
		},
		{
			httpResp:  &http.Response{StatusCode: 201},
			wantError: false,
		},
		{
			httpResp: &http.Response{
				StatusCode: 500,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"error": "foo"}`)),
			},
			wantError: true,
		},
		{
			httpResp: &http.Response{
				StatusCode: 103,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`Totally not JSON`)),
			},
			wantError: true,
		},
	}

	for _, tCase := range testCases {
		got := checkResponse(tCase.httpResp)
		if tCase.wantError && got == nil {
			t.Errorf("checkResponse should've returned error, got %+v", got)
		}
		if !tCase.wantError && got != nil {
			t.Errorf("checkResponse should've returned nil, got %+v", got)
		}
	}
}

func TestDo(t *testing.T) {
	setup()
	defer cleanup()

	type testJSON struct {
		TestString string `json:"test_string"`
	}

	const reqMethod = "GET"

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != reqMethod {
			t.Errorf("Request method = %v, want %v", r.Method, reqMethod)
		}
		fmt.Fprint(w, `{"test_string":"foobar"}`)
	})

	req, _ := client.NewRequest("GET", "/", nil)
	var out testJSON
	client.Do(req, &out)

	want := testJSON{"foobar"}
	if !reflect.DeepEqual(out, want) {
		t.Errorf("Response out = %v, want %v", out, want)
	}
}

func TestNewClientWithDebug(t *testing.T) {
	testCases := []struct {
		args   []string
		runID  int
		debug  bool
		tag    string
		token  string
		method string
	}{
		{
			args:   []string{"rainforest", "--token", "testToken123", "--debug", "run", "--tag", "star"},
			runID:  564,
			debug:  true,
			tag:    "star",
			token:  "testToken123",
			method: "GET",
		},
		{
			args:   []string{"rainforest", "--token", "testToken123", "run", "--tag", "star"},
			runID:  4335,
			debug:  false,
			tag:    "star",
			token:  "testToken123",
			method: "POST",
		},
	}

	for _, testCase := range testCases {
		client := NewClient(testCase.token, testCase.debug)
		client.BaseURL, _ = url.Parse("https://example.org")
		req, _ := client.NewRequest(testCase.method, "/", nil)
		if out := req.URL; out.String() != "https://example.org/" {
			t.Errorf("NewRequest didn't set proper URL %+v, want %+v", out, "https://example.org/")
		}
		client.Do(req, nil)

		checkString := strings.Join(testCase.args, " ")
		if out := strings.Contains(checkString, "debug"); out != client.DebugFlag {
			t.Errorf("It is %+v that the --debug flag was in the command line arguments. However, the value was actually %+v.", out, client.DebugFlag)
		}
	}
}
