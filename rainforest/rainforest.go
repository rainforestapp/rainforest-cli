// Package rainforest is a golang client for the Rainforest QA API
package rainforest

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

const (
	// Version of the lib in SemVer
	libVersion = "2.0.0"

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

	// Client token used for authenticating requests made to the RF
	ClientToken string
}

// NewClient constructs a new rainforest API Client. As a parameter takes client token
// which is used for authentication and is available in the rainforest web app.
func NewClient(token string) *Client {
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

	return &Client{client: http.DefaultClient, BaseURL: baseURL, ClientToken: token}
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
	if c.ClientToken != "" {
		req.Header.Set(authTokenHeader, c.ClientToken)
	}

	// Set UserAgent header with appended library version, will look like:
	// "rainforest-cli/2.1.0 [rainforest golang lib/2.0.0]"
	composedUserAgent := c.UserAgent + " [rainforest golang lib/" + libVersion + "]"
	req.Header.Set("User-Agent", composedUserAgent)

	return req, nil
}

// checkResponse checks if we received vaild response with code 200,
// returns error otherwise
func checkResponse(res *http.Response) error {
	// If we are on a happy path just return nil
	if res.StatusCode >= 200 && res.StatusCode < 300 {
		return nil
	}

	// Otherwise we return error from the API or general one if we can't decode it
	type simpleErrorResponse struct {
		Err string `json:"error"`
	}
	var out simpleErrorResponse
	err := json.NewDecoder(res.Body).Decode(&out)
	if err == nil {
		return errors.New(out.Err)
	}

	return errors.New("RF API Error")
}

// Do sends out the request to the API and unpacks JSON response to the out variable.
func (c *Client) Do(req *http.Request, out interface{}) (*http.Response, error) {
	// Send out http request
	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	// Close the body after we're done with it, to allow connection reuse.
	defer func() {
		io.Copy(ioutil.Discard, res.Body)
		res.Body.Close()
	}()

	// We check response for potential errors and return them to the caller.
	// We do not nil the response, as a caller might want to inspect the response in case of an error.
	err = checkResponse(res)
	if err != nil {
		return res, err
	}

	// Here we check for the out pointer, and if it exists we unmarshall JSON there and return any
	// potential errors to the caller.
	if out != nil {
		err = json.NewDecoder(res.Body).Decode(out)
	}

	return res, err
}
