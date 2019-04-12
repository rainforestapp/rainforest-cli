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
	getPresignedPOST      func(fileExt string, siteID int, environmentID int, appSlot int) (*rainforest.RFPresignedPostData, error)
	uploadToS3            func(postData *rainforest.RFPresignedPostData, filePath string) error
	setSiteEnvironmentURL func(siteID int, environmentID int, appSlot int, newURL string) error
}

func (f fakeMobileUploadAPI) GetPresignedPOST(fileExt string, siteID int, environmentID int, appSlot int) (*rainforest.RFPresignedPostData, error) {
	if f.getPresignedPOST != nil {
		return f.getPresignedPOST(fileExt, siteID, environmentID, appSlot)
	}
	return nil, nil
}

func (f fakeMobileUploadAPI) UploadToS3(postData *rainforest.RFPresignedPostData, filePath string) error {
	if f.uploadToS3 != nil {
		return f.uploadToS3(postData, filePath)
	}
	return nil
}

func (f fakeMobileUploadAPI) UpdateURL(siteID int, environmentID int, appSlot int, newURL string) error {
	if f.setSiteEnvironmentURL != nil {
		return f.setSiteEnvironmentURL(siteID, environmentID, appSlot, newURL)
	}
	return nil
}

func TestUploadMobileApp(t *testing.T) {
	siteID := 123
	environmentID := 456
	appSlot := 2

	callCount := make(map[string]int)
	f := fakeMobileUploadAPI{
		getPresignedPOST: func(fileExt string, siteID int, environmentID int, appSlot int) (*rainforest.RFPresignedPostData, error) {
			callCount["getPresignedPOST"] = callCount["getPresignedPOST"] + 1
			return &rainforest.RFPresignedPostData{}, nil
		},
		uploadToS3: func(postData *rainforest.RFPresignedPostData, filePath string) error {
			callCount["uploadToS3"] = callCount["uploadToS3"] + 1
			return nil
		},
		setSiteEnvironmentURL: func(siteID int, environmentID int, appSlot int, newURL string) error {
			callCount["setSiteEnvironmentURL"] = callCount["setSiteEnvironmentURL"] + 1
			return nil
		},
	}

	err := uploadMobileApp(f, testMobileAppPath, siteID, environmentID, appSlot)
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
	appSlot := "2"

	callCount := make(map[string]int)
	f := fakeMobileUploadAPI{
		getPresignedPOST: func(fileExt string, siteID int, environmentID int, appSlot int) (*rainforest.RFPresignedPostData, error) {
			callCount["getPresignedPOST"] = callCount["getPresignedPOST"] + 1
			return &rainforest.RFPresignedPostData{}, nil
		},
		uploadToS3: func(postData *rainforest.RFPresignedPostData, filePath string) error {
			callCount["uploadToS3"] = callCount["uploadToS3"] + 1
			return nil
		},
		setSiteEnvironmentURL: func(siteID int, environmentID int, appSlot int, newURL string) error {
			callCount["setSiteEnvironmentURL"] = callCount["setSiteEnvironmentURL"] + 1
			return nil
		},
	}
	fakeContext := newFakeContext(map[string]interface{}{
		"site-id":        siteID,
		"environment-id": environmentID,
		"app-slot":       appSlot,
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

	// non-int envorionment-id flag
	fakeContext = newFakeContext(map[string]interface{}{
		"site-id":         siteID,
		"envorionment-id": "test",
	}, cli.Args{"./file.zip"})
	err = mobileAppUpload(fakeContext, f)
	if _, ok := err.(*cli.ExitError); !ok &&
		!strings.Contains(err.Error(), "environment-id must be an integer") {
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

	// non-int site-id flag
	fakeContext = newFakeContext(map[string]interface{}{
		"site-id":        "test",
		"environment-id": environmentID,
	}, cli.Args{"./file.zip"})
	err = mobileAppUpload(fakeContext, f)
	if _, ok := err.(*cli.ExitError); !ok &&
		!strings.Contains(err.Error(), "site-id must be an integer") {
		t.Errorf("Not erroring on missing site-id.")
	}

	// app-slot flag is set to a non-int
	fakeContext = newFakeContext(map[string]interface{}{
		"site-id":        siteID,
		"environment-id": environmentID,
		"app-slot":       "test",
	}, cli.Args{"./file.zip"})
	err = mobileAppUpload(fakeContext, f)
	if _, ok := err.(*cli.ExitError); !ok &&
		!strings.Contains(err.Error(), "app-slot must be an integer (1 to 5)") {
		t.Errorf("Not erroring on invalid app-slot.")
	}

	// app-slot flag is set to a non 1-5 int
	fakeContext = newFakeContext(map[string]interface{}{
		"site-id":        siteID,
		"environment-id": environmentID,
		"app-slot":       10,
	}, cli.Args{"./file.zip"})
	err = mobileAppUpload(fakeContext, f)
	if _, ok := err.(*cli.ExitError); !ok &&
		!strings.Contains(err.Error(), "app-slot must be an integer (1 to 5)") {
		t.Errorf("Not erroring on invalid app-slot.")
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
