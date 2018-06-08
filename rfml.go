package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/gyuho/goraph"
	"github.com/rainforestapp/rainforest-cli/rainforest"
	"github.com/satori/go.uuid"
	"github.com/urfave/cli"
)

type rfmlAPI interface {
	GetRFMLIDs() (*rainforest.TestIDMappings, error)
	GetTests(*rainforest.RFTestFilters) ([]rainforest.RFTest, error)
	GetTest(int) (*rainforest.RFTest, error)
	CreateTest(*rainforest.RFTest) error
	UpdateTest(*rainforest.RFTest) error
	ParseEmbeddedFiles(*rainforest.RFTest) error
	ClientToken() string
}

// parseError is a custom error implementing error interface for reporting RFML parsing errors.
type fileParseError struct {
	filePath   string
	parseError error
}

func (e fileParseError) Error() string {
	return fmt.Sprintf("%v: %v", e.filePath, e.parseError.Error())
}

// validateRFML is a wrapper around two other validation functions
// first one for the single file and the other for whole directory
func validateRFML(c cliContext, api rfmlAPI) error {
	if path := c.Args().First(); path != "" {
		err := validateSingleRFMLFile(path)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		return nil
	}
	tests, err := readRFMLFiles([]string{c.String("test-folder")})
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	err = validateRFMLFiles(tests, false, api)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	return nil
}

// readRFMLFiles takes in a list of files and/or directories and a list of tags
// and returns a list of the parsed tests, or an error if it is encountered. To
// allow all tags, pass in nil for tags.
func readRFMLFiles(files []string) ([]*rainforest.RFTest, error) {
	fileList := []string{}
	for _, file := range files {
		stat, err := os.Stat(file)
		if err != nil {
			return nil, err
		}
		if !stat.IsDir() {
			if strings.HasSuffix(file, ".rfml") {
				fileList = append(fileList, file)
				continue
			} else {
				log.Printf("%s is not a valid RFML file", file)
				continue
			}
		}

		// We have a directory, walk through and find RFML files
		err = filepath.Walk(file, func(path string, f os.FileInfo, err error) error {
			if strings.HasSuffix(path, ".rfml") {
				fileList = append(fileList, path)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	tests := []*rainforest.RFTest{}
	seenPaths := map[string]bool{}
	for _, filePath := range fileList {
		// No dups!
		if seenPaths[filePath] {
			continue
		}
		seenPaths[filePath] = true
		test, err := readRFMLFile(filePath)
		if err != nil {
			return nil, err
		}
		tests = append(tests, test)
	}
	return tests, nil
}

// anyMember is one of those things that would probably be in the stdlib if
// there were generics. I hate golang sometimes. In any case, it returns true if
// any of needles are in haystack. It's O(n*m), so only put small stuff in
// there!
func anyMember(haystack []string, needles []string) bool {
	for _, n := range needles {
		for _, h := range haystack {
			if h == n {
				return true
			}
		}
	}

	return false
}

func readRFMLFile(filePath string) (*rainforest.RFTest, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rfmlReader := rainforest.NewRFMLReader(f)
	var pTest *rainforest.RFTest
	pTest, err = rfmlReader.ReadAll()
	if err != nil {
		return nil, fileParseError{filePath, err}
	}

	pTest.RFMLPath = filePath
	return pTest, err
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

var errValidation = errors.New("Validation failed")

// validateRFMLFiles validates RFML file syntax, embedded rfml ids, checks for
// circular dependiences and all other cool things in the specified directory
func validateRFMLFiles(parsedTests []*rainforest.RFTest, localOnly bool, api rfmlAPI) error {
	// parse all of them files
	var validationErrors []error
	var err error
	dependencyGraph := goraph.NewGraph()

	// check for rfml_id uniqueness
	rfmlIDToTest := make(map[string]*rainforest.RFTest)
	for _, pTest := range parsedTests {
		if conflictingTest, ok := rfmlIDToTest[pTest.RFMLID]; ok {
			err = fmt.Errorf(" duplicate RFML id %v, also found in: %v", pTest.RFMLID, conflictingTest.RFMLPath)
			validationErrors = append(validationErrors, fileParseError{pTest.RFMLPath, err})
		} else {
			rfmlIDToTest[pTest.RFMLID] = pTest
			dependencyGraph.AddNode(goraph.NewNode(pTest.RFMLID))
		}
	}

	// check for embedded tests id validity
	// start with pulling the external test ids to validate against them as well
	if !localOnly && api.ClientToken() != "" {
		var testIDMappings *rainforest.TestIDMappings
		testIDMappings, err = api.GetRFMLIDs()
		if err != nil {
			return err
		}
		for _, testIDs := range testIDMappings.Pairs {
			if _, ok := rfmlIDToTest[testIDs.RFMLID]; !ok {
				rfmlIDToTest[testIDs.RFMLID] = &rainforest.RFTest{}
				dependencyGraph.AddNode(goraph.NewNode(testIDs.RFMLID))
			}
		}
	}
	// go through all the tests
	for _, pTest := range parsedTests {
		// and steps...
		for stepNum, step := range pTest.Steps {
			// then check if it's embeddedTest
			if embeddedTest, ok := step.(rainforest.RFEmbeddedTest); ok {
				// if so, check if its rfml id exists
				if _, ok := rfmlIDToTest[embeddedTest.RFMLID]; !ok {
					if localOnly || api.ClientToken() != "" {
						err = fmt.Errorf("step %v - embeddedTest RFML id %v not found", stepNum+1, embeddedTest.RFMLID)
					} else {
						err = fmt.Errorf("step %v - embeddedTest RFML id %v not found. Specify token_id to check against external tests", stepNum+1, embeddedTest.RFMLID)
					}
					validationErrors = append(validationErrors, fileParseError{pTest.RFMLPath, err})
				} else {
					pNode := dependencyGraph.GetNode(goraph.StringID(pTest.RFMLID))
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
		return errValidation
	}

	log.Print("All files are valid!")
	return nil
}

func newRFMLTest(c cliContext) error {
	testDirectory := c.String("test-folder")

	absTestDirectory, err := prepareTestDirectory(testDirectory)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	fileName := c.Args().First()
	title := fileName

	if fileName == "" {
		fileName = "Unnamed Test.rfml"
		title = "Unnamed Test"
	} else if strings.HasSuffix(fileName, ".rfml") {
		title = strings.TrimSuffix(title, ".rfml")
	} else {
		fileName = fileName + ".rfml"
	}

	filePath := filepath.Join(absTestDirectory, fileName)

	// Make sure that the file is unique
	basePath := strings.TrimSuffix(filePath, ".rfml")
	fileIdentifier := 0
	var identStr string
	for {
		if fileIdentifier == 0 {
			identStr = ""
		} else {
			identStr = fmt.Sprintf(" (%v)", strconv.Itoa(fileIdentifier))
		}

		testPath := basePath + identStr + ".rfml"

		_, err = os.Stat(testPath)
		if !os.IsNotExist(err) {
			fileIdentifier = fileIdentifier + 1
		} else {
			filePath = testPath
			break
		}
	}

	test := rainforest.RFTest{
		RFMLID:   uuid.NewV4().String(),
		Title:    title,
		StartURI: "/",
		Execute:  true,
		Steps: []interface{}{
			rainforest.RFTestStep{
				Action:   "This is a step action.",
				Response: "This is a step question?",
				Redirect: true,
			},
			rainforest.RFTestStep{
				Action:   "This is another step action.",
				Response: "This is another step question?",
				Redirect: true,
			},
		},
	}

	f, err := os.Create(filePath)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	writer := rainforest.NewRFMLWriter(f)
	err = writer.WriteRFMLTest(&test)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

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
	if err != nil {
		errMsg := fmt.Sprintf("Error removing test at '%v': %v", filePath, err.Error())
		return cli.NewExitError(errMsg, 1)
	}

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
func uploadRFML(c cliContext, api rfmlAPI) error {
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
	tests, err := readRFMLFiles([]string{c.String("test-folder")})
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	err = uploadRFMLFiles(tests, false, api)
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
	parsedTest.RFMLPath = filePath

	// Check if the test already exists in RF so we can decide between updating and creating new one
	mappings, err := api.GetRFMLIDs()
	if err != nil {
		return err
	}

	testID, ok := mappings.GetID(parsedTest.RFMLID)
	if ok {
		parsedTest.TestID = testID
	} else {
		// Create an empty test
		log.Printf("Creating new test: %v", parsedTest.RFMLID)

		emptyTest := rainforest.RFTest{
			RFMLID: parsedTest.RFMLID,
			Title:  parsedTest.Title,
		}

		err = emptyTest.PrepareToUploadFromRFML(mappings)
		if err != nil {
			return err
		}

		err = api.CreateTest(&emptyTest)
		if err != nil {
			return err
		}
		log.Printf("Created new test: %v", parsedTest.RFMLID)
		// Refresh mappings
		mappings, err = api.GetRFMLIDs()
		if err != nil {
			return err
		}
		// Assign test ID
		testID, ok := mappings.GetID(parsedTest.RFMLID)
		if ok {
			parsedTest.TestID = testID
		} else {
			panic(fmt.Sprintf("Unable to map RFML ID to a primary ID: %v", parsedTest.RFMLID))
		}
	}

	if parsedTest.HasUploadableFiles() {
		err = api.ParseEmbeddedFiles(parsedTest)
		if err != nil {
			return err
		}
	}

	err = parsedTest.PrepareToUploadFromRFML(mappings)
	if err != nil {
		return err
	}

	// Update the steps
	log.Printf("Updating steps for test: %v", parsedTest.RFMLID)
	err = api.UpdateTest(parsedTest)
	if err != nil {
		return err
	}
	return nil
}

func uploadRFMLFiles(tests []*rainforest.RFTest, localOnly bool, api rfmlAPI) error {
	err := validateRFMLFiles(tests, localOnly, api)
	if err != nil {
		return err
	}

	// walk through the specifed directory (also subdirs) and pick the .rfml files
	// This will be used over and over again
	mappings, err := api.GetRFMLIDs()
	if err != nil {
		return err
	}
	var newTests []*rainforest.RFTest
	var parsedTests []*rainforest.RFTest

	for _, pTest := range tests {
		parsedTests = append(parsedTests, pTest)
		// Check if it's a new test or an existing one, because they need different treatment
		// to ensure we first add new ones and have IDs for potential embedds
		if _, ok := mappings.GetID(pTest.RFMLID); !ok {
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
		if err = <-errorsChan; err != nil {
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
		testID, ok := mappings.GetID(testToUpdate.RFMLID)
		if ok {
			testToUpdate.TestID = testID
		} else {
			panic(fmt.Sprintf("Unable to map RFML ID to primary ID: %v", testToUpdate.RFMLID))
		}

		if testToUpdate.HasUploadableFiles() {
			err = api.ParseEmbeddedFiles(testToUpdate)
			if err != nil {
				return err
			}
		}

		err = testToUpdate.PrepareToUploadFromRFML(mappings)
		if err != nil {
			return err
		}

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

func downloadRFML(c cliContext, client rfmlAPI) error {
	testDirectory := c.String("test-folder")
	absTestDirectory, err := prepareTestDirectory(testDirectory)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	var testIDs []int
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
		var tests []rainforest.RFTest
		filters := rainforest.RFTestFilters{
			Tags: c.StringSlice("tag"),
		}
		if c.Int("site-id") > 0 {
			filters.SiteID = c.Int("site-id")
		}
		if c.Int("folder-id") > 0 {
			filters.SmartFolderID = c.Int("folder-id")
		}
		if c.Int("feature-id") > 0 {
			filters.FeatureID = c.Int("feature-id")
		}
		if c.Int("run-group-id") > 0 {
			filters.RunGroupID = c.Int("run-group-id")
		}

		tests, err = client.GetTests(&filters)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		for _, t := range tests {
			testID := t.TestID
			testIDs = append(testIDs, testID)
		}
	}

	errorsChan := make(chan error)
	testIDChan := make(chan int, len(testIDs))
	testChan := make(chan *rainforest.RFTest, len(testIDs))

	for _, testID := range testIDs {
		testIDChan <- testID
	}
	close(testIDChan)

	for i := 0; i < rfmlDownloadConcurrency; i++ {
		go downloadRFTestWorker(testIDChan, errorsChan, testChan, client)
	}

	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	for i := 0; i < len(testIDs); i++ {
		select {
		case err = <-errorsChan:
			return cli.NewExitError(err.Error(), 1)
		case test := <-testChan:
			err = test.PrepareToWriteAsRFML(client, c.Bool("embed-tests"))
			if err != nil {
				return cli.NewExitError(err.Error(), 1)
			}

			paddedTestID := fmt.Sprintf("%010d", test.TestID)
			sanitizedTitle := sanitizeTestTitle(test.Title)
			fileName := fmt.Sprintf("%v_%v.rfml", paddedTestID, sanitizedTitle)
			rfmlFilePath := filepath.Join(absTestDirectory, fileName)

			var file *os.File
			file, err = os.Create(rfmlFilePath)
			if err != nil {
				return cli.NewExitError(err.Error(), 1)
			}

			writer := rainforest.NewRFMLWriter(file)
			err = writer.WriteRFMLTest(test)
			file.Close()
			if err != nil {
				return cli.NewExitError(err.Error(), 1)
			}

			log.Printf("Downloaded RFML test to %v", rfmlFilePath)
		}
	}

	return nil
}

func downloadRFTestWorker(testIDChan chan int, errorsChan chan error, testChan chan *rainforest.RFTest, client rfmlAPI) {
	for testID := range testIDChan {
		test, err := client.GetTest(testID)
		if err != nil {
			errorsChan <- err
			return
		}
		testChan <- test
	}
}

/*
	Helper Functions
*/

func prepareTestDirectory(testDir string) (string, error) {
	absTestDirectory, err := filepath.Abs(testDir)
	if err != nil {
		return "", err
	}

	dirStat, err := os.Stat(absTestDirectory)
	if os.IsNotExist(err) {
		log.Printf("Creating test directory: %v", absTestDirectory)
		os.MkdirAll(absTestDirectory, os.ModePerm)
	} else if err != nil {
		return "", err
	} else {
		if !dirStat.IsDir() {
			return "", fmt.Errorf("%v should be a directory", absTestDirectory)
		}
	}

	return absTestDirectory, nil
}

func sanitizeTestTitle(title string) string {
	title = strings.TrimSpace(title)
	title = strings.ToLower(title)

	// replace all non-alphanumeric character sequences with an underscore
	rep := regexp.MustCompile(`[^[[:alnum:]]+`)
	title = rep.ReplaceAllLiteralString(title, "_")

	if len(title) > 30 {
		return title[:30]
	}

	return title
}

func testCreationWorker(api rfmlAPI,
	testsToCreate <-chan *rainforest.RFTest, errorsChan chan<- error) {
	for test := range testsToCreate {
		log.Printf("Creating new test: %v", test.RFMLID)
		err := api.CreateTest(test)
		errorsChan <- err
	}
}

func testUpdateWorker(api rfmlAPI,
	testsToUpdate <-chan *rainforest.RFTest, errorsChan chan<- error) {
	for test := range testsToUpdate {
		log.Printf("Updating existing test: %v", test.RFMLID)
		err := api.UpdateTest(test)
		errorsChan <- err
	}
}
