package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gyuho/goraph"
	"github.com/rainforestapp/rainforest-cli/rainforest"
	"github.com/urfave/cli"
)

// parseError is a custom error implementing error interface for reporting RFML parsing errors.
type fileParseError struct {
	filePath   string
	parseError error
}

func (e fileParseError) Error() string {
	return fmt.Sprintf("%v:%v", e.filePath, e.parseError.Error())
}

// validateRFML is a wrapper around two other validation functions
// first one for the single file and the other for whole directory
func validateRFML(c cliContext) error {
	if path := c.Args().First(); path != "" {
		err := validateSingleRFMLFile(path)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		return nil
	}
	err := validateRFMLFilesInDirectory(c.String("test-folder"))
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	return nil
}

// validateSingleRFMLFile validates RFML file syntax by
// trying to parse the file and sending any parse errors to the caller
func validateSingleRFMLFile(filePath string) error {
	if !strings.Contains(filePath, ".rfml") {
		return errors.New("RFML files should have .rfml extension")
	}
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	rfmlReader := rainforest.NewRFMLReader(f)
	_, err = rfmlReader.ReadAll()
	if err != nil {
		return fileParseError{filePath, err}
	}
	log.Printf("%v's syntax is valid", filePath)
	return nil
}

// validateRFMLFilesInDirectory validates RFML file syntax, embedded rfml ids,
// checks for circular dependiences and all other cool things in the specified directory
func validateRFMLFilesInDirectory(rfmlDirectory string) error {
	// first just make sure we are dealing with directory
	dirStat, err := os.Stat(rfmlDirectory)
	if err != nil {
		return err
	}
	if !dirStat.IsDir() {
		return fmt.Errorf("%v should be a directory", rfmlDirectory)
	}

	// walk through the specifed directory (also subdirs) and pick the .rfml files
	var fileList []string
	err = filepath.Walk(rfmlDirectory, func(path string, f os.FileInfo, err error) error {
		if strings.Contains(path, ".rfml") {
			fileList = append(fileList, path)
		}
		return nil
	})
	if err != nil {
		return err
	}

	// parse all of them files
	type parsedTest struct {
		filePath string
		content  *rainforest.RFTest
	}
	var validationErrors []error
	var parsedTests []parsedTest
	dependencyGraph := goraph.NewGraph()
	for _, filePath := range fileList {
		f, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer f.Close()
		rfmlReader := rainforest.NewRFMLReader(f)
		pTest, err := rfmlReader.ReadAll()
		if err != nil {
			validationErrors = append(validationErrors, fileParseError{filePath, err})
		} else {
			parsedTests = append(parsedTests, parsedTest{filePath, pTest})
		}
	}

	// check for rfml_id uniqueness
	rfmlIDToTest := make(map[string]parsedTest)
	for _, pTest := range parsedTests {
		if conflictingTest, ok := rfmlIDToTest[pTest.content.RFMLID]; ok {
			err = fmt.Errorf(" duplicate RFML id %v, also found in: %v", pTest.content.RFMLID, conflictingTest.filePath)
			validationErrors = append(validationErrors, fileParseError{pTest.filePath, err})
		} else {
			rfmlIDToTest[pTest.content.RFMLID] = pTest
			dependencyGraph.AddNode(goraph.NewNode(pTest.content.RFMLID))
		}
	}

	// check for embedded tests id validity
	// start with pulling the external test ids to validate against them as well
	if api.ClientToken != "" {
		externalTests, err := api.GetRFMLIDs()
		if err != nil {
			return err
		}
		for _, externalTest := range externalTests {
			if _, ok := rfmlIDToTest[externalTest.RFMLID]; !ok {
				rfmlIDToTest[externalTest.RFMLID] = parsedTest{"external", &rainforest.RFTest{}}
				dependencyGraph.AddNode(goraph.NewNode(externalTest.RFMLID))
			}
		}
	}
	// go through all the tests
	for _, pTest := range parsedTests {
		// and steps...
		for stepNum, step := range pTest.content.Steps {
			// then check if it's embeddedTest
			if embeddedTest, ok := step.(rainforest.RFEmbeddedTest); ok {
				// if so, check if its rfml id exists
				if _, ok := rfmlIDToTest[embeddedTest.RFMLID]; !ok {
					if api.ClientToken != "" {
						err = fmt.Errorf("step %v - embeddedTest RFML id %v not found", stepNum+1, embeddedTest.RFMLID)
					} else {
						err = fmt.Errorf("step %v - embeddedTest RFML id %v not found. Specify token_id to check against external tests", stepNum+1, embeddedTest.RFMLID)
					}
					validationErrors = append(validationErrors, fileParseError{pTest.filePath, err})
				} else {
					pNode := dependencyGraph.GetNode(goraph.StringID(pTest.content.RFMLID))
					eNode := dependencyGraph.GetNode(goraph.StringID(embeddedTest.RFMLID))
					dependencyGraph.AddEdge(pNode.ID(), eNode.ID(), 1)
				}
			}
		}
	}

	// validate circular dependiences probably using Tarjan's strongly connected components
	stronglyConnected := goraph.Tarjan(dependencyGraph)
	for _, circularTests := range stronglyConnected {
		if len(circularTests) > 1 {
			err = fmt.Errorf("Found circular dependiences between: %v", circularTests)
			validationErrors = append(validationErrors, err)
		}
	}

	if len(validationErrors) > 0 {
		for _, err := range validationErrors {
			log.Print(err.Error())
		}
		return errors.New("Validation failed")
	}

	log.Print("All files are valid!")
	return nil
}

func deleteRFML(c cliContext) error {
	filePath := c.Args().First()
	if !strings.Contains(filePath, ".rfml") {
		return cli.NewExitError("RFML files should have .rfml extension", 1)
	}
	f, err := os.Open(filePath)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	rfmlReader := rainforest.NewRFMLReader(f)
	parsedRFML, err := rfmlReader.ReadAll()
	if parsedRFML.RFMLID == "" {
		return cli.NewExitError("RFML file doesn't have RFML ID", 1)
	}

	// Close the file now so we can delete it
	f.Close()

	// Delete remote first
	err = api.DeleteTestByRFMLID(parsedRFML.RFMLID)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	// Then delete local file
	err = os.Remove(filePath)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	return nil
}

// uploadRFML is a wrapper around test creating/updating functions
func uploadRFML(c cliContext) error {
	if path := c.Args().First(); path != "" {
		err := uploadSingleRFMLFile(path)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		return nil
	}
	err := uploadRFMLFilesInDirectory(c.String("test-folder"))
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	return nil
}

// uploadSingleRFMLFile uploads RFML file syntax by
// trying to parse the file and sending any parse errors to the caller
func uploadSingleRFMLFile(filePath string) error {
	// Validate first before uploading
	err := validateSingleRFMLFile(filePath)
	if err != nil {
		return err
	}

	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	rfmlReader := rainforest.NewRFMLReader(f)
	parsedTest, err := rfmlReader.ReadAll()
	if err != nil {
		return fileParseError{filePath, err}
	}

	// Check if the test already exists in RF so we can decide between updating and creating new one
	mappings, err := api.GetRFMLIDs()
	if err != nil {
		return err
	}
	err = parsedTest.PrepareToUploadFromRFML(mappings)
	if err != nil {
		return err
	}
	if parsedTest.TestID == 0 {
		// Create new test
		return api.CreateTest(parsedTest)
	}

	// Or update the old one
	return api.UpdateTest(parsedTest)
}

func uploadRFMLFilesInDirectory(rfmlDirectory string) error {
	return nil
}

func downloadRFML(c cliContext) error {
	var testIDs []int
	var err error

	if len(c.Args()) > 0 {
		var testID int

		for _, arg := range c.Args() {
			testID, err = strconv.Atoi(arg)
			if err != nil {
				return cli.NewExitError(err.Error(), 1)
			}

			testIDs = append(testIDs, testID)
		}
	} else {
		var rfmlIDMapping rainforest.TestIDMappings
		rfmlIDMapping, err = api.GetRFMLIDs()
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		for _, testIDMap := range rfmlIDMapping {
			testID := testIDMap.ID
			testIDs = append(testIDs, testID)
		}
	}

	errorsChan := make(chan error)
	testIDChan := make(chan int, len(testIDs))

	var mappings rainforest.TestIDMappings
	mappings, err = api.GetRFMLIDs()
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	for _, testID := range testIDs {
		testIDChan <- testID
	}
	close(testIDChan)

	for i := 0; i < rfmlDownloadConcurrency; i++ {
		go downloadRFTestWorker(testIDChan, errorsChan, mappings)
	}

	for i := 0; i < len(testIDs); i++ {
		if err := <-errorsChan; err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
	}

	return nil
}

func downloadRFTestWorker(testIDChan chan int, errorsChan chan error, mappings rainforest.TestIDMappings) {
	for testID := range testIDChan {
		test, err := api.GetTest(testID)
		if err != nil {
			errorsChan <- err
			return
		}

		test.UnmapBrowsers()
		err = test.UnmarshalElements(mappings)
		if err != nil {
			errorsChan <- err
			return
		}

		// TODO: Merge with RFML Writer branch to write RFML tests

		// if err is nil, nothing will happen. If err is not nil, the program will end
		errorsChan <- err
	}
}
