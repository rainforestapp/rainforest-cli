%{
package rfml

import (
	"strings"
)

import "github.com/rainforestapp/rainforest-cli/rainforest"

var curTest *rainforest.RFTest
%}

%union {
  strval string
  boolval bool
  steplist []interface{}
  step rainforest.RFTestStep
  embedded_test rainforest.RFEmbeddedTest
}

%token <strval> _STRING
%token <boolval> _BOOL
%token _TITLE
%token _START_URI
%token _TAGS
%token _BROWSERS
%token _REDIRECT
%token _EXECUTE
%token _EOF

%type   <strval>          headerval
%type   <strval>          action
%type   <strval>          response
%type   <boolval>         redirect_header
%type   <steplist>        steps
%type   <step>           step
%type   <embedded_test>  embedded_test
%type   <steplist>       emptyline

%start file

%%

file : metadata '\n' steps _EOF { return finalizeTest($3) }
     ;

metadata : id_header headers
         ;

id_header : '#' '!' _STRING '\n' { curTest.RFMLID = $3 }
          ;

headers : /* empty */
        |       '#' header headers
                ;

header : headerval
        |       _TITLE ':' headerval { curTest.Title = $3 }
        |       _START_URI ':' headerval { curTest.StartURI = $3 }
        |       _TAGS ':' headerval { curTest.Tags = parseList($3) }
        |       _BROWSERS ':' headerval { curTest.Browsers = parseList($3) }
        |       _EXECUTE ':' _BOOL '\n' { curTest.Execute = $3 }
                ;

headerval : '\n' { $$ = "" }
        |       _STRING '\n' { $$ = $1 }
                ;

steps : /* empty */ { $$ = []interface{}{} }
        |       emptyline steps { $$ = $2 }
        |       embedded_test steps { $$ = append([]interface{}{$1}, $2...) }
        |       step steps { $$ = append([]interface{}{$1}, $2...) }
                ;

emptyline : '\n' { $$ = []interface{}{} }

embedded_test : redirect_header '-' _STRING '\n' { $$ = rainforest.RFEmbeddedTest{$3, $1} }
                ;

step :
                redirect_header action response { $$ = rainforest.RFTestStep{$2, $3, $1} }
                ;

redirect_header : /* empty */ { $$ = true }
        |       '#' _REDIRECT ':' _BOOL '\n' { $$ = $4 }
                ;

action : _STRING '\n' { $$ = $1 }
                ;

response : _STRING '\n' { $$ = $1 }
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

func finalizeTest(steps []interface{}) int {
  curTest.Steps = steps
  return 0
}
