package rainforest

import (
	"bytes"
	"io/ioutil"
	"strconv"
	"strings"
	"testing"
)

func TestWriteRFMLTest(t *testing.T) {
	var buffer bytes.Buffer
	writer := NewRFMLWriter(&buffer)

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
}
