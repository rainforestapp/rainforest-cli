package rainforest

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func TestReadAll(t *testing.T) {
	expectedSteps := []interface{}{
		RFTestStep{
			Action:   "First Action",
			Response: "First Question?",
			Redirect: true,
		},
		RFTestStep{
			Action:   "Second Action",
			Response: "Second Question?",
			Redirect: true,
		},
		RFEmbeddedTest{
			RFMLID:   "embedded_id",
			Redirect: true,
		},
	}

	expectedTestVals := RFTest{
		RFMLID:   "my_rfml_id",
		Title:    "my_title",
		StartURI: "/testing",
		SiteID:   12345,
		Tags:     []string{"foo", "bar"},
		Browsers: []string{"chrome", "firefox"},
		Steps:    expectedSteps,
	}

	testText := fmt.Sprintf(`#! %v
# title: %v
# start_uri: %v
# site_id: %v
# tags: %v
# browsers: %v

%v
%v

%v
%v

- %v`,
		expectedTestVals.RFMLID,
		expectedTestVals.Title,
		expectedTestVals.StartURI,
		expectedTestVals.SiteID,
		strings.Join(expectedTestVals.Tags, ", "),
		strings.Join(expectedTestVals.Browsers, ", "),
		expectedSteps[0].(RFTestStep).Action,
		expectedSteps[0].(RFTestStep).Response,
		expectedSteps[1].(RFTestStep).Action,
		expectedSteps[1].(RFTestStep).Response,
		expectedSteps[2].(RFEmbeddedTest).RFMLID,
	)

	r := strings.NewReader(testText)
	reader := NewRFMLReader(r)
	rfTest, err := reader.ReadAll()
	if err != nil {
		t.Fatal(err.Error())
	}

	if !reflect.DeepEqual(*rfTest, expectedTestVals) {
		t.Errorf("Incorrect values for RFTest.\nGot %#v\nWant %#v", rfTest, expectedTestVals)
	}

	// Two RFML IDs
	testText = fmt.Sprintf(`#! %v
# title: %v
# start_uri: %v
#! another_rfml_id

%v
%v`,
		expectedTestVals.RFMLID,
		expectedTestVals.Title,
		expectedTestVals.StartURI,
		expectedSteps[0].(RFTestStep).Action,
		expectedSteps[0].(RFTestStep).Response,
	)

	r = strings.NewReader(testText)
	reader = NewRFMLReader(r)
	_, err = reader.ReadAll()
	if err == nil {
		t.Fatal("Expected an error from ReadAll")
	} else if !strings.Contains(err.Error(), "line 4") {
		t.Errorf("Wrong line reported. Expected 4. Returned error: %v", err.Error())
	}

	// Comment with a colon
	expectedComment := "this_should: be a comment"
	testText = fmt.Sprintf(`#! %v
# title: %v
# start_uri: %v
# %v

%v
%v`,
		expectedTestVals.RFMLID,
		expectedTestVals.Title,
		expectedTestVals.StartURI,
		expectedComment,
		expectedSteps[0].(RFTestStep).Action,
		expectedSteps[0].(RFTestStep).Response,
	)

	r = strings.NewReader(testText)
	reader = NewRFMLReader(r)
	rfTest, err = reader.ReadAll()
	if err != nil {
		t.Fatal(err.Error())
	}

	if !strings.Contains(rfTest.Description, expectedComment) {
		t.Errorf("Description not found. Expected \"%v\". Description: %v", expectedComment, rfTest.Description)
	}
}

func TestWriteRFMLTest(t *testing.T) {
	var buffer bytes.Buffer
	writer := NewRFMLWriter(&buffer)

	// Just test the required metadata first
	rfmlID := "fake_rfml_id"
	title := "fake_title"
	startURI := "/path/to/nowhere"

	test := RFTest{
		RFMLID:   rfmlID,
		Title:    title,
		StartURI: startURI,
	}

	getOutput := func() string {
		writer.WriteRFMLTest(&test)

		raw, err := ioutil.ReadAll(&buffer)

		if err != nil {
			t.Fatal(err.Error())
		}
		return string(raw)
	}

	output := getOutput()

	mustHaves := []string{
		"#! " + rfmlID,
		"# title: " + title,
		"# start_uri: " + startURI,
	}

	for _, mustHave := range mustHaves {
		if !strings.Contains(output, mustHave) {
			t.Errorf("Missing expected string in writer output: %v", mustHave)
		}
	}

	mustNotHaves := []string{"site_id", "tags", "browsers"}

	for _, mustNotHave := range mustNotHaves {
		if strings.Contains(output, mustNotHave) {
			t.Errorf("%v found in writer output when omitted from RF test.", mustNotHave)
		}
	}

	// Now test all the headers
	buffer.Reset()

	siteID := 1989
	tags := []string{"foo", "bar"}
	browsers := []string{"chrome", "firefox"}
	description := "This is my description\nand it spans multiple\nlines!"

	test.SiteID = siteID
	test.Tags = tags
	test.Browsers = browsers
	test.Description = description

	output = getOutput()

	siteIDStr := "# site_id: " + strconv.Itoa(siteID)
	tagsStr := "# tags: " + strings.Join(tags, ", ")
	browsersStr := "# browsers: " + strings.Join(browsers, ", ")
	descStr := "# " + strings.Replace(description, "\n", "\n# ", -1)

	mustHaves = append(mustHaves, []string{siteIDStr, tagsStr, browsersStr, descStr}...)
	for _, mustHave := range mustHaves {
		if !strings.Contains(output, mustHave) {
			t.Errorf("Missing expected string in writer output: %v", mustHave)
		}
	}

	// Now test flattened steps
	buffer.Reset()

	firstStep := RFTestStep{
		Action:   "first action",
		Response: "first question",
		Redirect: true,
	}

	secondStep := RFTestStep{
		Action:   "second action",
		Response: "second question",
	}

	test.Steps = []interface{}{firstStep, secondStep}

	output = getOutput()

	expectedStepText := fmt.Sprintf("%v\n%v\n\n%v\n%v", firstStep.Action, firstStep.Response,
		secondStep.Action, secondStep.Response)
	if !strings.Contains(output, expectedStepText) {
		t.Error("Expected step text not found in writer output.")
		t.Logf("Output:\n%v", output)
		t.Logf("Expected:\n%v", expectedStepText)
	}

	// Test redirects for an embedded second step
	buffer.Reset()

	embeddedRFMLID := "embedded_test_rfml_id"
	embeddedStep := RFEmbeddedTest{RFMLID: embeddedRFMLID}

	test.Steps = []interface{}{firstStep, embeddedStep}

	output = getOutput()

	expectedStepText = fmt.Sprintf("%v\n%v\n\n# redirect: %v\n- %v", firstStep.Action, firstStep.Response,
		embeddedStep.Redirect, embeddedStep.RFMLID)
	if !strings.Contains(output, expectedStepText) {
		t.Error("Expected step text not found in writer output.")
		t.Logf("Output:\n%v", output)
		t.Logf("Expected:\n%v", expectedStepText)
	}

	// Test redirects for an embedded first step
	buffer.Reset()

	test.Steps = []interface{}{embeddedStep, firstStep}

	output = getOutput()

	expectedStepText = fmt.Sprintf("- %v\n\n# redirect: %v\n%v\n%v", embeddedStep.RFMLID,
		firstStep.Redirect, firstStep.Action, firstStep.Response)
	if !strings.Contains(output, expectedStepText) {
		t.Error("Expected step text not found in writer output.")
		t.Logf("Output:\n%v", output)
		t.Logf("Expected:\n%v", expectedStepText)
	}

	// Test redirects for a flat step after an embedded step that is not the first step
	buffer.Reset()

	test.Steps = []interface{}{firstStep, embeddedStep, secondStep}

	output = getOutput()

	expectedStepText = fmt.Sprintf("%v\n%v\n\n# redirect: %v\n- %v\n\n%v\n%v", firstStep.Action, firstStep.Response,
		embeddedStep.Redirect, embeddedStep.RFMLID, secondStep.Action, secondStep.Response)
	if !strings.Contains(output, expectedStepText) {
		t.Error("Expected step text not found in writer output.")
		t.Logf("Output:\n%v", output)
		t.Logf("Expected:\n%v", expectedStepText)
	}
}

func TestParseEmbeddedFiles(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err.Error())
	}

	test := RFTest{
		TestID: 5678,
		Steps: []interface{}{
			RFTestStep{
				Action:   "Embedding an existing screenshot {{file.screenshot(./assets/screenshot1.png)}}",
				Response: "Embedding a new screenshot {{ file.screenshot(./assets/screenshot2.png) }}",
			},
		},
		// Test does not exist, but this path is used to find the relative path to the
		// embedded files in the action and response.
		RFMLPath: filepath.Join(pwd, "./test/fake_test.rfml"),
	}

	existingFileID := 1234
	existingSignature := "existing_signature"
	newFileID := 4567
	newSignature := "new_signature"

	awsTestServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Response does not matter - just that the upload succeeded
		if r.Method != "POST" {
			t.Fatal("Unexpected request method to AWS")
			w.WriteHeader(http.StatusCreated)
		}
	}))
	defer awsTestServer.Close()
	awsURL := awsTestServer.URL

	setup()
	defer cleanup()

	// Set up response for uploaded files
	mux.HandleFunc(fmt.Sprintf("/tests/%v/files", test.TestID), func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		if r.Method == "GET" {
			testDir := filepath.Dir(test.RFMLPath)
			filePath := filepath.Join(testDir, "./assets/screenshot1.png")
			var file *os.File
			file, err = os.Open(filePath)
			if err != nil {
				t.Fatal(err.Error())
			}

			var data []byte
			data, err = ioutil.ReadAll(file)
			file.Close()
			if err != nil {
				t.Fatal(err.Error())
			}

			checksum := md5.Sum(data)
			fileDigest := hex.EncodeToString(checksum[:])

			files := []uploadedFile{
				{ID: existingFileID, Digest: fileDigest, Signature: existingSignature},
			}
			json.NewEncoder(w).Encode(files)
		} else if r.Method == "POST" {
			// Don't test body - already tested in files_test.go
			awsInfo := awsFileInfo{
				FileID:        newFileID,
				FileSignature: newSignature,
				URL:           awsURL,
				Key:           "key",
				AccessID:      "accessId",
				Policy:        "abc123",
				ACL:           "private",
				Signature:     "signature",
			}
			json.NewEncoder(w).Encode(awsInfo)
		}
	})

	err = client.ParseEmbeddedFiles(&test)
	if err != nil {
		t.Fatal(err.Error())
	}

	step := test.Steps[0].(RFTestStep)
	expectedStr := fmt.Sprintf("{{ file.screenshot(%v, %v) }}", existingFileID, existingSignature[0:6])
	if !strings.Contains(step.Action, expectedStr) {
		t.Errorf("Expected to find %v in %v", expectedStr, step.Action)
	}

	expectedStr = fmt.Sprintf("{{ file.screenshot(%v, %v) }}", newFileID, newSignature[0:6])
	if !strings.Contains(step.Response, expectedStr) {
		t.Errorf("Expected to find %v in %v", expectedStr, step.Response)
	}
}
