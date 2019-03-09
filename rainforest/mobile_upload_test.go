package rainforest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

func TestGetPresignedPOST(t *testing.T) {
	setup()
	defer cleanup()

	const reqMethod = "GET"

	siteID := 123
	environmentID := 456
	extension := ".zip"

	Key := "awsKey"
	AccessID := "accessID"
	ACL := "acl"
	Policy := "policy"
	Signature := "signature"

	requiredFields := map[string]string{
		"key":            Key,
		"AWSAccessKeyId": AccessID,
		"acl":            ACL,
		"policy":         Policy,
		"signature":      Signature,
	}
	const awsURL = "https://s3.aws.com/stuff"
	const rainforestURL = "https://private-apps.rainforestqa.com/123/456.zip"
	respBody := RFPresignedPostData{
		URL: awsURL, RequiredFields: requiredFields, RainforestURL: rainforestURL,
	}

	mux.HandleFunc("/uploads", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != reqMethod {
			t.Errorf("Request method = %v, want %v", r.Method, reqMethod)
		}
		jsonData, _ := json.Marshal(respBody)
		jsonString := string(jsonData)
		fmt.Fprint(w, jsonString)
	})

	presignedPost, err := client.GetPresignedPOST(extension, siteID, environmentID)
	if err != nil {
		t.Errorf(err.Error())
	}
	if presignedPost.URL != awsURL {
		t.Errorf("Incorrect aws URL returned on PresignedPost data")
	}
	if presignedPost.RainforestURL != rainforestURL {
		t.Errorf("Incorrect rainforest URL returned on PresignedPost data")
	}
	if !reflect.DeepEqual(presignedPost.RequiredFields, requiredFields) {
		t.Errorf("Incorrect required fields returned")
	}
}

func TestUploadToS3(t *testing.T) {
	awsMUX := http.NewServeMux()
	fakeAWSServer := httptest.NewServer(awsMUX)
	awsURL, _ := url.Parse(fakeAWSServer.URL)

	Key := "awsKey"
	AccessID := "accessID"
	ACL := "acl"
	Policy := "policy"
	Signature := "signature"

	requiredFields := map[string]string{
		"key":            Key,
		"AWSAccessKeyId": AccessID,
		"acl":            ACL,
		"policy":         Policy,
		"signature":      Signature,
	}

	const reqMethod = "POST"

	postData := RFPresignedPostData{
		URL:            awsURL.String() + "/stuff",
		RequiredFields: requiredFields,
		RainforestURL:  "https://private-apps.rainforestqa.com/1234/5678/910.ipa",
	}

	awsMUX.HandleFunc("/stuff", func(w http.ResponseWriter, req *http.Request) {
		if req.Method != reqMethod {
			t.Errorf("Request method = %v, want %v", req.Method, reqMethod)
		}

		var body []byte
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			t.Error(err.Error())
		}

		if req.URL.String() != "/stuff" {
			t.Errorf("Incorrect URL. Have %v, want %v", req.URL, postData.URL)
		}

		if req.ContentLength != int64(len(body)) {
			t.Errorf("Incorrect ContentLength for request. Have %v, want %v.", req.ContentLength, len(body))
		}

		stringBody := string(body)
		fileName := "testfile.txt"
		testString := "testing123"
		fileHeaderStr := fmt.Sprintf("Content-Disposition: form-data; name=\"file\"; filename=\"%v\"\r\n"+
			"Content-Type: application/octet-stream", fileName)
		if !strings.Contains(stringBody, fileHeaderStr) {
			t.Error("Incorrect file header in request body")
		}
		if !strings.Contains(stringBody, testString) {
			t.Errorf("Test string not found in body. Have %v, want %v", stringBody, testString)
		}

		fields := map[string]string{
			"key":            Key,
			"AWSAccessKeyId": AccessID,
			"acl":            ACL,
			"policy":         Policy,
			"signature":      Signature,
		}

		for k, v := range fields {
			keyStr := fmt.Sprintf("Content-Disposition: form-data; name=\"%v\"", k)

			if !strings.Contains(stringBody, keyStr) {
				t.Errorf("Required field not found in request body: %v", keyStr)
			}

			if !strings.Contains(stringBody, v) {
				t.Errorf("Required value not found in request body for %v: %v", k, v)
			}
		}

		fmt.Fprint(w, "OK")
	})

	err := client.UploadToS3(&postData, "test/testfile.txt")

	if err != nil {
		t.Error(err.Error())
	}

	fakeAWSServer.Close()
	awsMUX = http.NewServeMux()
	fakeAWSServer = httptest.NewServer(awsMUX)
	awsURL, _ = url.Parse(fakeAWSServer.URL)

	awsMUX.HandleFunc("/stuff", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Something bad happened!"))
	})

	postData.URL = awsURL.String() + "/stuff"

	err = client.UploadToS3(&postData, "test/testfile.txt")

	test := err.Error()
	if !strings.Contains(test, "500") || !strings.Contains(test, "There was an error uploading your file") {
		t.Error("Not raising error when AWS server returns bad things.")
	}
	fakeAWSServer.Close()
}

func TestUpdateURL(t *testing.T) {
	setup()
	defer cleanup()

	siteID := 123
	environmentID := 456

	respBody := SiteEnvironmentsData{
		SiteEnvironments: []SiteEnvironment{
			SiteEnvironment{
				ID:            1,
				SiteID:        siteID,
				EnvironmentID: environmentID,
				URL:           "https://oldurl.com/app1.zip",
			},
			SiteEnvironment{
				ID:            2,
				SiteID:        654,
				EnvironmentID: 321,
				URL:           "https://oldurl.com/app2.zip",
			},
		},
	}

	newURL := "https://newurl.com/app3.zip"

	mux.HandleFunc("/site_environments", func(w http.ResponseWriter, r *http.Request) {
		const reqMethod = "GET"
		if r.Method != reqMethod {
			t.Errorf("Request method = %v, want %v", r.Method, reqMethod)
		}
		jsonData, _ := json.Marshal(respBody)
		jsonString := string(jsonData)
		fmt.Fprint(w, jsonString)
	})

	mux.HandleFunc("/site_environments/1", func(w http.ResponseWriter, r *http.Request) {
		const reqMethod = "PUT"
		if r.Method != reqMethod {
			t.Errorf("Request method = %v, want %v", r.Method, reqMethod)
		}

		decoder := json.NewDecoder(r.Body)
		var data SiteEnvironmentUpdate
		err := decoder.Decode(&data)
		if err != nil {
			t.Errorf(err.Error())
		}
		if data.URL != newURL {
			t.Errorf("New URL not found in body of put")
		}
		fmt.Fprint(w, "OK")
	})

	err := client.UpdateURL(siteID, environmentID, newURL)
	if err != nil {
		t.Errorf(err.Error())
	}
}
