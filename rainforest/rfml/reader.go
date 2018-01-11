package rfml

//go:generate goyacc -o rfml.y.go rfml.y

import (
	"bufio"
	"errors"
	"io"
	"unicode"

	"github.com/rainforestapp/rainforest-cli/rainforest"
)

type Reader struct {
	r          *bufio.Reader
	parseError error

	// vars for internal state tracking
	atbol  bool
	atmeta bool
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

var keywords = map[string]int{
	"string":    _STRING,
	"title":     _TITLE,
	"start_uri": _START_URI,
	"tags":      _TAGS,
	"browsers":  _BROWSERS,
	"redirect":  _REDIRECT,
	"execute":   _EXECUTE,
}

func (r *Reader) Lex(lval *yySymType) int {
	var c rune
	var err error

	if r.parseError != nil {
		return 0
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
		if k, ok := keywords[candidate]; ok {
			return k
		}

		// we're not at a keyword, so just read to EOL as a string
		lval.str = candidate + r.readToEOL()
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

	// As a catch-all, read to the EOL as a string
	lval.str = string(c) + r.readToEOL()
	return _STRING
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

func (r *Reader) readToEOL() string {
	b, err := r.r.ReadBytes('\n')
	if err != nil {
		panic(err)
	}
	// Unread the newline
	err = r.r.UnreadByte()
	if err != nil {
		panic(err)
	}

	return string(b[:len(b)-1])
}

func (r *Reader) Error(e string) {
	r.parseError = errors.New(e)
}
