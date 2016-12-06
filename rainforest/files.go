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

// CreateFile creates a UploadedFile resource by sending file information to
// Rainforest. This information is used for uploading the actual file to AWS.
func (c *Client) CreateFile(testID int, file os.File) (UploadedFile, error) {
	var fileResp UploadedFile
	// fileName := file.Name()
	fileInfo, err := os.Stat(file.Name())

	if err != nil {
		return fileResp, err
	}

	data, err := ioutil.ReadAll(&file)
	if err != nil {
		return fileResp, err
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
		return fileResp, err
	}

	_, err = c.Do(req, &fileResp)
	return UploadedFile{}, err
}
