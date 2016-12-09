package rainforest

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
	"reflect"
	"strconv"
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

type fakeOSFile struct {
	name  string
	stats OSFileInfo
}

func (f *fakeOSFile) Name() string {
	return f.name
}

func (f *fakeOSFile) Stat() (OSFileInfo, error) {
	return f.stats, nil
}

type fakeOSFileInfo struct {
	size int64
}

func (fi *fakeOSFileInfo) Size() int64 {
	return fi.size
}

func TestCreateTestFile(t *testing.T) {
	testID := 1001
	fileSize := int64(1337)
	fileExt := ".txt"
	fileName := "files/my_file_name" + fileExt
	fileContents := []byte("my file contents")

	md5CheckSum := md5.Sum(fileContents)
	hexDigest := hex.EncodeToString(md5CheckSum[:16])

	fileInfo := fakeOSFileInfo{size: fileSize}
	file := fakeOSFile{name: fileName, stats: &fileInfo}

	setup()
	defer cleanup()

	url := fmt.Sprintf("/tests/%v/files", testID)

	awsInfo := AWSFileInfo{
		FileID:        1234,
		FileSignature: "file signature",
		AWSURL:        "https://f.rainforestqa.com/stuff",
		AWSKey:        "tests/1234",
		AWSAccessID:   "accessId",
		AWSPolicy:     "abc123",
		AWSACL:        "private",
		AWSSignature:  "signature",
	}

	expectedRequestBody := UploadedFile{
		MimeType: mime.TypeByExtension(fileExt),
		Size:     fileSize,
		Name:     fileName,
		Digest:   hexDigest,
	}

	mux.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Request method = %v, want POST", r.Method)
		}

		defer r.Body.Close()

		out := &UploadedFile{}
		err := json.NewDecoder(r.Body).Decode(out)

		if err != nil {
			t.Errorf("Error decoding request body: %v", err.Error())
		}

		if !reflect.DeepEqual(*out, expectedRequestBody) {
			t.Errorf("Unexpected parameters.\nActual: %#v\nExpected: %#v", *out, expectedRequestBody)
		}

		enc := json.NewEncoder(w)
		enc.Encode(awsInfo)
	})

	out, err := client.CreateTestFile(testID, &file, fileContents)

	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(*out, awsInfo) {
		t.Errorf("Unexpected response from CreateTestFile.\nActual: %#v\nExpected: %#v", *out, awsInfo)
	}
}
