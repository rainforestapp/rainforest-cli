package rfml

//go:generate goyacc -o rfml.y.go rfml.y

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"regexp"
	"strconv"
	"unicode"

	"github.com/rainforestapp/rainforest-cli/rainforest"
)

type Reader struct {
	r          *bufio.Reader
	parseError error

	// vars for internal state tracking
	atbol         bool
	atmeta        bool
	pseudonewline bool
}

func NewReader(r io.Reader) *Reader {
	return &Reader{
		r:      bufio.NewReader(r),
		atbol:  true,
		atmeta: false,
	}
}

func (r *Reader) ReadAll() (*rainforest.RFTest, error) {
	curTest = &rainforest.RFTest{
		State:   "enabled",
		Execute: true,
		Tags:    []string{},
		Steps:   []interface{}{},
	}
	r.atbol = true
	yyParse(r)
	if r.parseError != nil {
		return nil, r.parseError
	}
	return curTest, nil
}

var headers = map[string]int{
	"title":      _TITLE,
	"start_uri":  _START_URI,
	"tags":       _TAGS,
	"browsers":   _BROWSERS,
	"redirect":   _REDIRECT,
	"execute":    _EXECUTE,
	"site_id":    _SITE_ID,
	"feature_id": _FEATURE_ID,
	"state":      _STATE,
}

func (r *Reader) Lex(lval *yySymType) int {
	var c rune
	var err error

	if r.parseError != nil {
		return 0
	}

	if r.pseudonewline {
		r.pseudonewline = false
		return '\n'
	}

	// Ignore leading whitespace before all tokens and advance to the next
	// siginificant character
	for {
		c, _, err = r.r.ReadRune()
		if err != nil || c != ' ' {
			break
		}
	}

	if err == io.EOF {
		return _EOF
	}
	if err != nil {
		r.parseError = err
		return 0
	}

	// Newlines are their own tokens in most contexts
	if c == '\n' {
		r.atbol = true
		r.atmeta = false
		return int(c)
	}

	// No newline, so next token won't be at BOL for sure
	defer func() { r.atbol = false }()

	if r.atmeta {
		defer func() { r.atmeta = false }()

		// Special consideration for "shebang"
		if c == '!' {
			return int(c)
		}

		r.r.UnreadRune()
		candidate := r.readKeyword()
		if k, ok := headers[candidate]; ok {
			return k
		}

		// we're not at a keyword, so just read to EOL as a string
		lval.strval = candidate + r.readToEOL()
		return _STRING
	}

	if r.atbol {
		if c == '#' {
			r.atmeta = true
			return int(c)
		}
		if c == '-' {
			return int(c)
		}
	}

	if c == ':' {
		return int(c)
	}

	// As a catch-all, to EOL as a value
	val := string(c) + r.readToEOL()
	return r.parseVal(val, lval)
}

func (r *Reader) readKeyword() string {
	var str []rune

	for {
		b, _, err := r.r.ReadRune()
		if err != nil {
			panic(err)
		}

		if unicode.IsLower(b) || b == '_' {
			str = append(str, b)
		} else {
			r.r.UnreadRune()
			break
		}
	}

	return string(str)
}

var bools = map[string]bool{
	"true":  true,
	"false": false,
}
var numRegexp = regexp.MustCompile(`\A[0-9]+\z`)

func (r *Reader) parseVal(val string, lval *yySymType) int {
	if b, ok := bools[val]; ok {
		lval.boolval = b

		return _BOOL
	}

	if numRegexp.MatchString(val) {
		d, err := strconv.Atoi(val)
		if err != nil {
			panic(err)
		}

		lval.intval = d
		return _INTEGER
	}

	lval.strval = val
	return _STRING
}

func (r *Reader) readToEOL() string {
	b, err := r.r.ReadBytes('\n')
	if err != nil && err != io.EOF {
		panic(err)
	}
	if err == io.EOF {
		// Hack: we insert a "pseudo-newline" at EOF if there isn't one.
		r.pseudonewline = true
	} else {
		// Unread the newline
		err = r.r.UnreadByte()
		if err != nil {
			panic(err)
		}
	}

	return string(bytes.TrimSpace(b))
}

func (r *Reader) Error(e string) {
	r.parseError = errors.New(e)
}
