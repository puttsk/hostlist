package compress

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type TokenType int16

const (
	RootToken   TokenType = iota // Root node
	RuneToken                    // Token containing a single character
	NumberToken                  // Token containing a number, including leading zeroes
)

func (t TokenType) String() string {
	switch t {
	case RootToken:
		return "*"
	case RuneToken:
		return "R"
	case NumberToken:
		return "D"
	}
	return "unknown"
}

type Token struct {
	Value      string
	Type       TokenType
	Int        int  // Integer value. For NumberToken only
	ZeroPadded bool // True if the integer is zero padded. For NumberToken only
}

func (t Token) String() string {
	return fmt.Sprintf("{%s:%s}", t.Type, t.Value)
}

// IsNext returns true of Token a is a continuation of Token t, i.e., the Int value of
// Token a is greater than Token t by 1 and both have the same zeroes padding.
func (t Token) IsNext(a Token) bool {
	if t.Type != NumberToken {
		return false
	}

	if a.Type != NumberToken {
		return false
	}

	// Check if a-t == 1
	if a.Int-t.Int != 1 {
		return false
	}

	// The difference is 1 but length of a is less than t, e.g. a=10 and t=009.
	if len(a.Value) < len(t.Value) {
		return false
	}

	// The difference is 1 but length of a is greater than t and a is zero padded, e.g. a=0100 and t=99.
	if len(a.Value) > len(t.Value) && a.ZeroPadded {
		return false
	}

	return true
}

// NewToken initializes a Token of type `t` with the providing `args`
func NewToken(t TokenType, args ...string) Token {
	tok := Token{Type: t}
	switch t {
	case RootToken:
		tok.Value = "*"
	case RuneToken:
		tok.Value = args[0][:1] // Keep only the first character of the first string
	case NumberToken:
		tok.Value = args[0]
		v, _ := strconv.ParseInt(args[0], 10, 0)
		tok.Int = int(v)
		if args[0][0] == '0' {
			tok.ZeroPadded = true
		}
	}
	return tok
}

// Tokenize converts string to a list of tokens for hostlist expression.
// Token can be either a rune token, containing single character, or
// a number token, containing an integer.
func Tokenize(str string) []Token {
	result := []Token{}
	hasDigit := false
	builder := strings.Builder{}
	builder.WriteByte(str[0])
	if unicode.IsDigit(rune(str[0])) {
		hasDigit = true
	}

	for _, s := range str[1:] {
		if !unicode.IsDigit(s) {
			tok := RuneToken
			if hasDigit {
				tok = NumberToken
			}
			result = append(result, NewToken(tok, builder.String()))
			hasDigit = false
			builder.Reset()
		} else {
			if !hasDigit {
				tok := RuneToken
				if hasDigit {
					tok = NumberToken
				}
				result = append(result, NewToken(tok, builder.String()))
				hasDigit = true
				builder.Reset()
			}
		}
		builder.WriteRune(s)
	}
	if builder.Len() > 0 {
		tok := RuneToken
		if hasDigit {
			tok = NumberToken
		}
		result = append(result, NewToken(tok, builder.String()))
	}
	return result
}
