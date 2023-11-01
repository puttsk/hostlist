// Package hostlist provides utility function for working with hostlist expression
// Hostlist expression provides a way to define a range of hostnames without an explicit list.
package hostlist

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
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

// IsValidToken checks if rune is a valid for using in hostlist expression
func IsValidToken(r rune) bool {
	return (r >= 'a' && r <= 'z') ||
		(r >= 'A' && r <= 'Z') ||
		(r >= '0' && r <= '9') ||
		r == ',' || r == '[' || r == ']' || r == '-' || r == '_' || r == '.'
}

// SplitExpressions splits hostlist expression separated by ',' and returns a list of hostlist expressions
//
// For example:
//
//	`host-[001-003],node-[3,4,5-10]` will be converted to `["host-[001-003]","node-[3,4,5-10]"]`
func SplitExpressions(hostlist string) ([]string, error) {
	expressions := []string{}

	bracket := 0 // For check bracket level
	var exprBuilder strings.Builder

	// Collect and check hostlist expressions
	for i, s := range hostlist {
		if !(IsValidToken(s)) {
			return nil, ErrInvalidToken{s, i + 1}
		}

		// ',' split expression if bracket level is 0
		if s == ',' && bracket == 0 {
			expressions = append(expressions, exprBuilder.String())

			exprBuilder.Reset() // Reset string builder for next expression
			continue
		}

		// Check bracket for range expression
		if s == '[' {
			// Check if this is nested ranged.
			if bracket > 0 {
				return nil, ErrNestedRangeExpression
			}
			bracket = bracket + 1 // Increase bracket level
		} else if s == ']' {
			// Found ']' without matching bracket
			if bracket == 0 {
				return nil, ErrInvalidToken{s, i + 1}
			}
			bracket = bracket - 1 // Decrease bracket level
		}
		exprBuilder.WriteRune(s)
	}
	// Check if all brackets are closed
	if bracket > 0 {
		return nil, ErrExpectedCloseBracket
	}
	expressions = append(expressions, exprBuilder.String()) // Keep last expression in string builder

	return expressions, nil
}

var rangeExprRegex = regexp.MustCompile(`^(?P<start>\d+)\-(?P<end>\d+)$`)

// ExpandRangeExpression expand a range expression and return an array of hostnames of that expression
//
// For example:
//
//	`001-003` will be converted to `["001","002","003"]`
//	`02-03,a` will be converted to `["02","03","a"]`
func ExpandRangeExpression(expression string) ([]string, error) {
	if expression == "" {
		return nil, ErrEmptyExpression
	}

	rangeList := []string{}

	for _, expr := range strings.Split(expression, ",") {
		m := rangeExprRegex.FindSubmatch([]byte(expr))
		if m == nil {
			rangeList = append(rangeList, expr)
		} else {
			// Extract start and end condition from range expression
			start := string(m[1])
			end := string(m[2])

			// Check if there is leading zeroes
			leadingZeroes := 0
			if start[0] == '0' || end[0] == '0' {
				leadingZeroes = max(len(start), len(end))
			}

			s, err := strconv.ParseInt(start, 10, 64)
			if err != nil {
				return nil, err
			}
			e, err := strconv.ParseInt(end, 10, 64)
			if err != nil {
				return nil, err
			}
			if e < s {
				return nil, ErrInvalidRange
			}

			rangeFormat := fmt.Sprintf("%%0%dd", leadingZeroes)
			for i := s; i <= e; i++ {
				rangeList = append(rangeList, fmt.Sprintf(rangeFormat, i))
			}
		}
	}
	return rangeList, nil
}

// ExpandExpression expand a single hostlist expression and return an array of hostnames of that expression
//
// For example:
//
//	`host-[001-003]` will be converted to `["host-001", "host-002", "host-003"]`
//	`host-1,host-2` will return ErrNotSingleExpression
func ExpandExpression(expression string) ([]string, error) {
	if expression == "" {
		return nil, ErrEmptyExpression
	}

	hosts := []string{}
	rangeExpr := []string{}

	bracket := 0 // For check bracket level
	var prefixBuilder strings.Builder
	var rangeBuilder strings.Builder

	// Collect and check hostlist expressions
	for i, s := range expression {
		if !(IsValidToken(s)) {
			return nil, ErrInvalidToken{s, i + 1}
		}

		// Detect another hostlist expression
		if s == ',' && bracket == 0 {
			return nil, ErrNotSingleExpression
		}

		// Check bracket for range expression
		if s == '[' {
			// Check if this is nested ranged.
			if bracket > 0 {
				return nil, ErrNestedRangeExpression
			}
			bracket = bracket + 1 // Increase bracket level
			continue
		} else if s == ']' {
			// Found ']' without matching bracket
			if bracket == 0 {
				return nil, ErrInvalidToken{']', i + 1}
			}
			bracket = bracket - 1 // Decrease bracket level

			// Range expression is closed, collect range expression
			if bracket == 0 {
				rangeExpr = append(rangeExpr, rangeBuilder.String())
				rangeBuilder.Reset()
				prefixBuilder.WriteString("%s")
				continue
			}
		}
		if bracket == 0 {
			prefixBuilder.WriteRune(s)
		} else {
			rangeBuilder.WriteRune(s)
		}
	}

	// Check if all brackets are closed
	if bracket > 0 {
		return nil, ErrExpectedCloseBracket
	}

	if len(rangeExpr) > 0 {
		rList := [][]interface{}{}
		for _, expr := range rangeExpr {
			r, err := ExpandRangeExpression(expr)
			if err != nil {
				return nil, err
			}
			p := make([]interface{}, len(r))
			for i := range r {
				p[i] = r[i]
			}
			rList = append(rList, p)
		}

		hostFormat := prefixBuilder.String()
		rProduct := CartesianProduct(rList)
		for _, r := range rProduct {
			hosts = append(hosts, fmt.Sprintf(hostFormat, r...))
		}
	} else {
		hosts = append(hosts, prefixBuilder.String())
	}

	return hosts, nil
}

// ExpandHostlist expands hostnames from hostlist expression and return an array of hostnames.
//
// For example:
//
//	`host-[001-003]` will be converted to `["host-001", "host-002", "host-003"]`
func ExpandHostlist(expression string) ([]string, error) {
	hostlist := []string{}

	if expression == "" {
		return nil, ErrEmptyExpression
	}

	expressions, err := SplitExpressions(expression)
	if err != nil {
		return nil, err
	}

	for _, expr := range expressions {
		hosts, err := ExpandExpression(expr)
		if err != nil {
			return nil, err
		}
		hostlist = append(hostlist, hosts...)
	}

	return hostlist, nil
}
