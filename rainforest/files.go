package rainforest

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/urfave/cli"
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

// CreateTestFile creates a UploadedFile resource by sending file information to
// Rainforest. This information is used for uploading the actual file to AWS.
func (c *Client) CreateTestFile(testID int, file *os.File) (*AWSFileInfo, error) {
	awsFileInfo := &AWSFileInfo{}
	fileName := file.Name()
	fileInfo, err := file.Stat()

	if err != nil {
		return awsFileInfo, err
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return awsFileInfo, err
	}

	md5CheckSum := md5.Sum(data)
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
		return awsFileInfo, err
	}

	_, err = c.Do(req, awsFileInfo)
	return awsFileInfo, err
}

// UploadTestFile is a function that uploads the actual file contents to AWS
func (c *Client) UploadTestFile(file *os.File, awsFileInfo *AWSFileInfo) error {
	fileName := filepath.Base(file.Name())
	fileExt := filepath.Ext(fileName)

	buffer := new(bytes.Buffer)
	writer := multipart.NewWriter(buffer)

	writer.WriteField("key", awsFileInfo.AWSKey)
	writer.WriteField("AWSAccessKeyId", awsFileInfo.AWSAccessID)
	writer.WriteField("acl", awsFileInfo.AWSACL)
	writer.WriteField("policy", awsFileInfo.AWSPolicy)
	writer.WriteField("signature", awsFileInfo.AWSSignature)
	writer.WriteField("Content-Type", mime.TypeByExtension(fileExt))

	// For some reason it's impossible to set Content-Type and Content-Transfer-Encoding
	// by simply using writer.CreateFormField, so doing it manually.
	part, err := createFormFile(writer, fileName)

	var contents []byte
	contents, err = ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	part.Write(contents)
	buffer.Write([]byte("--" + writer.Boundary() + "--"))

	x, _ := ioutil.ReadAll(buffer)
	y, _ := os.Create("test.txt")
	s := fmt.Sprintf("%q", x)
	y.Write([]byte(s))

	url := awsFileInfo.AWSURL

	var req *http.Request
	req, err = http.NewRequest("POST", url, buffer)

	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	var resp *http.Response
	resp, err = http.DefaultClient.Do(req)

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

		errMsg := fmt.Sprintf("There was an error uploading your file - %v: %v", file.Name(), string(body))
		return cli.NewExitError(errMsg, 1)
	}

	return nil
}

func createFormFile(w *multipart.Writer, fileName string) (io.Writer, error) {
	quoteEscaper := strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

	h := make(textproto.MIMEHeader)
	contentDisposition := fmt.Sprintf(`form-data; name="file"; filename="%s"`, quoteEscaper.Replace(fileName))
	h.Set("Content-Disposition", contentDisposition)
	contentType := mime.TypeByExtension(filepath.Ext(fileName))
	h.Set("Content-Type", contentType)
	h.Set("Content-Transfer-Encoding", "binary")
	return w.CreatePart(h)
}
