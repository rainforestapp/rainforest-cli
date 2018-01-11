%{
package rfml

import (
	"fmt"
	"strings"
)

import "github.com/rainforestapp/rainforest-cli/rainforest"

var curTest *rainforest.RFTest
%}

%union {
  str string
  bool bool
}

%token _STRING
%token _TITLE
%token _START_URI
%token _TAGS
%token _BROWSERS
%token _REDIRECT
%token _EXECUTE
%token _EOF

%%

file:
                metadata '\n' steps _EOF { return 0 }
                ;

metadata:
                id_header headers
                ;

id_header:
                '#' '!' _STRING '\n' { curTest.RFMLID = $3.str }
                ;

headers:
                /* empty */
        |       '#' header headers
                ;

header:
                headerval
        |       _TITLE ':' headerval { curTest.Title = $3.str }
        |       _START_URI ':' headerval { curTest.StartURI = $3.str }
        |       _TAGS ':' headerval { curTest.Tags = parseList($3.str) }
        |       _BROWSERS ':' headerval { curTest.Browsers = parseList($3.str) }
        |       _EXECUTE ':' headerval { curTest.Execute = parseBool($3.str) }
                ;

headerval:
                '\n' { $$.str = "" }
        |       _STRING '\n' { $$.str = $1.str }
                ;

steps:
                /* empty */
        |       '\n' steps
        |       embedded_test steps
        |       step steps
                ;

embedded_test:
                redirect_header '-' _STRING '\n' { appendEmbeddedTest($3.str, $1.bool) }
                ;

step:
                redirect_header action response { appendStep($2.str, $3.str, $1.bool) }
                ;

redirect_header:
                /* empty */ { $$.bool = true }
        |       '#' _REDIRECT ':' _STRING '\n' { $$.bool = parseBool($4.str) }
                ;

action:
                _STRING '\n' { $$.str = $1.str }
                ;

response:
                _STRING '\n' { $$.str = $1.str }
                ;

%%

func parseList(str string) []string {
	result := []string{}

	for _, s := range strings.Split(str, ",") {
      s = strings.TrimSpace(s)
      if s != "" {
          result = append(result, s)
      }
	}

	return result
}

func parseBool(str string) bool {
	s := strings.TrimSpace(str)
	if s == "true" {
		return true
	}
	if s == "false" {
		return false
	}
  panic(fmt.Sprintf("invalid boolean value: %s", str))
}

func appendEmbeddedTest(id string, redirect bool) {
	step := rainforest.RFEmbeddedTest{RFMLID: id, Redirect: redirect}
  curTest.Steps = append(curTest.Steps, step)
}

func appendStep(action, response string, redirect bool) {
	step := rainforest.RFTestStep{Action: action, Response: response, Redirect: redirect}
  curTest.Steps = append(curTest.Steps, step)
}
