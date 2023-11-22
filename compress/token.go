package compress

import (
	"fmt"
	"strings"
	"unicode"
)

type TokenType int16

const (
	RootToken   TokenType = iota // Root node
	RuneToken                    // Rune token
	NumberToken                  // Number token
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
	Value string
	Type  TokenType
}

func (t Token) String() string {
	return fmt.Sprintf("{%s:%s}", t.Type, t.Value)
}

// NewToken initializes a Token of type `t` with the providing `args`
func NewToken(t TokenType, args ...string) Token {
	tok := Token{Type: t}
	switch t {
	case RootToken:
		tok.Value = "*"
	case RuneToken:
		tok.Value = args[0][:1]
	case NumberToken:
		tok.Value = args[0]
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
