package rainforest

import (
	"bufio"
	"fmt"
	"io"
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
	line   int
	reason string
}

func (e *parseError) Error() string {
	return fmt.Sprintf("RFML parsing error in line %v: %v", e.line, e.reason)
}

// NewRFMLReader returns RFML parser based on passed io.Reader - typically a RFML file.
func NewRFMLReader(r io.Reader) *RFMLReader {
	return &RFMLReader{
		r:               bufio.NewReader(r),
		Version:         1,
		RedirectDefault: false,
	}
}

// ReadAll parses whole RFML file using RFML version specified by Version parameter of reader
// and returns resulting RFTest
func (r *RFMLReader) ReadAll() (*RFTest, error) {
	parsedRFTest := &RFTest{}
	// Set up a new scanner to read in data line by line
	scanner := bufio.NewScanner(r.r)
	lineNum := 0
	// Temp variables where we put stuff while parsing
	currStep := make([]string, 0, 2)
	currStepRedirect := r.RedirectDefault
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "#!") {
			// Handle shebang
			parsedRFTest.RFMLID = strings.TrimSpace(line[2:])
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
						return parsedRFTest, &parseError{lineNum, "Site ID must be a valid integer"}
					}
					parsedRFTest.SiteID = siteID
				case "tags":
					splitTags := strings.Split(value, ",")
					strippedTags := make([]string, len(splitTags))
					for i, tag := range splitTags {
						strippedTags[i] = strings.TrimSpace(tag)
					}
					parsedRFTest.Tags = strippedTags
				case "browsers":
					splitBrowsers := strings.Split(value, ",")
					strippedBrowsers := make([]string, len(splitBrowsers))
					for i, tag := range splitBrowsers {
						strippedBrowsers[i] = strings.TrimSpace(tag)
					}
					parsedRFTest.Browsers = strippedBrowsers
				case "redirect":
					redirect, err := strconv.ParseBool(value)
					if err != nil {
						return parsedRFTest, &parseError{lineNum, "Redirect value must be a valid boolean"}
					}
					currStepRedirect = redirect
				default:
					// If it doesn't match known key add it to description
					parsedRFTest.Description += content + "\n"
				}
			} else {
				// If it'a a hashed line without key-value pair add it as a comment
				parsedRFTest.Description += content + "\n"
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
					return parsedRFTest, &parseError{lineNum, "Each step must contain a question, with a `?`"}
				}
			case 2:
				if line == "" {
					parsedStep := RFTestStep{currStep[0], currStep[1], currStepRedirect}
					parsedRFTest.Steps = append(parsedRFTest.Steps, parsedStep)
					// Reset temp vars to defaults
					currStep = make([]string, 0, 2)
					currStepRedirect = r.RedirectDefault
				} else {
					return parsedRFTest, &parseError{lineNum, "Steps must be separated with empty lines"}
				}
			}
		}
	}
	if parsedRFTest.RFMLID == "" {
		return parsedRFTest, &parseError{1, "RFML ID is required for .rfml files, specify it using #!"}
	}
	parsedRFTest.mapBrowsers()
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

	// TODO: Steps

	// Writes buffered data to the underlying io.Writer
	err = writer.Flush()

	if err != nil {
		return err
	}

	return nil
}
