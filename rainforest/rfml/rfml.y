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
    intval int
    steplist []interface{}
    step rainforest.RFTestStep
    embedded_test rainforest.RFEmbeddedTest
}

%token  <strval>    _STRING
%token  <boolval>   _BOOL
%token  <intval>    _INTEGER
%token  _TITLE
%token  _START_URI
%token  _TAGS
%token  _BROWSERS
%token  _REDIRECT
%token  _EXECUTE
%token  _FEATURE_ID
%token  _SITE_ID
%token  _STATE
%token  _EOF

%type   <strval>            action
%type   <strval>            response
%type   <boolval>           redirect_header
%type   <steplist>          steplist
%type   <steplist>          steps
%type   <step>              step
%type   <embedded_test>     embedded_test
%type   <steplist>          emptyline
%type   <strval>            headerstr
%type   <intval>            headerint

%start file

%%

file        :   metadata steplist _EOF          { return finalizeTest($2) }
            ;

metadata    :   id_header headers
            ;

id_header   :   '#' '!' _STRING '\n'            { curTest.RFMLID = $3 }
            ;

headers     :   /* empty */
            |   '#' header headers
            ;

header      :   headerstr                       { curTest.Description += $1 + "\n" }
            |   _TITLE ':' headerstr            { curTest.Title = $3 }
            |   _START_URI ':' headerstr        { curTest.StartURI = $3 }
            |   _TAGS ':' headerstr             { curTest.Tags = parseList($3) }
            |   _BROWSERS ':' headerstr         { curTest.Browsers = parseList($3) }
            |   _EXECUTE ':' _BOOL '\n'         { curTest.Execute = $3 }
            |   _SITE_ID ':' headerint          { curTest.SiteID = $3 }
            |   _STATE ':' headerstr            { curTest.State = $3 }
            |   _FEATURE_ID ':' '\n'            { curTest.FeatureID = rainforest.FeatureIDInt(-1) }
            |   _FEATURE_ID ':' _INTEGER '\n'   { curTest.FeatureID = rainforest.FeatureIDInt($3) }
            ;

headerstr   :   '\n'                            { $$ = "" }
            |   _STRING '\n'                    { $$ = $1 }
            ;

headerint   :   '\n'                            { $$ = 0 }
            |   _INTEGER '\n'                   { $$ = $1 }
            ;

steplist    :   /* empty */                     { $$ = []interface{}{} }
            |   '\n' steps                      { $$ = $2 }
            ;

steps       :   /* empty */                     { $$ = []interface{}{} }
            |   emptyline steps                 { $$ = $2 }
            |   embedded_test steps             { $$ = append([]interface{}{$1}, $2...) }
            |   step steps                      { $$ = append([]interface{}{$1}, $2...) }
            ;

emptyline   :   '\n'                            { $$ = []interface{}{} }

embedded_test : redirect_header '-' _STRING step_end
                                                { $$ = rainforest.RFEmbeddedTest{$3, $1} }
            ;

step        :   redirect_header action response { $$ = rainforest.RFTestStep{$2, $3, $1} }
            ;

redirect_header : /* empty */                   { $$ = true }
            |   '#' _REDIRECT ':' _BOOL '\n'    { $$ = $4 }
            ;

action          :   _STRING '\n'                { $$ = $1 }
            ;

response    :   _STRING step_end                { $$ = $1 }
            ;

step_end    :   '\n' | _EOF
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
