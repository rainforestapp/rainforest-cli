package rainforest

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func TestGetUploadedFiles(t *testing.T) {
	setup()
	defer cleanup()

	const reqMethod = "GET"
	const testID = 1337

	files := []UploadedFile{
		{ID: 123, Signature: "file_sig1", Digest: "digest1"},
		{ID: 456, Signature: "file_sig2", Digest: "digest2"},
	}

	mux.HandleFunc("/tests/"+strconv.Itoa(testID)+"/files", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != reqMethod {
			t.Errorf("Request method = %v, want %v", r.Method, reqMethod)
		}

		enc := json.NewEncoder(w)
		enc.Encode(files)
	})

	out, _ := client.GetUploadedFiles(testID)

	if !reflect.DeepEqual(files, out) {
		t.Errorf("Response expected = %v, actual %v", files, out)
	}
}

func TestMultipartFormRequest(t *testing.T) {
	Key := "awsKey"
	AccessID := "accessID"
	ACL := "acl"
	Policy := "policy"
	Signature := "signature"

	aws := AWSFileInfo{
		URL:       "https://f.rainforestqa.com/stuff",
		Key:       Key,
		AccessID:  AccessID,
		ACL:       ACL,
		Policy:    Policy,
		Signature: Signature,
	}

	fileName := "my_file.txt"
	fileContents := []byte("This is in my file")

	req, err := aws.MultipartFormRequest(fileName, fileContents)

	if err != nil {
		t.Error(err.Error())
	}

	var body []byte
	body, err = ioutil.ReadAll(req.Body)

	if req.URL.String() != aws.URL {
		t.Errorf("Incorrect URL. Have %v, want %v", req.URL, aws.URL)
	}

	if req.ContentLength != int64(len(body)) {
		t.Errorf("Incorrect ContentLength for request. Have %v, want %v.", req.ContentLength, len(body))
	}

	stringBody := string(body)

	fileExt := filepath.Ext(fileName)
	fields := map[string]string{
		"key":            Key,
		"AWSAccessKeyId": AccessID,
		"acl":            ACL,
		"policy":         Policy,
		"signature":      Signature,
		"Content-Type":   mime.TypeByExtension(fileExt),
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

	fileHeaderStr := fmt.Sprintf("Content-Disposition: form-data; name=\"file\"; filename=\"%v\"\r\n"+
		"Content-Type: application/octet-stream", fileName)

	if !strings.Contains(stringBody, fileHeaderStr) {
		t.Error("Incorrect file header in request body")
	}

	if !strings.Contains(stringBody, string(fileContents)) {
		t.Error("File contents not found in request body")
	}
}

func TestCreateTestFile(t *testing.T) {
	testID := 1001
	fileExt := ".txt"
	fileName := "my_file_name" + fileExt
	fileContents := []byte("my file contents")

	md5CheckSum := md5.Sum(fileContents)
	hexDigest := hex.EncodeToString(md5CheckSum[:16])

	file, err := os.Create(fileName)

	if err != nil {
		t.Fatal(err.Error())
	}

	defer func() {
		file.Close()
		os.Remove(fileName)
	}()

	_, err = file.Write(fileContents)

	if err != nil {
		t.Fatal(err.Error())
	}

	setup()
	defer cleanup()

	url := fmt.Sprintf("/tests/%v/files", testID)

	awsInfo := AWSFileInfo{
		FileID:        1234,
		FileSignature: "file signature",
		URL:           "https://f.rainforestqa.com/stuff",
		Key:           "tests/1234",
		AccessID:      "accessId",
		Policy:        "abc123",
		ACL:           "private",
		Signature:     "signature",
	}

	expectedRequestBody := UploadedFile{
		MimeType: mime.TypeByExtension(fileExt),
		Size:     int64(len(fileContents)),
		Name:     fileName,
		Digest:   hexDigest,
	}

	mux.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Request method = %v, want POST", r.Method)
		}

		defer r.Body.Close()

		out := &UploadedFile{}
		err = json.NewDecoder(r.Body).Decode(out)

		if err != nil {
			t.Errorf("Error decoding request body: %v", err.Error())
		}

		if !reflect.DeepEqual(*out, expectedRequestBody) {
			t.Errorf("Unexpected parameters.\nActual: %#v\nExpected: %#v", *out, expectedRequestBody)
		}

		enc := json.NewEncoder(w)
		enc.Encode(awsInfo)
	})

	out, err := client.CreateTestFile(testID, file, fileContents)

	if err != nil {
		t.Error(err.Error())
	} else if !reflect.DeepEqual(*out, awsInfo) {
		t.Errorf("Unexpected response from CreateTestFile.\nActual: %#v\nExpected: %#v", *out, awsInfo)
	}
}
