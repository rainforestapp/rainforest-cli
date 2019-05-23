package rainforest

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
)

// RFMLReader reads form RFML formatted file.
// It exports some settings that can be set before parsing.
type RFMLReader struct {
	r *bufio.Reader
	// Version sets the RFML spec version, it's set by NewRFMLReader to the newest one.
	Version int
	// Sets the default value of redirect, that's used when it's not specified in RFML
	RedirectDefault bool
}

// parseError is a custom error implementing error interface for reporting RFML parsing errors.
type parseError struct {
	location string
	reason   string
}

func (e *parseError) Error() string {
	if _, err := strconv.Atoi(e.location); err == nil {
		// location is a line number
		return fmt.Sprintf("RFML parsing error in line %v: %v", e.location, e.reason)
	}
	// location is a field
	return fmt.Sprintf("RFML parsing error for test field \"%v\": %v", e.location, e.reason)
}

// NewRFMLReader returns RFML parser based on passed io.Reader - typically a RFML file.
func NewRFMLReader(r io.Reader) *RFMLReader {
	return &RFMLReader{
		r:               bufio.NewReader(r),
		Version:         1,
		RedirectDefault: true,
	}
}

// ReadAll parses whole RFML file using RFML version specified by Version parameter of reader
// and returns resulting RFTest
func (r *RFMLReader) ReadAll() (*RFTest, error) {
	// Default values
	parsedRFTest := &RFTest{
		State:   "enabled",
		Execute: true,
	}
	// Set up a new scanner to read in data line by line
	scanner := bufio.NewScanner(r.r)
	lineNum := 0
	// Temp variables where we put stuff while parsing
	currStep := make([]string, 0, 2)
	currStepRedirect := r.RedirectDefault
	for scanner.Scan() {
		lineNum++
		lineNumStr := strconv.Itoa(lineNum)

		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "#!") {
			if parsedRFTest.RFMLID != "" {
				return parsedRFTest, &parseError{lineNumStr, "Only one RFML ID may be specified"}
			}
			// Trim shebang and then take only first part of id before any spaces
			rfmlIDLine := strings.TrimSpace(line[2:])
			parsedRFTest.RFMLID = strings.Split(rfmlIDLine, " ")[0]
		} else if strings.HasPrefix(line, "#") {
			// Handle hashed lines
			content := line[1:]
			if strings.Contains(content, ":") {
				// Handle the key value pair
				split := strings.SplitN(content, ":", 2)
				key := strings.TrimSpace(split[0])
				value := strings.TrimSpace(split[1])
				switch key {
				case "title":
					parsedRFTest.Title = value
				case "start_uri":
					parsedRFTest.StartURI = value
				case "site_id":
					siteID, err := strconv.Atoi(value)
					if err != nil {
						return parsedRFTest, &parseError{lineNumStr, "Site ID must be a valid integer"}
					}
					parsedRFTest.SiteID = siteID
				case "tags":
					// If you split the empty string instead, you will get: []string{""}
					if len(value) == 0 {
						parsedRFTest.Tags = []string{}
						continue
					}

					splitTags := strings.Split(value, ",")
					strippedTags := make([]string, len(splitTags))
					for i, tag := range splitTags {
						strippedTags[i] = strings.TrimSpace(tag)
					}
					parsedRFTest.Tags = strippedTags
				case "browsers":
					// If you split the empty string instead, you will get: []string{""}
					if len(value) == 0 {
						parsedRFTest.Browsers = []string{}
						continue
					}

					splitBrowsers := strings.Split(value, ",")
					strippedBrowsers := make([]string, len(splitBrowsers))
					for i, tag := range splitBrowsers {
						strippedBrowsers[i] = strings.TrimSpace(tag)
					}
					parsedRFTest.Browsers = strippedBrowsers
				case "redirect":
					redirect, err := strconv.ParseBool(value)
					if err != nil {
						return parsedRFTest, &parseError{lineNumStr, "Redirect value must be a valid boolean"}
					}
					currStepRedirect = redirect
				case "feature_id":
					if value == "" {
						// If the value is empty, delete the feature
						parsedRFTest.FeatureID = deleteFeature
						continue
					}

					featureID, err := strconv.Atoi(value)
					if err != nil {
						return parsedRFTest, &parseError{lineNumStr, "Feature ID must be a valid integer"}
					}
					parsedRFTest.FeatureID = FeatureIDInt(featureID)
				case "state":
					parsedRFTest.State = value
				case "priority":
					value = strings.ToUpper(value)
					switch value {
					case "P1", "P2", "P3", "":
						parsedRFTest.Priority = value
					default:
						return parsedRFTest, &parseError{lineNumStr, "Priority value must be one of P1, P2, P3"}
					}
				case "execute":
					execute, err := strconv.ParseBool(value)
					if err != nil {
						return parsedRFTest, &parseError{lineNumStr, "Execute value must be a valid boolean"}
					}
					parsedRFTest.Execute = execute
				default:
					// If it doesn't match known key add it to description
					parsedRFTest.Description += strings.TrimSpace(content) + "\n"
				}
			} else {
				// If it'a a hashed line without key-value pair add it as a comment
				parsedRFTest.Description += strings.TrimSpace(content) + "\n"
			}
		} else {
			// Handle non prefixed lines
			// Here what we do depends on the fact if we have some step data collected already
			switch len(currStep) {
			case 0:
				if strings.HasPrefix(line, "-") {
					embeddedID := strings.TrimSpace(line[strings.Index(line, "-")+1:])
					embeddedStep := RFEmbeddedTest{embeddedID, currStepRedirect}
					parsedRFTest.Steps = append(parsedRFTest.Steps, embeddedStep)
					// Reset currStepRedirect
					currStepRedirect = r.RedirectDefault
				} else if line != "" {
					currStep = append(currStep, line)
				}
			case 1:
				if strings.Contains(line, "?") {
					currStep = append(currStep, line)
				} else {
					return parsedRFTest, &parseError{lineNumStr, "Each step must contain a question, with a `?`"}
				}
			case 2:
				if line == "" {
					parsedStep := RFTestStep{currStep[0], currStep[1], currStepRedirect}
					parsedRFTest.Steps = append(parsedRFTest.Steps, parsedStep)
					// Reset temp vars to defaults
					currStep = make([]string, 0, 2)
					currStepRedirect = r.RedirectDefault
				} else {
					return parsedRFTest, &parseError{lineNumStr, "Steps must be separated with empty lines"}
				}
			}
		}
	}

	// Check if parsing stopped before adding a step
	if len(currStep) == 1 {
		return parsedRFTest, &parseError{strconv.Itoa(lineNum), "Must have a corresponding question with your action."}
	}

	if len(currStep) == 2 {
		parsedStep := RFTestStep{currStep[0], currStep[1], currStepRedirect}
		parsedRFTest.Steps = append(parsedRFTest.Steps, parsedStep)
	}

	if parsedRFTest.RFMLID == "" {
		return parsedRFTest, &parseError{"#!", "RFML ID is required for .rfml files. Specify it using #! followed by a unique RFML ID"}
	}

	if parsedRFTest.Title == "" {
		return parsedRFTest, &parseError{"# title", "Title is required for .rfml files. Specify it using \"# title: \" followed by your test's title."}
	}

	return parsedRFTest, nil
}

// RFMLWriter writes a RFML formatted test to a given file.
type RFMLWriter struct {
	w *bufio.Writer
	// Version sets the RFML spec version
	Version int
}

// NewRFMLWriter returns RFML writer based on passed io.Writer - typically a RFML file.
func NewRFMLWriter(w io.Writer) *RFMLWriter {
	return &RFMLWriter{
		w:       bufio.NewWriter(w),
		Version: 1,
	}
}

// WriteRFMLTest writes a given RFTest to its writer in the given RFML version.
func (r *RFMLWriter) WriteRFMLTest(test *RFTest) error {
	writer := r.w
	headerTemplate := `#! %v
# title: %v
# start_uri: %v
`

	header := fmt.Sprintf(headerTemplate, test.RFMLID, test.Title, test.StartURI)
	_, err := writer.WriteString(header)

	if err != nil {
		return err
	}

	if test.SiteID != 0 {
		_, err = writer.WriteString("# site_id: " + strconv.Itoa(test.SiteID) + "\n")

		if err != nil {
			return err
		}
	}

	if featureID := int(test.FeatureID); featureID != 0 {
		_, err = writer.WriteString("# feature_id: " + strconv.Itoa(featureID) + "\n")

		if err != nil {
			return err
		}
	}

	if len(test.Tags) > 0 {
		tags := strings.Join(test.Tags, ", ")
		tagsHeader := fmt.Sprintf("# tags: %v\n", tags)

		_, err = writer.WriteString(tagsHeader)

		if err != nil {
			return err
		}
	}

	if len(test.Browsers) > 0 {
		browsers := strings.Join(test.Browsers, ", ")
		browsersHeader := fmt.Sprintf("# browsers: %v\n", browsers)

		_, err = writer.WriteString(browsersHeader)

		if err != nil {
			return err
		}
	}

	if test.State == "disabled" {
		_, err = writer.WriteString("# state: disabled\n")
		if err != nil {
			return err
		}
	}

	switch test.Priority {
	case "P1", "P2", "P3":
		_, err = writer.WriteString(fmt.Sprintf("# priority: %s\n", test.Priority))
		if err != nil {
			return err
		}
	}

	if !test.Execute {
		_, err = writer.WriteString("# execute: false\n")
		if err != nil {
			return err
		}
	}

	if test.Description != "" {
		_, err = writer.WriteString("# " + strings.Replace(test.Description, "\n", "\n# ", -1) + "\n")

		if err != nil {
			return err
		}
	}

	firstStepProcessed := false
	processStep := func(idx int, step RFTestStep) string {
		stepText := ""
		if idx > 0 && firstStepProcessed == false {
			stepText = stepText + fmt.Sprintf("# redirect: %v\n", step.Redirect)
		}
		action := strings.Replace(step.Action, "\n", " ", -1)
		response := strings.Replace(step.Response, "\n", " ", -1)
		firstStepProcessed = true

		return stepText + action + "\n" + response
	}

	for idx, step := range test.Steps {
		var stepText string
		switch step.(type) {
		case RFTestStep:
			stepText = processStep(idx, step.(RFTestStep))
		case RFEmbeddedTest:
			embeddedTest := step.(RFEmbeddedTest)
			if idx > 0 {
				stepText = "# redirect: " + strconv.FormatBool(embeddedTest.Redirect) + "\n"
			}
			stepText = stepText + "- " + embeddedTest.RFMLID
		}

		_, err = writer.WriteString("\n" + stepText + "\n")

		if err != nil {
			return err
		}
	}

	// Writes buffered data to the underlying io.Writer
	err = writer.Flush()

	if err != nil {
		return err
	}

	return nil
}

// ParseEmbeddedFiles replaces file step variable paths with values expected
// by Rainforest. eg: {{ file.screenshot(my_screenshot.gif) }} would be translated
// to the format {{ file.screenshot(FILE_ID, FILE_SIGNATURE) }}.
func (c *Client) ParseEmbeddedFiles(test *RFTest) error {
	if test.TestID == 0 {
		return fmt.Errorf("Cannot parse embedded files without a test ID.")
	}

	uploadedFiles, err := c.getUploadedFiles(test.TestID)
	if err != nil {
		return err
	}

	digestToFileMap := map[string]uploadedFile{}
	for _, f := range uploadedFiles {
		digestToFileMap[f.Digest] = f
	}

	replaceEmbeddedFilePaths := func(text string, embeddedFiles []embeddedFile) (string, error) {
		out := text
		for _, embed := range embeddedFiles {
			filePath := embed.path
			if strings.HasPrefix(filePath, "~/") {
				var usr *user.User
				usr, err = user.Current()
				if err != nil {
					return "", err
				}
				filePath = filepath.Join(usr.HomeDir, filePath[2:])
			} else if test.RFMLPath == "" {
				return "", fmt.Errorf("Cannot parse relative file path %v for RFML test %v. RFMLPath field cannot be blank.", filePath, test.RFMLID)
			} else {
				rfmlDirectory := filepath.Dir(test.RFMLPath)
				filePath = filepath.Join(rfmlDirectory, filePath)
			}

			filePath, err = filepath.Abs(filePath)
			if err != nil {
				return "", err
			}

			var file *os.File
			file, err = os.Open(filePath)

			if os.IsNotExist(err) {
				log.Printf("Error for test: %v\n\tNo such file exists: %v\n", test.RFMLPath, filePath)
				continue
			} else if err != nil {
				return "", err
			}

			var data []byte
			data, err = ioutil.ReadAll(file)
			defer file.Close()
			if err != nil {
				return "", err
			}

			checksum := md5.Sum(data)
			fileDigest := hex.EncodeToString(checksum[:])
			uploadedFileInfo, ok := digestToFileMap[fileDigest]
			// TODO: Check mime type as well
			if !ok {
				// File has not been uploaded before
				// Upload to RF
				var awsInfo *awsFileInfo
				awsInfo, err = c.createTestFile(test.TestID, file, data)
				if err != nil {
					return "", err
				}
				// Upload to AWS
				err = c.uploadEmbeddedFile(filepath.Base(filePath), data, awsInfo)
				if err != nil {
					return "", err
				}
				uploadedFileInfo = uploadedFile{
					ID:        awsInfo.FileID,
					Signature: awsInfo.FileSignature,
					Digest:    fileDigest,
				}
				// Add to the mappings for future reference
				digestToFileMap[fileDigest] = uploadedFileInfo
			}

			sig := uploadedFileInfo.Signature[0:6]
			var replacement string
			if embed.stepVar == "screenshot" {
				replacement = fmt.Sprintf("{{ file.screenshot(%v, %v) }}", uploadedFileInfo.ID, sig)
			} else if embed.stepVar == "download" {
				replacement = fmt.Sprintf("{{ file.download(%v, %v, %v) }}", uploadedFileInfo.ID, sig, filepath.Base(filePath))
			}

			out = strings.Replace(out, embed.text, replacement, 1)
		}

		return out, nil
	}

	for idx, step := range test.Steps {
		s, ok := step.(RFTestStep)
		if ok && s.hasUploadableFiles() {
			if embeddedFiles := s.embeddedFilesInAction(); len(embeddedFiles) > 0 {
				s.Action, err = replaceEmbeddedFilePaths(s.Action, embeddedFiles)
				if err != nil {
					return err
				}
			}

			if embeddedFiles := s.embeddedFilesInResponse(); len(embeddedFiles) > 0 {
				s.Response, err = replaceEmbeddedFilePaths(s.Response, embeddedFiles)
				if err != nil {
					return err
				}
			}
			test.Steps[idx] = s
		}
	}

	return nil
}
