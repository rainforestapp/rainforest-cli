package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
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
	if c.Bool("synchronous-upload") {
		rfmlUploadConcurrency = 1
	}
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
		log.Printf("Creating new test: %v", parsedTest.RFMLID)
		err = api.CreateTest(parsedTest)
		if err != nil {
			return err
		}
		log.Printf("Created new test: %v", parsedTest.RFMLID)
		return nil
	}

	// Or update the old one
	log.Printf("Updating existing test: %v", parsedTest.RFMLID)
	err = api.UpdateTest(parsedTest)
	if err != nil {
		return err
	}
	log.Printf("Updated test: %v", parsedTest.RFMLID)
	return nil
}

func uploadRFMLFilesInDirectory(rfmlDirectory string) error {
	// Validate files first
	err := validateRFMLFilesInDirectory(rfmlDirectory)
	if err != nil {
		return err
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

	// This will be used over and over again
	mappings, err := api.GetRFMLIDs()
	if err != nil {
		return err
	}
	rfmlidToID := mappings.MapRFMLIDtoID()
	var parsedTests []*rainforest.RFTest
	var newTests []*rainforest.RFTest

	for _, filePath := range fileList {
		f, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer f.Close()
		rfmlReader := rainforest.NewRFMLReader(f)
		pTest, err := rfmlReader.ReadAll()
		if err != nil {
			return err
		}
		parsedTests = append(parsedTests, pTest)
		// Check if it's a new test or an existing one, because they need different treatment
		// to ensure we first add new ones and have IDs for potential embedds
		if _, ok := rfmlidToID[pTest.RFMLID]; !ok {
			newTests = append(newTests, pTest)
		}
	}
	// chan to gather errors from workers
	errorsChan := make(chan error)

	// prepare empty tests to upload, we will fill the steps later on in case there are some
	// dependiences between them, we want all of the IDs in place
	testsToCreate := make(chan *rainforest.RFTest, len(newTests))
	for _, newTest := range newTests {
		emptyTest := rainforest.RFTest{
			RFMLID:      newTest.RFMLID,
			Description: newTest.Description,
			Title:       newTest.Title,
		}
		err = emptyTest.PrepareToUploadFromRFML(mappings)
		if err != nil {
			return err
		}
		testsToCreate <- &emptyTest
	}
	close(testsToCreate)

	// spawn workers to create the tests
	for i := 0; i < rfmlUploadConcurrency; i++ {
		go testCreationWorker(api, testsToCreate, errorsChan)
	}

	// Read out the workers results
	for i := 0; i < len(newTests); i++ {
		if err := <-errorsChan; err != nil {
			return err
		}
	}

	// Refresh the mappings so we have all of the new tests
	mappings, err = api.GetRFMLIDs()
	if err != nil {
		return err
	}

	// And here we update all of the tests
	testsToUpdate := make(chan *rainforest.RFTest, len(parsedTests))
	for _, testToUpdate := range parsedTests {
		testToUpdate.PrepareToUploadFromRFML(mappings)
		testsToUpdate <- testToUpdate
	}
	close(testsToUpdate)

	// spawn workers to create the tests
	for i := 0; i < rfmlUploadConcurrency; i++ {
		go testUpdateWorker(api, testsToUpdate, errorsChan)
	}

	// Read out the workers results
	for i := 0; i < len(parsedTests); i++ {
		if err := <-errorsChan; err != nil {
			return err
		}
	}

	return nil
}

func testCreationWorker(api *rainforest.Client,
	testsToCreate <-chan *rainforest.RFTest, errorsChan chan<- error) {
	for test := range testsToCreate {
		log.Printf("Creating new test: %v", test.RFMLID)
		error := api.CreateTest(test)
		if error == nil {
			log.Printf("Created new test: %v", test.RFMLID)
		}
		errorsChan <- error
	}
}

func testUpdateWorker(api *rainforest.Client,
	testsToUpdate <-chan *rainforest.RFTest, errorsChan chan<- error) {
	for test := range testsToUpdate {
		log.Printf("Updating existing test: %v", test.RFMLID)
		error := api.UpdateTest(test)
		if error == nil {
			log.Printf("Updated existing test: %v", test.RFMLID)
		}
		errorsChan <- error
	}
}
