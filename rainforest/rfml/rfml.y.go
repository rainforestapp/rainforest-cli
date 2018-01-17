//line rfml.y:2
package rfml

import __yyfmt__ "fmt"

//line rfml.y:2
import (
	"strings"
)

import "github.com/rainforestapp/rainforest-cli/rainforest"

var curTest *rainforest.RFTest

//line rfml.y:13
type yySymType struct {
	yys           int
	strval        string
	boolval       bool
	intval        int
	steplist      []interface{}
	step          rainforest.RFTestStep
	embedded_test rainforest.RFEmbeddedTest
}

const _STRING = 57346
const _BOOL = 57347
const _INTEGER = 57348
const _TITLE = 57349
const _START_URI = 57350
const _TAGS = 57351
const _BROWSERS = 57352
const _REDIRECT = 57353
const _EXECUTE = 57354
const _FEATURE_ID = 57355
const _SITE_ID = 57356
const _STATE = 57357
const _EOF = 57358

var yyToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"_STRING",
	"_BOOL",
	"_INTEGER",
	"_TITLE",
	"_START_URI",
	"_TAGS",
	"_BROWSERS",
	"_REDIRECT",
	"_EXECUTE",
	"_FEATURE_ID",
	"_SITE_ID",
	"_STATE",
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

//line rfml.y:110

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

//line yacctab:1
var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
	-1, 5,
	16, 19,
	-2, 26,
	-1, 10,
	16, 19,
	-2, 26,
	-1, 11,
	16, 19,
	-2, 26,
	-1, 12,
	16, 19,
	-2, 26,
}

const yyPrivate = 57344

const yyLast = 72

var yyAct = [...]int{

	63, 17, 58, 27, 35, 52, 18, 19, 20, 21,
	45, 22, 25, 23, 24, 44, 26, 43, 42, 60,
	8, 33, 41, 40, 39, 38, 13, 15, 65, 64,
	59, 27, 7, 4, 70, 69, 68, 51, 47, 46,
	53, 54, 55, 56, 26, 5, 61, 6, 62, 29,
	9, 66, 36, 67, 57, 50, 48, 28, 16, 3,
	2, 30, 31, 32, 37, 1, 10, 11, 12, 14,
	49, 34,
}
var yyPact = [...]int{

	15, -1000, 28, 14, 1, 9, -1000, -1, 53, 33,
	9, 9, 9, -1000, 0, 41, 14, -1000, 5, 4,
	3, 2, -2, -3, -5, -10, -1000, 22, 21, -1000,
	-1000, -1000, -1000, 52, 51, 20, -15, -1000, 27, 27,
	27, 27, 49, 13, 27, 13, -1000, -1000, 12, -1000,
	12, -1000, 48, -1000, -1000, -1000, -1000, 19, -1000, -1000,
	18, -1000, -1000, -1000, -1000, -1000, -1000, 17, -1000, -1000,
	-1000,
}
var yyPgo = [...]int{

	0, 71, 70, 69, 50, 68, 67, 66, 1, 2,
	65, 60, 59, 47, 58, 0,
}
var yyR1 = [...]int{

	0, 10, 11, 12, 13, 13, 14, 14, 14, 14,
	14, 14, 14, 14, 14, 8, 8, 9, 9, 4,
	4, 4, 4, 7, 6, 5, 3, 3, 1, 2,
	15, 15,
}
var yyR2 = [...]int{

	0, 4, 2, 4, 0, 3, 1, 3, 3, 3,
	3, 4, 3, 3, 3, 1, 2, 1, 2, 0,
	2, 2, 2, 1, 4, 3, 0, 5, 2, 2,
	1, 1,
}
var yyChk = [...]int{

	-1000, -10, -11, -12, 18, 17, -13, 18, 19, -4,
	-7, -6, -5, 17, -3, 18, -14, -8, 7, 8,
	9, 10, 12, 14, 15, 13, 17, 4, 4, 16,
	-4, -4, -4, 21, -1, 4, 11, -13, 20, 20,
	20, 20, 20, 20, 20, 20, 17, 17, 4, -2,
	4, 17, 20, -8, -8, -8, -8, 5, -9, 17,
	6, -8, -9, -15, 17, 16, -15, 5, 17, 17,
	17,
}
var yyDef = [...]int{

	0, -2, 0, 4, 0, -2, 2, 0, 0, 0,
	-2, -2, -2, 23, 0, 0, 4, 6, 0, 0,
	0, 0, 0, 0, 0, 0, 15, 0, 0, 1,
	20, 21, 22, 0, 0, 0, 0, 5, 0, 0,
	0, 0, 0, 0, 0, 0, 16, 3, 0, 25,
	0, 28, 0, 7, 8, 9, 10, 0, 12, 17,
	0, 13, 14, 24, 30, 31, 29, 0, 11, 18,
	27,
}
var yyTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	17, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 19, 3, 18, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 21, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 20,
}
var yyTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16,
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
		//line rfml.y:50
		{
			return finalizeTest(yyDollar[3].steplist)
		}
	case 3:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line rfml.y:56
		{
			curTest.RFMLID = yyDollar[3].strval
		}
	case 7:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line rfml.y:64
		{
			curTest.Title = yyDollar[3].strval
		}
	case 8:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line rfml.y:65
		{
			curTest.StartURI = yyDollar[3].strval
		}
	case 9:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line rfml.y:66
		{
			curTest.Tags = parseList(yyDollar[3].strval)
		}
	case 10:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line rfml.y:67
		{
			curTest.Browsers = parseList(yyDollar[3].strval)
		}
	case 11:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line rfml.y:68
		{
			curTest.Execute = yyDollar[3].boolval
		}
	case 12:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line rfml.y:69
		{
			curTest.SiteID = yyDollar[3].intval
		}
	case 13:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line rfml.y:70
		{
			curTest.State = yyDollar[3].strval
		}
	case 14:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line rfml.y:71
		{
			curTest.FeatureID = rainforest.FeatureIDInt(yyDollar[3].intval)
		}
	case 15:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line rfml.y:74
		{
			yyVAL.strval = ""
		}
	case 16:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line rfml.y:75
		{
			yyVAL.strval = yyDollar[1].strval
		}
	case 17:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line rfml.y:78
		{
			yyVAL.intval = 0
		}
	case 18:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line rfml.y:79
		{
			yyVAL.intval = yyDollar[1].intval
		}
	case 19:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line rfml.y:82
		{
			yyVAL.steplist = []interface{}{}
		}
	case 20:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line rfml.y:83
		{
			yyVAL.steplist = yyDollar[2].steplist
		}
	case 21:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line rfml.y:84
		{
			yyVAL.steplist = append([]interface{}{yyDollar[1].embedded_test}, yyDollar[2].steplist...)
		}
	case 22:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line rfml.y:85
		{
			yyVAL.steplist = append([]interface{}{yyDollar[1].step}, yyDollar[2].steplist...)
		}
	case 23:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line rfml.y:88
		{
			yyVAL.steplist = []interface{}{}
		}
	case 24:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line rfml.y:90
		{
			yyVAL.embedded_test = rainforest.RFEmbeddedTest{yyDollar[3].strval, yyDollar[1].boolval}
		}
	case 25:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line rfml.y:94
		{
			yyVAL.step = rainforest.RFTestStep{yyDollar[2].strval, yyDollar[3].strval, yyDollar[1].boolval}
		}
	case 26:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line rfml.y:97
		{
			yyVAL.boolval = true
		}
	case 27:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line rfml.y:98
		{
			yyVAL.boolval = yyDollar[4].boolval
		}
	case 28:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line rfml.y:101
		{
			yyVAL.strval = yyDollar[1].strval
		}
	case 29:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line rfml.y:104
		{
			yyVAL.strval = yyDollar[1].strval
		}
	}
	goto yystack /* stack new state and value */
}
