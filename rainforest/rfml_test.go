package rainforest

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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
	validSteps := []interface{}{
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

	validTestValues := RFTest{
		RFMLID:    "my_rfml_id",
		Title:     "my_title",
		StartURI:  "/testing",
		SiteID:    12345,
		FeatureID: 98765,
		State:     "enabled",
		Priority:  "P1",
		Tags:      []string{"foo", "bar"},
		Browsers:  []string{"chrome", "firefox"},
		Steps:     validSteps,
		Execute:   true,
	}

	testText := fmt.Sprintf(`#! %v
# title: %v
# start_uri: %v
# site_id: %v
# tags: %v
# browsers: %v
# feature_id: %v
# state: %v
# priority: %v

%v
%v

%v
%v

- %v`,
		validTestValues.RFMLID,
		validTestValues.Title,
		validTestValues.StartURI,
		validTestValues.SiteID,
		strings.Join(validTestValues.Tags, ", "),
		strings.Join(validTestValues.Browsers, ", "),
		validTestValues.FeatureID,
		validTestValues.State,
		validTestValues.Priority,
		validSteps[0].(RFTestStep).Action,
		validSteps[0].(RFTestStep).Response,
		validSteps[1].(RFTestStep).Action,
		validSteps[1].(RFTestStep).Response,
		validSteps[2].(RFEmbeddedTest).RFMLID,
	)

	r := strings.NewReader(testText)
	reader := NewRFMLReader(r)
	rfTest, err := reader.ReadAll()
	if err != nil {
		t.Fatal(err.Error())
	}

	if !reflect.DeepEqual(*rfTest, validTestValues) {
		t.Errorf("Incorrect values for RFTest.\nGot %#v\nWant %#v", rfTest, validTestValues)
	}

	// Test state is present
	testText = fmt.Sprintf(`#! %v
# title: %v
# start_uri: %v
# state: %v
`,
		validTestValues.RFMLID,
		validTestValues.Title,
		validTestValues.StartURI,
		"some_state",
	)

	r = strings.NewReader(testText)
	reader = NewRFMLReader(r)
	rfTest, err = reader.ReadAll()
	if err != nil {
		t.Fatal(err.Error())
	}

	if rfTest.State != "some_state" {
		t.Errorf("Incorrect test state. Got %v, Want %v", rfTest.State, "some_state")
	}

	// Test state is omitted
	testText = fmt.Sprintf(`#! %v
# title: %v
# start_uri: %v
`,
		validTestValues.RFMLID,
		validTestValues.Title,
		validTestValues.StartURI,
	)

	r = strings.NewReader(testText)
	reader = NewRFMLReader(r)
	rfTest, err = reader.ReadAll()
	if err != nil {
		t.Fatal(err.Error())
	}

	if rfTest.State != "enabled" {
		t.Errorf("Incorrect test state. Got %v, Want %v", rfTest.State, "enabled")
	}

	// Test priority is present
	testText = fmt.Sprintf(`#! %v
# title: %v
# start_uri: %v
# priority: %v
`,
		validTestValues.RFMLID,
		validTestValues.Title,
		validTestValues.StartURI,
		"P1",
	)

	r = strings.NewReader(testText)
	reader = NewRFMLReader(r)
	rfTest, err = reader.ReadAll()
	if err != nil {
		t.Fatal(err.Error())
	}

	if rfTest.Priority != "P1" {
		t.Errorf("Incorrect test priority. Got %v, Want %v", rfTest.Priority, "priority")
	}

	// Test priority is omitted
	testText = fmt.Sprintf(`#! %v
# title: %v
# start_uri: %v
`,
		validTestValues.RFMLID,
		validTestValues.Title,
		validTestValues.StartURI,
	)

	r = strings.NewReader(testText)
	reader = NewRFMLReader(r)
	rfTest, err = reader.ReadAll()
	if err != nil {
		t.Fatal(err.Error())
	}

	if rfTest.Priority != "" {
		t.Errorf("Incorrect test priority. Got %v, Want empty string", rfTest.Priority)
	}

	// Comment with a colon
	expectedComment := "this_should: be a comment"
	testText = fmt.Sprintf(`#! %v
# title: %v
# start_uri: %v
# %v

%v
%v`,
		validTestValues.RFMLID,
		validTestValues.Title,
		validTestValues.StartURI,
		expectedComment,
		validSteps[0].(RFTestStep).Action,
		validSteps[0].(RFTestStep).Response,
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

	// Non-executed test
	testText = fmt.Sprintf(`#! %v
# title: %v
# start_uri: %v
# execute: false

%v
%v`,
		validTestValues.RFMLID,
		validTestValues.Title,
		validTestValues.StartURI,
		validSteps[0].(RFTestStep).Action,
		validSteps[0].(RFTestStep).Response)
	r = strings.NewReader(testText)
	reader = NewRFMLReader(r)
	rfTest, err = reader.ReadAll()
	if err != nil {
		t.Error(err.Error())
	}

	if rfTest.Execute {
		t.Errorf("`execute: false` was not parsed correctly")
	}

	// missing RFML ID
	testText = fmt.Sprintf(`# title: %v
# start_uri: %v

%v
%v`,
		validTestValues.Title,
		validTestValues.StartURI,
		validSteps[0].(RFTestStep).Action,
		validSteps[0].(RFTestStep).Response,
	)

	r = strings.NewReader(testText)
	reader = NewRFMLReader(r)
	_, err = reader.ReadAll()
	if err == nil {
		t.Fatal("Expected an error from ReadAll")
	} else if !strings.Contains(err.Error(), "#!") {
		t.Errorf("Wrong error reported. Expected error for RFML ID field. Returned error: %v", err.Error())
	}

	// Two RFML IDs
	testText = fmt.Sprintf(`#! %v
# title: %v
# start_uri: %v
#! another_rfml_id

%v
%v`,
		validTestValues.RFMLID,
		validTestValues.Title,
		validTestValues.StartURI,
		validSteps[0].(RFTestStep).Action,
		validSteps[0].(RFTestStep).Response,
	)

	r = strings.NewReader(testText)
	reader = NewRFMLReader(r)
	_, err = reader.ReadAll()
	if err == nil {
		t.Fatal("Expected an error from ReadAll")
	} else if !strings.Contains(err.Error(), "line 4") {
		t.Errorf("Wrong line reported. Expected 4. Returned error: %v", err.Error())
	}

	// Missing Title
	testText = fmt.Sprintf(`#! %v
# start_uri: %v

%v
%v`,
		validTestValues.RFMLID,
		validTestValues.StartURI,
		validSteps[0].(RFTestStep).Action,
		validSteps[0].(RFTestStep).Response,
	)

	r = strings.NewReader(testText)
	reader = NewRFMLReader(r)
	_, err = reader.ReadAll()
	if err == nil {
		t.Fatal("Expected an error from ReadAll")
	} else if !strings.Contains(err.Error(), "# title") {
		t.Errorf("Wrong error reported. Expected error for title field. Returned error: %v", err.Error())
	}

	// empty feature_id, browser list, and tag list
	testText = fmt.Sprintf(`#! %v
# title: %v
# start_uri: %v
# browsers:
# tags:
# feature_id:

%v
%v`,
		validTestValues.RFMLID,
		validTestValues.Title,
		validTestValues.StartURI,
		validSteps[0].(RFTestStep).Action,
		validSteps[0].(RFTestStep).Response,
	)

	r = strings.NewReader(testText)
	reader = NewRFMLReader(r)
	rfTest, err = reader.ReadAll()
	if err != nil {
		t.Fatalf("Unexpected error from ReadAll: %v", err.Error())
	}

	if browserCount := len(rfTest.Browsers); browserCount != 0 {
		t.Errorf("Unexpected browsers, expected 0, got %v: %v", browserCount, rfTest.Browsers)
	}

	if tagCount := len(rfTest.Tags); tagCount != 0 {
		t.Errorf("Unexpected tags, expected 0, got %v: %v", tagCount, rfTest.Tags)
	}

	if rfTest.FeatureID != deleteFeature {
		t.Errorf("Unexpected feature ID, expected %v, got %v", deleteFeature, rfTest.FeatureID)
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
		Execute:  true,
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

	mustNotHaves := []string{"site_id", "tags", "browsers", "execute"}

	for _, mustNotHave := range mustNotHaves {
		if strings.Contains(output, mustNotHave) {
			t.Errorf("%v found in writer output when omitted from RF test.", mustNotHave)
		}
	}

	// Now test all the headers
	buffer.Reset()

	siteID := 1989
	featureID := 2017
	tags := []string{"foo", "bar"}
	browsers := []string{"chrome", "firefox"}
	description := "This is my description\nand it spans multiple\nlines!"

	test.SiteID = siteID
	test.Tags = tags
	test.Browsers = browsers
	test.Description = description
	test.FeatureID = FeatureIDInt(featureID)
	test.State = "disabled"
	test.Priority = "P1"

	output = getOutput()

	siteIDStr := "# site_id: " + strconv.Itoa(siteID)
	featureIDStr := "# feature_id: " + strconv.Itoa(featureID)
	tagsStr := "# tags: " + strings.Join(tags, ", ")
	browsersStr := "# browsers: " + strings.Join(browsers, ", ")
	descStr := "# " + strings.Replace(description, "\n", "\n# ", -1)
	stateStr := "# state: " + test.State
	priorityStr := "# priority: " + test.Priority

	mustHaves = append(mustHaves, []string{siteIDStr, featureIDStr, tagsStr, browsersStr, descStr, stateStr, priorityStr}...)
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

	// Test execute: false
	buffer.Reset()
	test.Execute = false
	output = getOutput()
	if !strings.Contains(output, "\n# execute: false") {
		t.Error("execute: false not found in writer output.")
		t.Logf("Output:\n%v", output)
	}

	// Test enabled state omitted
	buffer.Reset()
	test.State = "enabled"
	output = getOutput()
	if strings.Contains(output, "state") {
		t.Error("Test state field not expected to appear")
		t.Logf("Output:\n%v", output)
	}

	// Test empty state omitted (rainforest new)
	buffer.Reset()
	test.State = ""
	output = getOutput()
	if strings.Contains(output, "state") {
		t.Error("Test state field not expected to appear")
		t.Logf("Output:\n%v", output)
	}

	// Test draft state included
	buffer.Reset()
	test.State = "draft"
	output = getOutput()
	if !strings.Contains(output, "state") {
		t.Error("Test state expected to appear")
		t.Logf("Output:\n%v", output)
	}

	// Test priority omitted
	buffer.Reset()
	test.Priority = ""
	output = getOutput()
	if strings.Contains(output, "priority") {
		t.Error("Test priority field not expected to appear")
		t.Logf("Output:\n%v", output)
	}
}

func TestParseEmbeddedFiles(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err.Error())
	}

	testDir := filepath.Join(pwd, "./test")

	existingScreenshotPath := "./assets/screenshot1.png"
	newScreenshotPath := "./assets/screenshot2.png"
	nonexistentScreenshotPath := "./foo"
	existingDownloadPath := "./existing.txt"
	newDownloadPath := "./new.txt"
	nonexistentFilePath := "./bar"

	test := RFTest{
		TestID: 5678,
		Steps: []interface{}{
			RFTestStep{
				Action:   fmt.Sprintf("Embedding an existing screenshot {{file.screenshot(%v)}}", existingScreenshotPath),
				Response: fmt.Sprintf("Embedded an existing download {{ file.download(%v)}}", existingDownloadPath),
			},
			RFTestStep{
				Action:   fmt.Sprintf("Embedding a new screenshot {{ file.screenshot(%v) }}", newScreenshotPath),
				Response: fmt.Sprintf("Embedded a new download {{file.download(%v) }}", newDownloadPath),
			},
			RFTestStep{
				Action:   fmt.Sprintf("Embedding a non-existent screenshot {{  file.screenshot(%v) }}", nonexistentScreenshotPath),
				Response: fmt.Sprintf("Embedding a non-existent download {{ file.download(%v)  }}", nonexistentFilePath),
			},
		},
		// Test does not exist, but this path is used to find the relative path to the
		// embedded files in the action and response.
		RFMLPath: filepath.Join(testDir, "./fake_test.rfml"),
	}

	// Save existing screenshot digest
	var file *os.File
	file, err = os.Open(filepath.Join(testDir, existingScreenshotPath))
	if err != nil {
		t.Fatal(err.Error())
	}

	var contents []byte
	contents, err = ioutil.ReadAll(file)
	file.Close()
	if err != nil {
		t.Fatal(err.Error())
	}

	checksum := md5.Sum(contents)
	screenshotDigest := hex.EncodeToString(checksum[:])

	// Create existing downloaded file and save digest
	existingDownloadFile, err := os.Create(filepath.Join(testDir, existingDownloadPath))
	defer os.Remove(existingDownloadFile.Name())
	if err != nil {
		t.Fatal(err.Error())
	}
	contents = []byte("This has already been uploaded!")
	_, err = existingDownloadFile.Write(contents)
	if err != nil {
		t.Fatal(err.Error())
	}
	existingDownloadFile.Close()

	checksum = md5.Sum(contents)
	downloadDigest := hex.EncodeToString(checksum[:])

	// Create new downloaded file
	newDownloadFile, err := os.Create(filepath.Join(testDir, newDownloadPath))
	defer os.Remove(newDownloadFile.Name())
	if err != nil {
		t.Fatal(err.Error())
	}
	newDownloadFile.WriteString("This has not yet been uploaded!")
	newDownloadFile.Close()

	// Values retrieved from Rainforest API
	existingScreenshotID := 1234
	existingScreenshotSignature := "existing_screenshot_signature"
	newScreenshotID := 4567
	newScreenshotSignature := "new_screenshot_signature"

	existingDownloadID := 4321
	existingDownloadSignature := "existing_download_signature"
	newDownloadID := 7654
	newDownloadSignature := "new_download_signature"

	// Set up fake AWS server for uploads
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

	// Set up fake Rainforest server for GET and POST to /tests/:id/files
	mux.HandleFunc(fmt.Sprintf("/tests/%v/files", test.TestID), func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		switch r.Method {
		case "GET":
			files := []uploadedFile{
				{ID: existingScreenshotID, Digest: screenshotDigest, Signature: existingScreenshotSignature},
				{ID: existingDownloadID, Digest: downloadDigest, Signature: existingDownloadSignature},
			}
			json.NewEncoder(w).Encode(files)

		case "POST":
			var reqBody []byte
			reqBody, err = ioutil.ReadAll(r.Body)
			if err != nil {
				t.Fatal(err.Error())
			}

			var awsInfo awsFileInfo
			if bytes.Contains(reqBody, []byte(filepath.Base(newScreenshotPath))) {
				awsInfo = awsFileInfo{
					FileID:        newScreenshotID,
					FileSignature: newScreenshotSignature,
					URL:           awsURL,
					Key:           "key",
					AccessID:      "accessId",
					Policy:        "abc123",
					ACL:           "private",
					Signature:     "signature",
				}
			} else if bytes.Contains(reqBody, []byte(filepath.Base(newDownloadPath))) {
				awsInfo = awsFileInfo{
					FileID:        newDownloadID,
					FileSignature: newDownloadSignature,
					URL:           awsURL,
					Key:           "key",
					AccessID:      "accessId",
					Policy:        "abc123",
					ACL:           "private",
					Signature:     "signature",
				}
			}

			json.NewEncoder(w).Encode(awsInfo)
		}
	})

	out := &bytes.Buffer{}
	log.SetOutput(out)
	defer log.SetOutput(os.Stdout)

	err = client.ParseEmbeddedFiles(&test)
	if err != nil {
		t.Fatal(err.Error())
	}

	// Check screenshot values
	step := test.Steps[0].(RFTestStep)
	expectedStr := fmt.Sprintf("{{ file.screenshot(%v, %v) }}", existingScreenshotID, existingScreenshotSignature[0:6])
	if !strings.Contains(step.Action, expectedStr) {
		t.Errorf("Expected to find %v in %v", expectedStr, step.Action)
	}

	expectedStr = fmt.Sprintf("{{ file.download(%v, %v, %v) }}", existingDownloadID,
		existingDownloadSignature[0:6], filepath.Base(existingDownloadPath))
	if !strings.Contains(step.Response, expectedStr) {
		t.Errorf("Expected to find %v in %v", expectedStr, step.Response)
	}

	// Check download values
	step = test.Steps[1].(RFTestStep)
	expectedStr = fmt.Sprintf("{{ file.screenshot(%v, %v) }}", newScreenshotID, newScreenshotSignature[0:6])
	if !strings.Contains(step.Action, expectedStr) {
		t.Errorf("Expected to find %v in %v", expectedStr, step.Action)
	}

	expectedStr = fmt.Sprintf("{{ file.download(%v, %v, %v) }}", newDownloadID, newDownloadSignature[0:6],
		filepath.Base(newDownloadPath))
	if !strings.Contains(step.Response, expectedStr) {
		t.Errorf("Expected to find %v in %v", expectedStr, step.Response)
	}

	// Check for logger warning for nonexistent files
	logs := out.String()
	// -- File download error
	absPath := filepath.Join(testDir, nonexistentFilePath)
	expectedError := fmt.Sprintf("Error for test: %v\n\tNo such file exists: %v", test.RFMLPath, absPath)
	if !strings.Contains(logs, expectedError) {
		t.Errorf("Expecting nonexistent file error for %v. Got: %v", test.RFMLPath, logs)
	}

	// -- Screenshot error
	absPath = filepath.Join(testDir, nonexistentScreenshotPath)
	expectedError = fmt.Sprintf("Error for test: %v\n\tNo such file exists: %v", test.RFMLPath, absPath)
	if !strings.Contains(logs, expectedError) {
		t.Errorf("Expecting nonexistent file error for %v. Got: %v", test.RFMLPath, logs)
	}
}
