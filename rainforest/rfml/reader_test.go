package rfml

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/rainforestapp/rainforest-cli/rainforest"
)

var fixtures = map[string]*rainforest.RFTest{
	"a": &rainforest.RFTest{
		RFMLID:      "test-a",
		Title:       "Test A",
		StartURI:    "/",
		State:       "enabled",
		Description: "",
		Execute:     true,
		Tags:        []string{"run-me", "fixme"},
		Browsers:    []string{"chrome_1440_900"},
		Steps: []interface{}{
			rainforest.RFEmbeddedTest{
				RFMLID:   "login",
				Redirect: true,
			},
			rainforest.RFTestStep{
				Action:   "Do a thing.",
				Response: "Did something happen?",
				Redirect: false,
			},
			rainforest.RFTestStep{
				Action:   "Do something else.",
				Response: "Did something else happen?",
				Redirect: true,
			},
		},
	},
	"login-test": &rainforest.RFTest{
		RFMLID:      "login",
		Title:       "Login",
		StartURI:    "/login",
		SiteID:      12345,
		State:       "enabled",
		Description: "",
		Execute:     false,
		Tags:        []string{},
		Steps: []interface{}{
			rainforest.RFTestStep{
				Action:   "Log in with a username and password.",
				Response: "Did you log in successfully?",
				Redirect: true,
			},
		},
	},
}

func TestLex(t *testing.T) {
	// This is really just a sanity check, the actual parser will break if the
	// lexer is broken.
	f, err := os.Open("fixtures/login-test.rfml")
	if err != nil {
		t.Fatal(err)
	}

	want := []int{
		'#', '!', _STRING, '\n',
		'#', _TITLE, ':', _STRING, '\n',
		'#', _START_URI, ':', _STRING, '\n',
		'#', _SITE_ID, ':', _INTEGER, '\n',
		'#', _TAGS, ':', '\n',
		'#', _EXECUTE, ':', _BOOL, '\n',
		'#', '\n',
		'\n',
		_STRING, '\n',
		_STRING, '\n',
		_EOF,
	}
	got := []int{}

	r := NewReader(f)

	for {
		tok := r.Lex(&yySymType{})
		if r.parseError != nil {
			t.Fatal("Parse error:", r.parseError)
		}
		got = append(got, tok)
		if tok == _EOF {
			break
		}
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("Lex error: want:\n %v\ngot:\n%v\n", want, got)
	}
}

func TestReadFixtures(t *testing.T) {
	for fname, want := range fixtures {
		f, err := os.Open(fmt.Sprintf("fixtures/%v.rfml", fname))
		if err != nil {
			t.Error("Open error:", err)
			continue
		}

		r := NewReader(f)
		got, err := r.ReadAll()
		if err != nil {
			t.Error("Read error:", err)
			continue
		}

		if !reflect.DeepEqual(want, got) {
			t.Errorf("RFML parse error!\nwant:\n %+v\n\ngot:\n %+v", want, got)
		}
	}
}

func TestReadAll(t *testing.T) {
	const deleteFeature = -1

	validSteps := []interface{}{
		rainforest.RFTestStep{
			Action:   "First Action",
			Response: "First Question?",
			Redirect: true,
		},
		rainforest.RFTestStep{
			Action:   "Second Action",
			Response: "Second Question?",
			Redirect: true,
		},
		rainforest.RFEmbeddedTest{
			RFMLID:   "embedded_id",
			Redirect: true,
		},
	}

	validTestValues := rainforest.RFTest{
		RFMLID:    "my_rfml_id",
		Title:     "my_title",
		StartURI:  "/testing",
		SiteID:    12345,
		FeatureID: 98765,
		State:     "enabled",
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
		validSteps[0].(rainforest.RFTestStep).Action,
		validSteps[0].(rainforest.RFTestStep).Response,
		validSteps[1].(rainforest.RFTestStep).Action,
		validSteps[1].(rainforest.RFTestStep).Response,
		validSteps[2].(rainforest.RFEmbeddedTest).RFMLID,
	)

	r := strings.NewReader(testText)
	reader := NewReader(r)
	rfTest, err := reader.ReadAll()
	if err != nil {
		t.Fatal(err.Error())
	}

	if !reflect.DeepEqual(*rfTest, validTestValues) {
		t.Errorf("Incorrect values for RFTest.\nGot %#v\nWant %#v", rfTest, validTestValues)
	}

	// Test is disabled
	testText = fmt.Sprintf(`#! %v
# title: %v
# start_uri: %v
# state: %v
`,
		validTestValues.RFMLID,
		validTestValues.Title,
		validTestValues.StartURI,
		"disabled",
	)

	r = strings.NewReader(testText)
	reader = NewReader(r)
	rfTest, err = reader.ReadAll()
	if err != nil {
		t.Fatal(err.Error())
	}

	if rfTest.State != "disabled" {
		t.Errorf("Incorrect test state. Got %v, Want %v", rfTest.State, "disabled")
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
	reader = NewReader(r)
	rfTest, err = reader.ReadAll()
	if err != nil {
		t.Fatal(err.Error())
	}

	if rfTest.State != "enabled" {
		t.Errorf("Incorrect test state. Got %v, Want %v", rfTest.State, "enabled")
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
		validSteps[0].(rainforest.RFTestStep).Action,
		validSteps[0].(rainforest.RFTestStep).Response,
	)

	r = strings.NewReader(testText)
	reader = NewReader(r)
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
		validSteps[0].(rainforest.RFTestStep).Action,
		validSteps[0].(rainforest.RFTestStep).Response)
	r = strings.NewReader(testText)
	reader = NewReader(r)
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
		validSteps[0].(rainforest.RFTestStep).Action,
		validSteps[0].(rainforest.RFTestStep).Response,
	)

	r = strings.NewReader(testText)
	reader = NewReader(r)
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
		validSteps[0].(rainforest.RFTestStep).Action,
		validSteps[0].(rainforest.RFTestStep).Response,
	)

	r = strings.NewReader(testText)
	reader = NewReader(r)
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
		validSteps[0].(rainforest.RFTestStep).Action,
		validSteps[0].(rainforest.RFTestStep).Response,
	)

	r = strings.NewReader(testText)
	reader = NewReader(r)
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
		validSteps[0].(rainforest.RFTestStep).Action,
		validSteps[0].(rainforest.RFTestStep).Response,
	)

	r = strings.NewReader(testText)
	reader = NewReader(r)
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
