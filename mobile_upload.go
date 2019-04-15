package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/rainforestapp/rainforest-cli/rainforest"
	"github.com/urfave/cli"
)

// mobileUploadAPI is part of the API connected to mobile uploads
type mobileUploadAPI interface {
	GetPresignedPOST(fileExt string, siteID int, environmentID int, appSlot int) (*rainforest.RFPresignedPostData, error)
	UploadToS3(postData *rainforest.RFPresignedPostData, filePath string) error
	UpdateURL(siteID int, environmentID int, appSlot int, newURL string) error
}

// uploadMobileApp takes a path to a mobile app file and uploads it to S3 then sets
// the site-id specified's URL to the magic url
func uploadMobileApp(api mobileUploadAPI, filePath string, siteID int, environmentID int, appSlot int) error {

	presignedPostData, err := api.GetPresignedPOST(filepath.Ext(filePath), siteID, environmentID, appSlot)
	if err != nil {
		return err
	}

	err = api.UploadToS3(presignedPostData, filePath)
	if err != nil {
		return err
	}

	err = api.UpdateURL(siteID, environmentID, appSlot, presignedPostData.RainforestURL)
	if err != nil {
		return err
	}

	return nil
}

func isAllowedExtension(extension string) bool {
	switch extension {
	case ".apk", ".ipa", ".zip", ".gz", ".tar.gz":
		return true
	}
	return false
}

const allowedExtensionsPretty = ".apk, .ipa, .zip, .gz, .tar.gz"

// mobileAppUpload is a wrapper around uploadMobileApp to function with mobile-upload cli command
func mobileAppUpload(c cliContext, api mobileUploadAPI) error {
	// Get the filepath from the arg
	filePath := c.Args().First()
	if filePath == "" {
		return cli.NewExitError("Mobile app file path not specified", 1)
	}

	// verify extension
	fileExt := strings.ToLower(filepath.Ext(filePath))
	if !isAllowedExtension(fileExt) {
		return cli.NewExitError(fmt.Sprintf("Invalid file extension. - %v. Allowed Extensions: %v", fileExt, allowedExtensionsPretty), 1)
	}

	siteIDString := c.String("site-id")
	if siteIDString == "" {
		return cli.NewExitError("site-id flag required", 1)
	}
	siteID, err := strconv.Atoi(siteIDString)
	if err != nil {
		return cli.NewExitError("site-id must be an integer", 1)
	}

	envIDstring := c.String("environment-id")
	if envIDstring == "" {
		return cli.NewExitError("environment-id flag required", 1)
	}
	environmentID, err := strconv.Atoi(envIDstring)
	if err != nil {
		return cli.NewExitError("environment-id must be an integer", 1)
	}

	appSlot := 1 // Default to 1, optional param
	appSlotString := c.String("app-slot")
	if appSlotString != "" {
		appSlot, err = strconv.Atoi(appSlotString)
		if err != nil || appSlot < 1 || appSlot > 5 {
			return cli.NewExitError("app-slot must be an integer (1 to 5)", 1)
		}
	}

	// Open app and return early with an error if we fail
	f, err := os.Open(filePath)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	defer f.Close()

	err = uploadMobileApp(api, filePath, siteID, environmentID, appSlot)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	return nil
}
