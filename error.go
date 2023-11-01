package hostlist

import (
	"errors"
	"fmt"
)

var ErrEmptyExpression = errors.New("expression cannot be empty string")
var ErrNestedRangeExpression = errors.New("range expression cannot be nested")
var ErrExpectedCloseBracket = errors.New("cannot find matching ']'")
var ErrNotSingleExpression = errors.New("more than single expression detected")
var ErrInvalidRange = errors.New("end value must be greater than start")

type ErrInvalidToken struct {
	Token    rune
	Position int
}

func (e ErrInvalidToken) Error() string {
	return fmt.Sprintf("invalid character '%c' at position %d", e.Token, e.Position)
}
