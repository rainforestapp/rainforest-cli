package rainforest

import (
	"crypto/md5"
	"io/ioutil"
	"mime"
	"os"
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

// GetUploadedFiles returns information for all all files uploaded to the
// given test before.
func (c *Client) GetUploadedFiles(fileID int) ([]UploadedFile, error) {
	req, err := c.NewRequest("GET", "tests/"+strconv.Itoa(fileID)+"/files", nil)
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
	AWSURL        string `json:"aws_url"`
	AWSKey        string `json:"aws_key"`
	AWSAccessID   string `json:"aws_access_id"`
	AWSPolicy     string `json:"aws_policy"`
	AWSACL        string `json:"aws_acl"`
	AWSSignature  string `json:"aws_signature"`
}

// CreateFile creates a UploadedFile resource by sending file information to
// Rainforest. This information is used for uploading the actual file to AWS.
func (c *Client) CreateFile(testID int, file os.File) (*AWSFileInfo, error) {
	var awsFileInfo *AWSFileInfo
	fileInfo, err := os.Stat(file.Name())

	if err != nil {
		return awsFileInfo, err
	}

	data, err := ioutil.ReadAll(&file)
	if err != nil {
		return awsFileInfo, err
	}

	md5CheckSum := md5.Sum(data)

	body := UploadedFile{
		MimeType: mime.TypeByExtension(fileInfo.Name()),
		Size:     fileInfo.Size(),
		Name:     fileInfo.Name(),
		Digest:   string(md5CheckSum[:16]),
	}

	url := "tests/" + strconv.Itoa(testID) + "/files"
	req, err := c.NewRequest("POST", url, body)
	if err != nil {
		return awsFileInfo, err
	}

	_, err = c.Do(req, awsFileInfo)
	return awsFileInfo, err
}
