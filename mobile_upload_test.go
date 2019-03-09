package main

import (
	"strings"
	"testing"

	"github.com/rainforestapp/rainforest-cli/rainforest"
	"github.com/urfave/cli"
)

const (
	testMobileAppPath = "test/testing.zip"
)

type fakeMobileUploadAPI struct {
	getPresignedPOST      func(fileExt string, siteID int, environmentID int) (*rainforest.RFPresignedPostData, error)
	uploadToS3            func(postData *rainforest.RFPresignedPostData, filePath string) error
	setSiteEnvironmentURL func(siteID int, environmentID int, newURL string) error
}

func (f fakeMobileUploadAPI) GetPresignedPOST(fileExt string, siteID int, environmentID int) (*rainforest.RFPresignedPostData, error) {
	if f.getPresignedPOST != nil {
		return f.getPresignedPOST(fileExt, siteID, environmentID)
	}
	return nil, nil
}

func (f fakeMobileUploadAPI) UploadToS3(postData *rainforest.RFPresignedPostData, filePath string) error {
	if f.uploadToS3 != nil {
		return f.uploadToS3(postData, filePath)
	}
	return nil
}

func (f fakeMobileUploadAPI) UpdateURL(siteID int, environmentID int, newURL string) error {
	if f.setSiteEnvironmentURL != nil {
		return f.setSiteEnvironmentURL(siteID, environmentID, newURL)
	}
	return nil
}

func TestUploadMobileApp(t *testing.T) {
	siteID := 123
	environmentID := 456

	callCount := make(map[string]int)
	f := fakeMobileUploadAPI{
		getPresignedPOST: func(fileExt string, siteID int, environmentID int) (*rainforest.RFPresignedPostData, error) {
			callCount["getPresignedPOST"] = callCount["getPresignedPOST"] + 1
			return &rainforest.RFPresignedPostData{}, nil
		},
		uploadToS3: func(postData *rainforest.RFPresignedPostData, filePath string) error {
			callCount["uploadToS3"] = callCount["uploadToS3"] + 1
			return nil
		},
		setSiteEnvironmentURL: func(siteID int, environmentID int, newURL string) error {
			callCount["setSiteEnvironmentURL"] = callCount["setSiteEnvironmentURL"] + 1
			return nil
		},
	}

	err := uploadMobileApp(f, testMobileAppPath, siteID, environmentID)
	if err != nil {
		t.Errorf("Unexpected error: %v", err.Error())
	}
	if expected := 1; callCount["uploadToS3"] != expected {
		t.Errorf("api.uploadToS3 called invalid number of times: %v, expected %v", callCount["uploadToS3"], expected)
	}
	if expected := 1; callCount["getPresignedPOST"] != expected {
		t.Errorf("api.getPresignedPOST called invalid number of times: %v, expected %v", callCount["getPresignedPOST"], expected)
	}
	if expected := 1; callCount["setSiteEnvironmentURL"] != expected {
		t.Errorf("api.setSiteEnvironmentURL called invalid number of times: %v, expected %v", callCount["setSiteEnvironmentURL"], expected)
	}
}
func TestMobileAppUpload(t *testing.T) {
	siteID := "123"
	environmentID := "456"

	callCount := make(map[string]int)
	f := fakeMobileUploadAPI{
		getPresignedPOST: func(fileExt string, siteID int, environmentID int) (*rainforest.RFPresignedPostData, error) {
			callCount["getPresignedPOST"] = callCount["getPresignedPOST"] + 1
			return &rainforest.RFPresignedPostData{}, nil
		},
		uploadToS3: func(postData *rainforest.RFPresignedPostData, filePath string) error {
			callCount["uploadToS3"] = callCount["uploadToS3"] + 1
			return nil
		},
		setSiteEnvironmentURL: func(siteID int, environmentID int, newURL string) error {
			callCount["setSiteEnvironmentURL"] = callCount["setSiteEnvironmentURL"] + 1
			return nil
		},
	}
	fakeContext := newFakeContext(map[string]interface{}{
		"site-id":        siteID,
		"environment-id": environmentID,
	}, cli.Args{testMobileAppPath})

	err := mobileAppUpload(fakeContext, f)
	if err != nil {
		t.Errorf("Unexpected error: %v", err.Error())
	}
	if expected := 1; callCount["uploadToS3"] != expected {
		t.Errorf("api.uploadToS3 called invalid number of times: %v, expected %v", callCount["uploadToS3"], expected)
	}
	if expected := 1; callCount["getPresignedPOST"] != expected {
		t.Errorf("api.getPresignedPOST called invalid number of times: %v, expected %v", callCount["getPresignedPOST"], expected)
	}
	if expected := 1; callCount["setSiteEnvironmentURL"] != expected {
		t.Errorf("api.setSiteEnvironmentURL called invalid number of times: %v, expected %v", callCount["setSiteEnvironmentURL"], expected)
	}

	// Bad extension
	fakeContext = newFakeContext(map[string]interface{}{
		"site-id":        siteID,
		"environment-id": environmentID,
	}, cli.Args{"./file.exe"})
	err = mobileAppUpload(fakeContext, f)
	if _, ok := err.(*cli.ExitError); !ok &&
		!strings.Contains(err.Error(), "Invalid file extension") {
		t.Errorf("Not erroring on invalid extension.")
	}

	// missing envorionment-id flag
	fakeContext = newFakeContext(map[string]interface{}{
		"site-id": siteID,
	}, cli.Args{"./file.zip"})
	err = mobileAppUpload(fakeContext, f)
	if _, ok := err.(*cli.ExitError); !ok &&
		!strings.Contains(err.Error(), "environment-id flag required") {
		t.Errorf("Not erroring on missing environment-id.")
	}

	// missing envorionment-id flag
	fakeContext = newFakeContext(map[string]interface{}{
		"site-id":         siteID,
		"envorionment-id": "test",
	}, cli.Args{"./file.zip"})
	err = mobileAppUpload(fakeContext, f)
	if _, ok := err.(*cli.ExitError); !ok &&
		!strings.Contains(err.Error(), "environment-id flag required") {
		t.Errorf("envorionment-id must be an integer.")
	}

	// missing site-id flag
	fakeContext = newFakeContext(map[string]interface{}{
		"environment-id": environmentID,
	}, cli.Args{"./file.zip"})
	err = mobileAppUpload(fakeContext, f)
	if _, ok := err.(*cli.ExitError); !ok &&
		!strings.Contains(err.Error(), "site-id flag required") {
		t.Errorf("Not erroring on missing site-id.")
	}

	// missing site-id flag
	fakeContext = newFakeContext(map[string]interface{}{
		"site-id":        "test",
		"environment-id": environmentID,
	}, cli.Args{"./file.zip"})
	err = mobileAppUpload(fakeContext, f)
	if _, ok := err.(*cli.ExitError); !ok &&
		!strings.Contains(err.Error(), "site-id must be an integer") {
		t.Errorf("Not erroring on missing site-id.")
	}

	// missing filepath args
	fakeContext = newFakeContext(map[string]interface{}{
		"site-id":        siteID,
		"environment-id": environmentID,
	}, cli.Args{""})
	err = mobileAppUpload(fakeContext, f)
	if _, ok := err.(*cli.ExitError); !ok &&
		!strings.Contains(err.Error(), "Mobile app filepath not specified") {
		t.Errorf("Should have errored on missing filepath")
	}

	// bad filepath
	fakeContext = newFakeContext(map[string]interface{}{
		"site-id":        siteID,
		"environment-id": environmentID,
	}, cli.Args{"./bad_file_name.zip"})
	err = mobileAppUpload(fakeContext, f)
	if _, ok := err.(*cli.ExitError); !ok &&
		!strings.Contains(err.Error(), "no such file or directory") {
		t.Errorf("Should have errored on missing filepath")
	}
}
