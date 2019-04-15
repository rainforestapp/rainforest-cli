package rainforest

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// RFPresignedPostData represents the url and required fields we must POST to AWS
// in order to upload
type RFPresignedPostData struct {
	URL            string            `json:"url"`
	RequiredFields map[string]string `json:"url_fields"`
	RainforestURL  string            `json:"rainforest_url"`
}

// GetPresignedPOST requests the presigned POST data from Rainforest so that we can upload the mobile
// app to S3
func (c *Client) GetPresignedPOST(fileExt string, siteID int, environmentID int, appSlot int) (*RFPresignedPostData, error) {
	var data RFPresignedPostData
	url := "uploads"
	req, err := c.NewRequest("GET", url, nil)
	q := req.URL.Query()
	q.Add("app_site_id", strconv.Itoa(siteID))
	q.Add("environment_id", strconv.Itoa(environmentID))
	q.Add("extension", fileExt)
	q.Add("app_slot", strconv.Itoa(appSlot))
	req.URL.RawQuery = q.Encode()

	_, err = c.Do(req, &data)

	return &data, err
}

// UploadToS3 creates a http.Request containing the required body for
// uploading a file to AWS given the values stored in the receiving awsFileInfo struct.
func (c *Client) UploadToS3(postData *RFPresignedPostData, filePath string) error {
	var req *http.Request
	fileName := filepath.Base(filePath)

	f, err := os.Open(filePath)
	defer f.Close()
	if err != nil {
		return err
	}

	fi, err := f.Stat()
	if err != nil {
		return err
	}
	fileSize := int64(fi.Size())

	contentLength := postData.emptyMultipartSize("file", fileName) + fileSize

	readBody, writeBody := io.Pipe()
	defer readBody.Close()

	writer := multipart.NewWriter(writeBody)

	// Do the writes async streamed
	errChan := make(chan error, 1)
	go func() {
		defer writeBody.Close()

		// Add the required fields from S3
		writeFields(postData.RequiredFields, writer)

		part, err := writer.CreateFormFile("file", fileName)
		if err != nil {
			errChan <- err
			return
		}
		if _, err := io.CopyN(part, f, fileSize); err != nil {
			errChan <- err
			return
		}
		errChan <- writer.Close()
	}()

	// Create the Request
	url := postData.URL
	req, err = http.NewRequest("POST", url, readBody)
	if err != nil {
		<-errChan
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.ContentLength = contentLength

	// Perform the upload
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		<-errChan
		return err
	}

	status := resp.StatusCode
	if status >= 300 {
		body, err := ioutil.ReadAll(resp.Body)

		<-errChan
		if err != nil {
			return err
		}

		return fmt.Errorf("There was an error uploading your file - %v: %v", fileName, string(body))
	}

	return nil
}

func writeFields(requiredFields map[string]string, writer *multipart.Writer) {
	for k, v := range requiredFields {
		writer.WriteField(k, v)
	}
}

func (postData *RFPresignedPostData) emptyMultipartSize(fieldname, filename string) int64 {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writeFields(postData.RequiredFields, writer)
	writer.CreateFormFile("file", filename)
	writer.Close()
	return int64(body.Len())
}

// UpdateURL fetches sites available to use during the RF runs.
func (c *Client) UpdateURL(siteID int, environmentID int, appSlot int, newURL string) error {
	siteEnvironment, err := c.getSiteEnvironment(siteID, environmentID)
	if err != nil {
		return err
	}

	index := appSlot - 1 // appSlot is 1-5, make it 0-index based
	splitURL := strings.Split(siteEnvironment.URL, "|")[:]
	if len(splitURL) < index+1 {
		newSplitURL := make([]string, index+1)
		copy(newSplitURL, splitURL)
		splitURL = newSplitURL
	}
	splitURL[index] = newURL
	updatedNewURL := strings.Join(splitURL, "|")

	err = c.setSiteEnvironmentURL(siteEnvironment.ID, updatedNewURL)
	if err != nil {
		return err
	}

	return nil
}

// SiteEnvironmentsData type represents the collection returned when calling the site_environments get endpoint
type SiteEnvironmentsData struct {
	SiteEnvironments []SiteEnvironment `json:"site_environments"`
}

// SiteEnvironment type represents a single SiteEnvironment returned by the API call for a list of sites
type SiteEnvironment struct {
	ID            int    `json:"id"`
	SiteID        int    `json:"site_id"`
	EnvironmentID int    `json:"environment_id"`
	URL           string `json:"url"`
}

func (c *Client) getSiteEnvironment(siteID int, environmentID int) (SiteEnvironment, error) {
	var siteEnvironment SiteEnvironment

	// Prepare request
	req, err := c.NewRequest("GET", "site_environments", nil)
	if err != nil {
		return siteEnvironment, err
	}

	// Send request and process response
	var resp SiteEnvironmentsData
	_, err = c.Do(req, &resp)
	if err != nil {
		return siteEnvironment, err
	}

	for _, siteEnvironment := range resp.SiteEnvironments {
		if siteEnvironment.SiteID == siteID && siteEnvironment.EnvironmentID == environmentID {
			return siteEnvironment, nil
		}
	}

	return siteEnvironment, fmt.Errorf("SiteEnvironment not found")
}

// SiteEnvironmentUpdate type is the body of site_environments PUT update for updating the URL
type SiteEnvironmentUpdate struct {
	URL string `json:"url"`
}

func (c *Client) setSiteEnvironmentURL(siteEnvironmentID int, newURL string) error {
	data := SiteEnvironmentUpdate{
		URL: newURL,
	}

	req, err := c.NewRequest("PUT", fmt.Sprintf("site_environments/%d", siteEnvironmentID), data)
	if err != nil {
		return err
	}

	resp, err := c.Do(req, nil)
	if err != nil {
		return err
	}

	status := resp.StatusCode
	if status >= 300 {
		var body []byte
		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		return fmt.Errorf("There was an error setting the mobile app's URL - %v", string(body))
	}
	return nil
}
