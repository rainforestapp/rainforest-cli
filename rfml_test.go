package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/rainforestapp/rainforest-cli/rainforest"
	"github.com/urfave/cli"
)

func TestNewRFMLTest(t *testing.T) {
	context := new(fakeContext)
	testDefaultSpecFolder := "testing/" + defaultSpecFolder

	/*
	   Declare reusable variables
	*/
	var contents []byte
	var rfmlText string
	var err error
	var expectedRFMLPath string
	var file *os.File

	/*
	   Helper functions
	*/
	testExpectation := func(filePath string, title string) {
		_, err = os.Stat(filePath)
		if os.IsNotExist(err) {
			t.Errorf("Expected RFML test does not exist: %v", filePath)
			return
		}

		contents, err = ioutil.ReadFile(filePath)
		if err != nil {
			t.Error(err.Error())
			return
		}

		rfmlText = string(contents)

		if !strings.Contains(rfmlText, title) {
			t.Errorf("Expected title \"%v\" to appear in RFML test", title)
		}
	}

	/*
	   No flags or args and spec folder doesn't exist
	*/
	context.mappings = map[string]interface{}{
		"test-folder": testDefaultSpecFolder,
	}

	err = newRFMLTest(context)
	if err != nil {
		t.Error(err.Error())
	}

	expectedRFMLPath = filepath.Join(testDefaultSpecFolder, "Unnamed Test.rfml")
	testExpectation(expectedRFMLPath, "Unnamed Test")
	err = os.RemoveAll("./testing")
	if err != nil {
		t.Fatal(err.Error())
	}

	/*
	   No flags or args and spec folder does exist
	*/
	err = os.MkdirAll(testDefaultSpecFolder, os.ModePerm)
	if err != nil {
		t.Fatal(err.Error())
	}

	err = newRFMLTest(context)
	if err != nil {
		t.Error(err.Error())
	}

	testExpectation(expectedRFMLPath, "Unnamed Test")
	err = os.RemoveAll("./testing")
	if err != nil {
		t.Fatal(err.Error())
	}

	/*
	   Test folder given
	*/
	specFolder := "./my_specs"
	err = os.MkdirAll(specFolder, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	context.mappings = map[string]interface{}{
		"test-folder": specFolder,
	}

	err = newRFMLTest(context)
	if err != nil {
		err = os.RemoveAll(specFolder)
		if err != nil {
			t.Fatal(err.Error())
		}
		t.Fatal(err)
	}

	expectedRFMLPath = filepath.Join(specFolder, "Unnamed Test.rfml")
	testExpectation(expectedRFMLPath, "Unnamed Test")
	err = os.RemoveAll(specFolder)
	if err != nil {
		t.Fatal(err.Error())
	}

	/*
	   Filename argument given
	*/
	context.mappings = map[string]interface{}{
		"test-folder": testDefaultSpecFolder,
	}

	context.args = []string{"my_file_name.rfml"}

	err = newRFMLTest(context)
	if err != nil {
		t.Fatal(err.Error())
	}

	expectedRFMLPath = filepath.Join(testDefaultSpecFolder, "my_file_name.rfml")
	testExpectation(expectedRFMLPath, "my_file_name")
	err = os.RemoveAll("./testing")
	if err != nil {
		t.Fatal(err.Error())
	}

	/*
	   Title argument given
	*/
	context.args = []string{"my_test_title"}

	err = newRFMLTest(context)
	if err != nil {
		t.Fatal(err.Error())
	}

	expectedRFMLPath = filepath.Join(testDefaultSpecFolder, "my_test_title.rfml")
	testExpectation(expectedRFMLPath, "my_test_title")
	err = os.RemoveAll("./testing")
	if err != nil {
		t.Fatal(err.Error())
	}

	/*
	   Test folder flag is actually a file
	*/
	dummyFolder := "./dummy"
	err = os.MkdirAll(dummyFolder, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	dummyFilePath := filepath.Join(dummyFolder, "dummy_file")
	file, err = os.Create(dummyFilePath)
	if err != nil {
		err = os.RemoveAll(dummyFolder)
		if err != nil {
			t.Fatal(err.Error())
		}
		t.Fatal(err)
	}

	file.Close()

	context.args = []string{}
	context.mappings = map[string]interface{}{
		"test-folder": dummyFilePath,
	}

	err = newRFMLTest(context)
	if err == nil {
		t.Error("Expecting an error, got nil")
	}

	os.RemoveAll(dummyFolder)

	/*
	   RFML file already exists
	*/
	context.mappings = map[string]interface{}{
		"test-folder": testDefaultSpecFolder,
	}

	err = os.MkdirAll(testDefaultSpecFolder, os.ModePerm)
	if err != nil {
		t.Fatal(err.Error())
	}

	existingRFMLPath := filepath.Join(testDefaultSpecFolder, "Unnamed Test.rfml")
	file, err = os.Create(existingRFMLPath)
	file.Close()
	if err != nil {
		t.Fatal(err.Error())
	}

	err = newRFMLTest(context)
	if err != nil {
		err = os.RemoveAll("./testing")
		if err != nil {
			t.Fatal(err.Error())
		}
		t.Fatal(err.Error())
	}

	expectedRFMLPath = filepath.Join(testDefaultSpecFolder, "Unnamed Test (1).rfml")
	_, err = os.Stat(expectedRFMLPath)
	if err != nil {
		t.Error(err.Error())
	}

	err = os.RemoveAll("./testing")
	if err != nil {
		t.Fatal(err.Error())
	}
}

type testRfmlAPI struct {
	mappings         rainforest.TestIDMappings
	tests            []rainforest.RFTest
	handleUpdateTest func(*rainforest.RFTest)
}

func (t *testRfmlAPI) GetRFMLIDs() (*rainforest.TestIDMappings, error) {
	return &t.mappings, nil
}

func (t *testRfmlAPI) GetTest(testID int) (*rainforest.RFTest, error) {
	for _, test := range t.tests {
		if test.TestID == testID {
			return &test, nil
		}
	}
	return nil, errors.New("Test ID not found")
}

func (t *testRfmlAPI) GetTests(*rainforest.RFTestFilters) ([]rainforest.RFTest, error) {
	return t.tests, nil
}

func (t *testRfmlAPI) ClientToken() string {
	return "abc123"
}

func (t *testRfmlAPI) CreateTest(_ *rainforest.RFTest) error {
	// implement when needed
	return errStub
}

func (t *testRfmlAPI) UpdateTest(test *rainforest.RFTest) error {
	t.handleUpdateTest(test)
	return nil
}

func (t *testRfmlAPI) ParseEmbeddedFiles(_ *rainforest.RFTest) error {
	// implement when needed
	return errStub
}

func createTestFolder(testFolderPath string) error {
	absTestFolderPath, err := filepath.Abs(testFolderPath)
	if err != nil {
		return err
	}

	dirStat, err := os.Stat(absTestFolderPath)
	if os.IsNotExist(err) {
		os.MkdirAll(absTestFolderPath, os.ModePerm)
	} else if err != nil {
		return err
	} else if !dirStat.IsDir() {
		return fmt.Errorf("Test folder path is not a directory: %v", absTestFolderPath)
	}

	return nil
}

func cleanUpTestFolder(testFolderPath string) error {
	_, err := os.Stat(testFolderPath)

	if err != nil && os.IsNotExist(err) {
		return err
	}

	err = os.RemoveAll(testFolderPath)
	if err != nil {
		return err
	}

	return nil
}

func TestUploadRFML(t *testing.T) {
	context := new(fakeContext)
	testAPI := new(testRfmlAPI)
	testDefaultSpecFolder := "testing/" + defaultSpecFolder

	defer func() {
		err := cleanUpTestFolder("testing")
		if err != nil {
			t.Fatal(err.Error())
		}
	}()

	context.mappings = map[string]interface{}{
		"test-folder": testDefaultSpecFolder,
	}

	testID := 666
	rfmlID := "unique_rfml_id"
	title := "a very descriptive title"
	var featureID rainforest.FeatureIDInt = 777

	err := createTestFolder(testDefaultSpecFolder)
	if err != nil {
		t.Fatal(err.Error())
	}

	testPath := filepath.Join(testDefaultSpecFolder, "valid_test.rfml")

	testAPI.mappings = rainforest.TestIDMappings{
		Pairs: []rainforest.TestIDMap{
			{ID: testID, RFMLID: rfmlID},
		},
	}

	// basic test
	testAPI.handleUpdateTest = func(rfTest *rainforest.RFTest) {
		testCases := []struct {
			fieldName string
			expected  interface{}
			got       interface{}
		}{
			{"test ID", testID, rfTest.TestID},
			{"RFML ID", rfmlID, rfTest.RFMLID},
			{"title", title, rfTest.Title},
			{"feature ID", featureID, rfTest.FeatureID},
			{"disabled state", "enabled", rfTest.State},
		}

		for _, testCase := range testCases {
			if testCase.got != testCase.expected {
				t.Errorf("Incorrect value for %v. Expected %v, Got %v", testCase.fieldName, testCase.expected, testCase.got)
			}
		}
	}

	testContents := fmt.Sprintf(`#! %v
# title: %v
# feature_id: %v
`, rfmlID, title, featureID)

	err = ioutil.WriteFile(testPath, []byte(testContents), os.ModePerm)
	if err != nil {
		t.Fatal(err.Error())
	}

	err = uploadRFML(context, testAPI)
	if err != nil {
		t.Fatal(err.Error())
	}

	// state is specified
	testAPI.handleUpdateTest = func(rfTest *rainforest.RFTest) {
		if rfTest.State != "disabled" {
			t.Errorf("Incorrect value for state. Expected \"disabled\", Got %v", rfTest.State)
		}
	}

	testContents = fmt.Sprintf(`#! %v
# title: %v
# state: disabled
`, rfmlID, title)

	err = ioutil.WriteFile(testPath, []byte(testContents), os.ModePerm)
	if err != nil {
		t.Fatal(err.Error())
	}

	err = uploadRFML(context, testAPI)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestDeleteRFML(t *testing.T) {
	// Test error in parsing file
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err.Error())
	}
	defer os.RemoveAll(dir)

	rfmlFilePath := filepath.Join(dir, "testing.rfml")
	fileContents := `#! testing
# title: hello
# site_id: a_string
`
	err = ioutil.WriteFile(rfmlFilePath, []byte(fileContents), 0666)
	if err != nil {
		t.Fatal(err.Error())
	}

	dummyMappings := map[string]interface{}{}
	args := cli.Args{rfmlFilePath}
	ctx := newFakeContext(dummyMappings, args)

	err = deleteRFML(ctx)
	if err == nil {
		t.Fatal("Expected parse error but received no error.")
	}

	errMsg := err.Error()
	if !strings.Contains(errMsg, rfmlFilePath) {
		t.Errorf("Expected error to contain file path \"%v\". Got:\n%v", rfmlFilePath, errMsg)
	}
}

func TestDownloadRFML(t *testing.T) {
	context := new(fakeContext)
	testAPI := new(testRfmlAPI)
	testDefaultSpecFolder := "testing/" + defaultSpecFolder

	defer func() {
		err := cleanUpTestFolder("testing")
		if err != nil {
			t.Fatal(err.Error())
		}
	}()

	testID := 112233
	rfmlID := "rfml_test_id"
	title := "My Test Title"
	featureID := 665544

	rfTest := rainforest.RFTest{
		TestID:    testID,
		RFMLID:    rfmlID,
		Title:     title,
		FeatureID: rainforest.FeatureIDInt(featureID),
		State:     "enabled",
	}

	testAPI.mappings = rainforest.TestIDMappings{
		Pairs: []rainforest.TestIDMap{
			{ID: testID, RFMLID: rfmlID},
		},
	}
	testAPI.tests = []rainforest.RFTest{rfTest}

	paddedTestID := fmt.Sprintf("%010d", testID)
	sanitizedTitle := "my_test_title"
	expectedFileName := fmt.Sprintf("%v_%v.rfml", paddedTestID, sanitizedTitle)
	expectedRFMLPath := filepath.Join(testDefaultSpecFolder, expectedFileName)

	context.mappings = map[string]interface{}{
		"test-folder": testDefaultSpecFolder,
	}

	// basic test
	err := downloadRFML(context, testAPI)
	if err != nil {
		t.Fatal(err.Error())
	}

	fileInfo, err := os.Stat(expectedRFMLPath)
	if os.IsNotExist(err) {
		t.Fatalf("Expected RFML test does not exist: %v", expectedRFMLPath)
	}

	if fileInfo.Name() != expectedFileName {
		t.Errorf("Expected RFML file path %v, got %v", expectedRFMLPath, fileInfo.Name())
	} else if err != nil {
		t.Fatalf(err.Error())
	}

	var contents []byte
	contents, err = ioutil.ReadFile(expectedRFMLPath)
	if err != nil {
		t.Fatalf(err.Error())
	}

	rfmlText := string(contents)

	if !strings.Contains(rfmlText, title) {
		t.Errorf("Expected title \"%v\" to appear in RFML test", title)
	}

	if !strings.Contains(rfmlText, rfmlID) {
		t.Errorf("Expected RFML ID \"%v\" to appear in RFML test", rfmlID)
	}

	if !strings.Contains(rfmlText, strconv.Itoa(featureID)) {
		t.Errorf("Expected Feature ID \"%v\" to appear in RFML test", featureID)
	}

	if strings.Contains(rfmlText, "state") {
		t.Errorf("Did not expect state field in RFML test. Got %v", rfmlText)
	}

	// Test is disabled
	rfTest.State = "disabled"
	testAPI.tests = []rainforest.RFTest{rfTest}

	err = downloadRFML(context, testAPI)
	if err != nil {
		t.Fatal(err.Error())
	}

	fileInfo, err = os.Stat(expectedRFMLPath)
	if os.IsNotExist(err) {
		t.Fatalf("Expected RFML test does not exist: %v", expectedRFMLPath)
	}

	contents, err = ioutil.ReadFile(expectedRFMLPath)
	if err != nil {
		t.Fatalf(err.Error())
	}
	rfmlText = string(contents)

	if !strings.Contains(rfmlText, "# state: disabled") {
		t.Errorf("Expected RFML test state to read disabled. Output: %v", rfmlText)
	}
}

func TestSanitizeTestTitle(t *testing.T) {
	// Test that it replaces non-alphanumeric characters with underscores
	illegalTitle := `Foo\123|*&bar `
	sanitizedTitle := sanitizeTestTitle(illegalTitle)
	expectedSanitizedTitle := "foo_123_bar"

	if sanitizedTitle != expectedSanitizedTitle {
		t.Errorf("Expected sanitized title to be %v, got %v", expectedSanitizedTitle, sanitizedTitle)
	}

	// Test that it truncates strings with more than 30 characters
	longTitle := strings.Repeat("abc", 11) // 33 characters
	sanitizedTitle = sanitizeTestTitle(longTitle)
	expectedSanitizedTitle = strings.Repeat("abc", 10) // 30 characters

	if sanitizedTitle != expectedSanitizedTitle {
		t.Errorf("Expected sanitized title to be %v, got %v", expectedSanitizedTitle, sanitizedTitle)
	}
}

func TestValidateEmbedded(t *testing.T) {
	t1 := rainforest.RFTest{
		TestID:   1,
		RFMLID:   "1",
		RFMLPath: "./1.rfml",
	}
	t2 := rainforest.RFTest{
		TestID:   2,
		RFMLID:   "2",
		RFMLPath: "./2.rfml",
		Steps: []interface{}{
			rainforest.RFEmbeddedTest{
				RFMLID:   "1",
				Redirect: true,
			},
		},
	}

	testAPI := new(testRfmlAPI)
	testAPI.mappings = rainforest.TestIDMappings{
		Pairs: []rainforest.TestIDMap{
			{ID: t1.TestID, RFMLID: t1.RFMLID},
			{ID: t2.TestID, RFMLID: t2.RFMLID},
		},
	}

	tests := []*rainforest.RFTest{&t2}
	var err error

	// With API access, the validation should be fine
	err = validateRFMLFiles(tests, false, testAPI)
	if err != nil {
		t.Error("Non-local validation failed:", err)
	}

	// With local-only and all embedded tests available, the validation should
	// be fine
	allTests := append(tests, &t1)
	err = validateRFMLFiles(allTests, true, testAPI)
	if err != nil {
		t.Error("Local-only validation with all tests failed:", err)
	}

	// With local-only and embedded tests not available locally, the validation
	// should fail
	err = validateRFMLFiles(tests, true, testAPI)
	if err == nil {
		t.Error("Local-only validation should have failed but didn't")
	} else if err != errValidation {
		t.Error("Unexpected error for local-only validation:", err)
	}
}

func TestReadRFMLFiles(t *testing.T) {
	dir := setupTestRFMLDir()
	defer os.RemoveAll(dir)

	var testCases = []struct {
		files     []string
		wantFiles []string
		wantError bool
	}{
		{
			files:     []string{"a/a1.rfml"},
			wantFiles: []string{"a/a1.rfml"},
		},
		{
			files:     []string{"a"},
			wantFiles: []string{"a/a1.rfml", "a/a2.rfml", "a/a3.rfml"},
		},
		{
			files:     []string{"a", "a/a1.rfml", "b/b1.rfml"},
			wantFiles: []string{"a/a1.rfml", "a/a2.rfml", "a/a3.rfml", "b/b1.rfml"},
		},
		{
			files: []string{""},
			wantFiles: []string{
				"a/a1.rfml",
				"a/a2.rfml",
				"a/a3.rfml",
				"b/b1.rfml",
				"b/a/b2.rfml",
				"b/b/b3.rfml",
				"b/b/b4.rfml",
				"b/b/b5.rfml",
				"standalone.rfml",
			},
		},
		{
			files:     []string{"c"},
			wantError: true,
		},
		{
			// We want to just ignore bogus files, rather than error, so that
			// shell expansions like foo/ab* still work
			files:     []string{"a/a1.rfml", "a/bogus.rf"},
			wantFiles: []string{"a/a1.rfml"},
		},
	}

	for _, tc := range testCases {
		fullFiles := []string{}
		for _, f := range tc.files {
			fullFiles = append(fullFiles, filepath.Join(dir, f))
		}

		pTests, err := readRFMLFiles(fullFiles)
		if err != nil && !tc.wantError {
			t.Error(err)
			continue
		} else if err == nil && tc.wantError {
			t.Errorf("Expected error when reading %v", tc.files)
			continue
		}

		gotFiles := []string{}
		for _, p := range pTests {
			gotFiles = append(gotFiles, p.RFMLPath)
		}
		sort.Strings(gotFiles)

		wantFiles := []string{}
		for _, f := range tc.wantFiles {
			wantFiles = append(wantFiles, filepath.Join(dir, f))
		}
		sort.Strings(wantFiles)

		if !reflect.DeepEqual(wantFiles, gotFiles) {
			t.Errorf("Unexpected files returned (want: %v, got: %v)", wantFiles, gotFiles)
		}
	}
}

func TestReadRFMLFile(t *testing.T) {
	// Test error in parsing file
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err.Error())
	}
	defer os.RemoveAll(dir)

	rfmlFilePath := filepath.Join(dir, "testing.rfml")
	fileContents := `#! testing
# title: hello
# site_id: a_string
`
	err = ioutil.WriteFile(rfmlFilePath, []byte(fileContents), 0666)
	if err != nil {
		t.Fatal(err.Error())
	}

	_, err = readRFMLFile(rfmlFilePath)
	if err == nil {
		t.Fatal("Expected parse error but received no error.")
	}

	errMsg := err.Error()
	if !strings.Contains(errMsg, rfmlFilePath) {
		t.Errorf("Expected error to contain file path \"%v\". Got:\n%v", rfmlFilePath, errMsg)
	}
}

func setupTestRFMLDir() string {
	var rfmlTests = []struct {
		path    string
		content *rainforest.RFTest
	}{
		{
			path: "a/a1.rfml",
			content: &rainforest.RFTest{
				RFMLID:  "a1",
				Tags:    []string{"foo", "baz"},
				Execute: true,
				Steps: []interface{}{
					rainforest.RFEmbeddedTest{
						RFMLID: "b4",
					},
				},
			},
		},
		{
			path: "a/a2.rfml",
			content: &rainforest.RFTest{
				RFMLID:  "a2",
				Execute: true,
			},
		},
		{
			path: "a/a3.rfml",
			content: &rainforest.RFTest{
				RFMLID:  "a3",
				Tags:    []string{"bar"},
				Execute: true,
			},
		},
		{
			path: "b/b1.rfml",
			content: &rainforest.RFTest{
				RFMLID:  "b1",
				Execute: true,
			},
		},
		{
			path: "b/a/b2.rfml",
			content: &rainforest.RFTest{
				RFMLID:  "b2",
				Execute: true,
			},
		},
		{
			path: "b/b/b3.rfml",
			content: &rainforest.RFTest{
				RFMLID:  "b3",
				Tags:    []string{"foo"},
				Execute: false,
			},
		},
		{
			path: "b/b/b4.rfml",
			content: &rainforest.RFTest{
				RFMLID:  "b4",
				Tags:    []string{},
				Execute: true,
				Steps: []interface{}{
					rainforest.RFEmbeddedTest{
						RFMLID: "b5",
					},
				},
			},
		},
		{
			path: "b/b/b5.rfml",
			content: &rainforest.RFTest{
				RFMLID:  "b5",
				Execute: true,
			},
		},
		{
			path: "standalone.rfml",
			content: &rainforest.RFTest{
				RFMLID:  "standalone",
				Execute: true,
			},
		},
		{
			path: "a/bogus.rf",
			content: &rainforest.RFTest{
				RFMLID:  "bogus",
				Execute: true,
			},
		},
	}

	dir, err := ioutil.TempDir("", "")
	if err != nil {
		log.Fatal(err)
	}

	for _, subdir := range []string{"a", "b/a", "b/b"} {
		err := os.MkdirAll(filepath.Join(dir, subdir), os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}

	for i, test := range rfmlTests {
		test.content.TestID = i
		test.content.Title = test.content.RFMLID
		p := filepath.Join(dir, test.path)
		test.content.RFMLPath = p
		f, err := os.Create(p)
		if err != nil {
			log.Fatal(err)
		}
		w := rainforest.NewRFMLWriter(f)
		err = w.WriteRFMLTest(test.content)
		if err != nil {
			log.Fatal(err)
		}
	}

	return dir
}
