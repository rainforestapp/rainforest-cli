package rainforest

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
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

	client = NewClient("testToken123")
	url, _ := url.Parse(fakeServer.URL)
	client.BaseURL = url
}

// close the fakeServer
func cleanup() {
	fakeServer.Close()
}

func TestNewClient(t *testing.T) {
	token := "testToken123"
	client = NewClient(token)
	if out := client.ClientToken; client.ClientToken != token {
		t.Errorf("NewClient didn't set proper token %+v, want %+v", out, token)
	}
}

func TestNewRequest(t *testing.T) {
	token := "testToken123"
	client = NewClient(token)
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
			httpResp:  &http.Response{StatusCode: 500},
			wantError: true,
		},
		{
			httpResp:  &http.Response{StatusCode: 103},
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
