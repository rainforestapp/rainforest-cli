package rainforest

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

// UploadedFile represents a file that has been uploaded to Rainforest
type UploadedFile struct {
	ID        int    `json:"id"`
	Signature string `json:"signature"`
	Digest    string `json:"digest"`
	MimeType  string `json:"mime_type"`
	Size      int64  `json:"size"`
	Name      string `json:"name"`
}

// getUploadedFiles returns information for all all files uploaded to the
// given test before.
func (c *Client) getUploadedFiles(testID int) ([]UploadedFile, error) {
	req, err := c.NewRequest("GET", "tests/"+strconv.Itoa(testID)+"/files", nil)
	if err != nil {
		return nil, err
	}

	var fileResp []UploadedFile
	_, err = c.Do(req, &fileResp)
	return fileResp, err
}

// AWSFileInfo represents the response when uploading new file data to Rainforest.
// It contains information used to upload data the file to AWS.
type AWSFileInfo struct {
	FileID        int    `json:"file_id"`
	FileSignature string `json:"file_signature"`
	URL           string `json:"aws_url"`
	Key           string `json:"aws_key"`
	AccessID      string `json:"aws_access_id"`
	Policy        string `json:"aws_policy"`
	ACL           string `json:"aws_acl"`
	Signature     string `json:"aws_signature"`
}

// multipartFormRequest creates a http.Request containing the required body for
// uploading a file to AWS given the values stored in the receiving AWSFileInfo struct.
func (aws *AWSFileInfo) multipartFormRequest(fileName string, fileContents []byte) (*http.Request, error) {
	var req *http.Request
	fileExt := filepath.Ext(fileName)

	buffer := new(bytes.Buffer)
	writer := multipart.NewWriter(buffer)

	writer.WriteField("key", aws.Key)
	writer.WriteField("AWSAccessKeyId", aws.AccessID)
	writer.WriteField("acl", aws.ACL)
	writer.WriteField("policy", aws.Policy)
	writer.WriteField("signature", aws.Signature)
	writer.WriteField("Content-Type", mime.TypeByExtension(fileExt))

	part, err := writer.CreateFormFile("file", fileName)
	part.Write(fileContents)

	url := aws.URL
	req, err = http.NewRequest("POST", url, buffer)
	if err != nil {
		return req, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	writer.Close()
	req.ContentLength = int64(buffer.Len())

	return req, nil
}

// createTestFile creates a UploadedFile resource by sending file information to
// Rainforest. This information is used for uploading the file contents to AWS.
func (c *Client) createTestFile(testID int, file *os.File, fileContents []byte) (*AWSFileInfo, error) {
	fileName := file.Name()
	fileInfo, err := file.Stat()

	if err != nil {
		return &AWSFileInfo{}, err
	}

	md5CheckSum := md5.Sum(fileContents)
	hexDigest := hex.EncodeToString(md5CheckSum[:16])

	body := UploadedFile{
		MimeType: mime.TypeByExtension(filepath.Ext(fileName)),
		Size:     fileInfo.Size(),
		Name:     fileName,
		Digest:   hexDigest,
	}

	url := "tests/" + strconv.Itoa(testID) + "/files"
	req, err := c.NewRequest("POST", url, body)
	if err != nil {
		return &AWSFileInfo{}, err
	}

	awsFileInfo := &AWSFileInfo{}
	_, err = c.Do(req, awsFileInfo)
	return awsFileInfo, err
}

// uploadTestFile is a function that uploads the file contents to AWS
func (c *Client) uploadTestFile(fileName string, fileContents []byte, awsFileInfo *AWSFileInfo) error {
	req, err := awsFileInfo.multipartFormRequest(fileName, fileContents)

	var resp *http.Response
	resp, err = c.client.Do(req)

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

		return fmt.Errorf("There was an error uploading your file - %v: %v", fileName, string(body))
	}

	return nil
}
