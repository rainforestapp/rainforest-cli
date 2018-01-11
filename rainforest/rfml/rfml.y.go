//line rfml.y:2
package rfml

import __yyfmt__ "fmt"

//line rfml.y:2
import (
	"fmt"
	"strings"
)

import "github.com/rainforestapp/rainforest-cli/rainforest"

var curTest *rainforest.RFTest

//line rfml.y:14
type yySymType struct {
	yys           int
	strval        string
	boolval       bool
	steplist      []interface{}
	step          rainforest.RFTestStep
	embedded_test rainforest.RFEmbeddedTest
}

const _STRING = 57346
const _TITLE = 57347
const _START_URI = 57348
const _TAGS = 57349
const _BROWSERS = 57350
const _REDIRECT = 57351
const _EXECUTE = 57352
const _EOF = 57353

var yyToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"_STRING",
	"_TITLE",
	"_START_URI",
	"_TAGS",
	"_BROWSERS",
	"_REDIRECT",
	"_EXECUTE",
	"_EOF",
	"'\\n'",
	"'#'",
	"'!'",
	"':'",
	"'-'",
}
var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

//line rfml.y:94

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

func finalizeTest(steps []interface{}) int {
	curTest.Steps = steps
	return 0
}

//line yacctab:1
var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
	-1, 5,
	11, 14,
	-2, 21,
	-1, 10,
	11, 14,
	-2, 21,
	-1, 11,
	11, 14,
	-2, 21,
	-1, 12,
	11, 14,
	-2, 21,
}

const yyPrivate = 57344

const yyLast = 57

var yyAct = [...]int{

	17, 24, 18, 19, 20, 21, 32, 22, 6, 23,
	46, 39, 33, 38, 37, 36, 35, 9, 30, 13,
	15, 55, 8, 24, 7, 34, 4, 53, 27, 28,
	29, 23, 52, 45, 41, 40, 47, 48, 49, 50,
	51, 5, 26, 54, 44, 42, 25, 16, 3, 2,
	1, 10, 11, 12, 14, 43, 31,
}
var yyPact = [...]int{

	13, -1000, 29, 11, 8, 7, -1000, -3, 42, 31,
	7, 7, 7, -1000, 2, 3, 11, -1000, 1, 0,
	-1, -2, -4, -1000, 23, 22, -1000, -1000, -1000, -1000,
	41, 40, 21, -5, -1000, 19, 19, 19, 19, 19,
	-1000, -1000, 20, -1000, 15, -1000, 39, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, 9, -1000,
}
var yyPgo = [...]int{

	0, 0, 56, 55, 54, 17, 53, 52, 51, 50,
	49, 48, 8, 47,
}
var yyR1 = [...]int{

	0, 9, 10, 11, 12, 12, 13, 13, 13, 13,
	13, 13, 1, 1, 5, 5, 5, 5, 8, 7,
	6, 4, 4, 2, 3,
}
var yyR2 = [...]int{

	0, 4, 2, 4, 0, 3, 1, 3, 3, 3,
	3, 3, 1, 2, 0, 2, 2, 2, 1, 4,
	3, 0, 5, 2, 2,
}
var yyChk = [...]int{

	-1000, -9, -10, -11, 13, 12, -12, 13, 14, -5,
	-8, -7, -6, 12, -4, 13, -13, -1, 5, 6,
	7, 8, 10, 12, 4, 4, 11, -5, -5, -5,
	16, -2, 4, 9, -12, 15, 15, 15, 15, 15,
	12, 12, 4, -3, 4, 12, 15, -1, -1, -1,
	-1, -1, 12, 12, 4, 12,
}
var yyDef = [...]int{

	0, -2, 0, 4, 0, -2, 2, 0, 0, 0,
	-2, -2, -2, 18, 0, 0, 4, 6, 0, 0,
	0, 0, 0, 12, 0, 0, 1, 15, 16, 17,
	0, 0, 0, 0, 5, 0, 0, 0, 0, 0,
	13, 3, 0, 20, 0, 23, 0, 7, 8, 9,
	10, 11, 19, 24, 0, 22,
}
var yyTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	12, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 14, 3, 13, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 16, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 15,
}
var yyTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
}
var yyTok3 = [...]int{
	0,
}

var yyErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

//line yaccpar:1

/*	parser for yacc output	*/

var (
	yyDebug        = 0
	yyErrorVerbose = false
)

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

type yyParser interface {
	Parse(yyLexer) int
	Lookahead() int
}

type yyParserImpl struct {
	lval  yySymType
	stack [yyInitialStackSize]yySymType
	char  int
}

func (p *yyParserImpl) Lookahead() int {
	return p.char
}

func yyNewParser() yyParser {
	return &yyParserImpl{}
}

const yyFlag = -1000

func yyTokname(c int) string {
	if c >= 1 && c-1 < len(yyToknames) {
		if yyToknames[c-1] != "" {
			return yyToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func yyStatname(s int) string {
	if s >= 0 && s < len(yyStatenames) {
		if yyStatenames[s] != "" {
			return yyStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func yyErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !yyErrorVerbose {
		return "syntax error"
	}

	for _, e := range yyErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + yyTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := yyPact[state]
	for tok := TOKSTART; tok-1 < len(yyToknames); tok++ {
		if n := base + tok; n >= 0 && n < yyLast && yyChk[yyAct[n]] == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if yyDef[state] == -2 {
		i := 0
		for yyExca[i] != -1 || yyExca[i+1] != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; yyExca[i] >= 0; i += 2 {
			tok := yyExca[i]
			if tok < TOKSTART || yyExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if yyExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += yyTokname(tok)
	}
	return res
}

func yylex1(lex yyLexer, lval *yySymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = yyTok1[0]
		goto out
	}
	if char < len(yyTok1) {
		token = yyTok1[char]
		goto out
	}
	if char >= yyPrivate {
		if char < yyPrivate+len(yyTok2) {
			token = yyTok2[char-yyPrivate]
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		token = yyTok3[i+0]
		if token == char {
			token = yyTok3[i+1]
			goto out
		}
	}

out:
	if token == 0 {
		token = yyTok2[1] /* unknown char */
	}
	if yyDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", yyTokname(token), uint(char))
	}
	return char, token
}

func yyParse(yylex yyLexer) int {
	return yyNewParser().Parse(yylex)
}

func (yyrcvr *yyParserImpl) Parse(yylex yyLexer) int {
	var yyn int
	var yyVAL yySymType
	var yyDollar []yySymType
	_ = yyDollar // silence set and not used
	yyS := yyrcvr.stack[:]

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yyrcvr.char = -1
	yytoken := -1 // yyrcvr.char translated into internal numbering
	defer func() {
		// Make sure we report no lookahead when not parsing.
		yystate = -1
		yyrcvr.char = -1
		yytoken = -1
	}()
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	if yyDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", yyTokname(yytoken), yyStatname(yystate))
	}

	yyp++
	if yyp >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyS[yyp] = yyVAL
	yyS[yyp].yys = yystate

yynewstate:
	yyn = yyPact[yystate]
	if yyn <= yyFlag {
		goto yydefault /* simple state */
	}
	if yyrcvr.char < 0 {
		yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
	}
	yyn += yytoken
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = yyAct[yyn]
	if yyChk[yyn] == yytoken { /* valid shift */
		yyrcvr.char = -1
		yytoken = -1
		yyVAL = yyrcvr.lval
		yystate = yyn
		if Errflag > 0 {
			Errflag--
		}
		goto yystack
	}

yydefault:
	/* default state action */
	yyn = yyDef[yystate]
	if yyn == -2 {
		if yyrcvr.char < 0 {
			yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
		}

		/* look through exception table */
		xi := 0
		for {
			if yyExca[xi+0] == -1 && yyExca[xi+1] == yystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			yyn = yyExca[xi+0]
			if yyn < 0 || yyn == yytoken {
				break
			}
		}
		yyn = yyExca[xi+1]
		if yyn < 0 {
			goto ret0
		}
	}
	if yyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			yylex.Error(yyErrorMessage(yystate, yytoken))
			Nerrs++
			if yyDebug >= 1 {
				__yyfmt__.Printf("%s", yyStatname(yystate))
				__yyfmt__.Printf(" saw %s\n", yyTokname(yytoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for yyp >= 0 {
				yyn = yyPact[yyS[yyp].yys] + yyErrCode
				if yyn >= 0 && yyn < yyLast {
					yystate = yyAct[yyn] /* simulate a shift of "error" */
					if yyChk[yystate] == yyErrCode {
						goto yystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if yyDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", yyS[yyp].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", yyTokname(yytoken))
			}
			if yytoken == yyEofCode {
				goto ret1
			}
			yyrcvr.char = -1
			yytoken = -1
			goto yynewstate /* try again in the same state */
		}
	}

	/* reduction by production yyn */
	if yyDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", yyn, yyStatname(yystate))
	}

	yynt := yyn
	yypt := yyp
	_ = yypt // guard against "declared and not used"

	yyp -= yyR2[yyn]
	// yyp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if yyp+1 >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyVAL = yyS[yyp+1]

	/* consult goto table to find next state */
	yyn = yyR1[yyn]
	yyg := yyPgo[yyn]
	yyj := yyg + yyS[yyp].yys + 1

	if yyj >= yyLast {
		yystate = yyAct[yyg]
	} else {
		yystate = yyAct[yyj]
		if yyChk[yystate] != -yyn {
			yystate = yyAct[yyg]
		}
	}
	// dummy call; replaced with literal code
	switch yynt {

	case 1:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line rfml.y:44
		{
			return finalizeTest(yyDollar[3].steplist)
		}
	case 3:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line rfml.y:50
		{
			curTest.RFMLID = yyDollar[3].strval
		}
	case 7:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line rfml.y:58
		{
			curTest.Title = yyDollar[3].strval
		}
	case 8:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line rfml.y:59
		{
			curTest.StartURI = yyDollar[3].strval
		}
	case 9:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line rfml.y:60
		{
			curTest.Tags = parseList(yyDollar[3].strval)
		}
	case 10:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line rfml.y:61
		{
			curTest.Browsers = parseList(yyDollar[3].strval)
		}
	case 11:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line rfml.y:62
		{
			curTest.Execute = parseBool(yyDollar[3].strval)
		}
	case 12:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line rfml.y:65
		{
			yyVAL.strval = ""
		}
	case 13:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line rfml.y:66
		{
			yyVAL.strval = yyDollar[1].strval
		}
	case 14:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line rfml.y:69
		{
			yyVAL.steplist = []interface{}{}
		}
	case 15:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line rfml.y:70
		{
			yyVAL.steplist = yyDollar[2].steplist
		}
	case 16:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line rfml.y:71
		{
			yyVAL.steplist = append([]interface{}{yyDollar[1].embedded_test}, yyDollar[2].steplist...)
		}
	case 17:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line rfml.y:72
		{
			yyVAL.steplist = append([]interface{}{yyDollar[1].step}, yyDollar[2].steplist...)
		}
	case 18:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line rfml.y:75
		{
			yyVAL.steplist = []interface{}{}
		}
	case 19:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line rfml.y:77
		{
			yyVAL.embedded_test = rainforest.RFEmbeddedTest{yyDollar[3].strval, yyDollar[1].boolval}
		}
	case 20:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line rfml.y:81
		{
			yyVAL.step = rainforest.RFTestStep{yyDollar[2].strval, yyDollar[3].strval, yyDollar[1].boolval}
		}
	case 21:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line rfml.y:84
		{
			yyVAL.boolval = true
		}
	case 22:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line rfml.y:85
		{
			yyVAL.boolval = parseBool(yyDollar[4].strval)
		}
	case 23:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line rfml.y:88
		{
			yyVAL.strval = yyDollar[1].strval
		}
	case 24:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line rfml.y:91
		{
			yyVAL.strval = yyDollar[1].strval
		}
	}
	goto yystack /* stack new state and value */
}
