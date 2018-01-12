package rfml

import (
	"fmt"
	"os"
	"reflect"
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

func TestReadAll(t *testing.T) {
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
