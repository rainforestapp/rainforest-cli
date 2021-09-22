// Package rainforest is a golang client for the Rainforest QA API
package rainforest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/rainforestapp/rainforest-cli/gittrigger"
	"github.com/ukd1/go.detectci"
	"github.com/whilp/git-urls"
)

const (
	// Version of the lib in SemVer
	libVersion = "2.0.1"

	currentBaseURL  = "https://app.rainforestqa.com/api/1/"
	authTokenHeader = "CLIENT_TOKEN"
)

// Client is responsible for communicating with Rainforest API
type Client struct {
	// http client used to connect to the Rainforest API
	client *http.Client

	// URL of a Rainforest API endpoint to be used by the client
	BaseURL *url.URL

	// String that will be set as an user agent with current library version appended to it
	UserAgent string

	// Save HTTP Response Headers
	LastResponseHeaders http.Header

	//Set debug flag to decide whether to return headers or not
	DebugFlag bool

	// Client token used for authenticating requests made to the RF
	clientToken string

	// Send telemetry with each API request using the user agent
	// this is used by Rainforest (and not shared or sold) to make
	// integrations better. See README for more details.
	SendTelemetry bool
}

// NewClient constructs a new rainforest API Client. As a parameter takes client token
// which is used for authentication and is available in the rainforest web app.
func NewClient(token string, debug bool) *Client {
	var baseURL *url.URL
	var err error
	if envURL := os.Getenv("RAINFOREST_API_URL"); envURL != "" {
		baseURL, err = url.Parse(envURL)
		if err != nil {
			log.Fatalf("Invalid URL set in $RAINFOREST_API_URL=%v", envURL)
		}
	} else {
		baseURL, _ = url.Parse(currentBaseURL)
	}

	return &Client{
		client:              http.DefaultClient,
		BaseURL:             baseURL,
		clientToken:         token,
		LastResponseHeaders: http.Header{},
		DebugFlag:           debug,
	}
}

// ClientToken returns the API authentication token for the client.
func (c *Client) ClientToken() string {
	return c.clientToken
}

// NewRequest creates an API request. Provided url will be resolved using ResolveReference,
// which works in a similar way to the hrefs in a browser (most important takeaway is to
// not add preceeding slash to the link as it resolves to a root path of domain).
// The body argument is JSON endoded and attached as a request body.
// This function also attaches auth token from the client to the request.
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	// Resolve the relative URL path
	relPath, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	endpointURL := c.BaseURL.ResolveReference(relPath)

	// Create buffer and fill it with body data encoded in JSON
	var b io.ReadWriter
	if body != nil {
		b = new(bytes.Buffer)
		err = json.NewEncoder(b).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	// Create new http request and set the headers
	req, err := http.NewRequest(method, endpointURL.String(), b)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Set the auth token to the one specified in the client
	if c.clientToken != "" {
		req.Header.Set(authTokenHeader, c.clientToken)
	} else {
		return nil, errors.New("Please provide your API Token with the RAINFOREST_API_TOKEN environment variable or --token global flag")
	}

	// Set UserAgent header with appended library version, will look like:
	// "rainforest-cli/2.1.0 [rainforest golang lib/2.0.0]"
	userAgent := []string{"rainforest", "golang", "lib/" + libVersion}

	if c.SendTelemetry {
		found, ci_name := detectci.WhichCI()
		if found {
			userAgent = append(userAgent, "ci/"+ci_name)
		}

		var remote string
		git, err := gitTrigger.NewGitTrigger()
		if err == nil {
			remote, err = git.GetRemote()
			if err == nil {
				u, err := giturls.Parse(remote)
				if err == nil {
					// Strip the user details, if any
					u.User = nil
					userAgent = append(userAgent, "repo/"+u.String())
				}
			}
		}
	}

	composedUserAgent := c.UserAgent + " [" + strings.Join(userAgent[:], " ") + "]"
	req.Header.Set("User-Agent", composedUserAgent)

	return req, nil
}

// checkResponse checks if we received vaild response with code 200,
// returns error otherwise
func checkResponse(res *http.Response, debugFlag bool) error {
	// If we are on a happy path just return nil
	if res.StatusCode >= 200 && res.StatusCode < 300 {
		return nil
	}

	errPrefix := fmt.Sprintf("RF API Error (%v)", res.StatusCode)

	// Otherwise we return error from the API or general one if we can't decode it
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return errors.New(errPrefix + " - Unable to read response: " + err.Error())
	}

	if contentType := res.Header.Get("Content-Type"); contentType != "application/json" {
		// Just print out the response body as a string
		return errors.New(errPrefix + ":\n" + string(body))
	}

	var out struct {
		Err string `json:"error"`
	}
	err = json.Unmarshal(body, &out)
	if err != nil {
		if debugFlag {
			fmt.Println("Cannot parse response:\n" + string(body))
		}

		return errors.New(errPrefix + " - Unable to parse response JSON: " + err.Error())
	}

	return errors.New(errPrefix + ": " + out.Err)
}

// Do sends out the request to the API and unpacks JSON response to the out variable.
func (c *Client) Do(req *http.Request, out interface{}) (*http.Response, error) {
	// Send out http request
	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if c.DebugFlag {
		log.Print("Trying ", res.Request.URL, "...")
	}

	// We check response for potential errors and return them to the caller.
	// We do not nil the response, as a caller might want to inspect the response in case of an error.
	err = checkResponse(res, c.DebugFlag)
	if err != nil {
		return res, err
	}

	// Here we check for the out pointer, and if it exists we unmarshall JSON there and return any
	// potential errors to the caller.
	if out != nil {
		// Close the body after we're done with it, to allow connection reuse.
		defer func() {
			io.Copy(ioutil.Discard, res.Body)
			res.Body.Close()
		}()

		err = json.NewDecoder(res.Body).Decode(out)

		if err != nil {
			log.Println("ERROR for ", req.Method, req.URL)
			log.Printf("ERROR PARSING JSON : %v\n\n", err.Error())
			return res, err
		}
	}

	c.LastResponseHeaders = res.Header

	if c.DebugFlag {
		log.Println("connected")
		printRequestHeaders(res)
	}

	return res, err
}

func printRequestHeaders(res *http.Response) {
	log.Println(res.Request.Method, res.Request.Proto)
	log.Println("User Agent:", res.Request.UserAgent())
	log.Println("Host:", res.Request.Host)
	log.Println("")
}
