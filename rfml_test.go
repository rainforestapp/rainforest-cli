package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rainforestapp/rainforest-cli/rainforest"
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
	mappings     rainforest.TestIDMappings
	testMappings map[int]rainforest.RFTest
}

func (t *testRfmlAPI) GetRFMLIDs() (rainforest.TestIDMappings, error) {
	return t.mappings, nil
}

func (t *testRfmlAPI) GetTest(testID int) (*rainforest.RFTest, error) {
	test, ok := t.testMappings[testID]
	if !ok {
		return nil, errors.New("Test ID not found")
	}
	return &test, nil
}

func TestDownloadRFML(t *testing.T) {
	context := new(fakeContext)
	testAPI := new(testRfmlAPI)
	testDefaultSpecFolder := "testing/" + defaultSpecFolder

	testID := 112233
	rfmlID := "rfml_test_id"
	title := "My Test Title"

	rfTest := rainforest.RFTest{
		TestID: testID,
		RFMLID: rfmlID,
		Title:  title,
	}

	testAPI.mappings = rainforest.TestIDMappings{{ID: testID, RFMLID: rfmlID}}
	testAPI.testMappings = map[int]rainforest.RFTest{
		testID: rfTest,
	}

	context.mappings = map[string]interface{}{
		"test-folder": testDefaultSpecFolder,
	}

	err := downloadRFML(context, testAPI)
	if err != nil {
		t.Fatal(err.Error())
	}

	paddedTestID := fmt.Sprintf("%010d", testID)
	sanitizedTitle := strings.TrimSpace(title)
	expectedFileName := fmt.Sprintf("%v_%v.rfml", paddedTestID, sanitizedTitle)
	expectedRFMLPath := filepath.Join(testDefaultSpecFolder, expectedFileName)

	_, err = os.Stat(expectedRFMLPath)
	if os.IsNotExist(err) {
		t.Errorf("Expected RFML test does not exist: %v", expectedRFMLPath)
		return
	}

	var contents []byte
	contents, err = ioutil.ReadFile(expectedRFMLPath)
	if err != nil {
		t.Error(err.Error())
		return
	}

	rfmlText := string(contents)

	if !strings.Contains(rfmlText, title) {
		t.Errorf("Expected title \"%v\" to appear in RFML test", title)
	}

	if !strings.Contains(rfmlText, rfmlID) {
		t.Errorf("Expected RFML ID \"%v\" to appear in RFML test", rfmlID)
	}

	err = os.RemoveAll("./testing")
	if err != nil {
		t.Fatal(err.Error())
	}
}
