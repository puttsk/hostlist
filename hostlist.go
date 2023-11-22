// Package hostlist provides utility function for working with hostlist expression
// Hostlist expression provides a way to define a range of hostnames without an explicit list.
package hostlist

import (
	"fmt"
	"slices"

	"github.com/puttsk/hostlist/compress"
	"github.com/puttsk/hostlist/expand"
)

// ExpandHostlist expands hostnames from hostlist expression and return an array of hostnames.
//
// For example:
//
//	`host-[001-003]` will be converted to `["host-001", "host-002", "host-003"]`
func ExpandHostlist(expression string) ([]string, error) {
	hostlist := []string{}

	if expression == "" {
		return nil, expand.ErrEmptyExpression
	}

	expressions, err := expand.SplitExpressions(expression)
	if err != nil {
		return nil, err
	}

	for _, expr := range expressions {
		hosts, err := expand.ExpandSingleExpression(expr)
		if err != nil {
			return nil, err
		}
		hostlist = append(hostlist, hosts...)
	}

	return hostlist, nil
}

// CompressHostlist return hostlist expression from a list of host
func CompressHostlist(hosts []string) (string, error) {
	tree := compress.NewHostlistExpressionTree()

	slices.Sort(hosts)

	for _, h := range hosts {
		tree.AddHost(h)
	}

	return fmt.Sprint(tree), nil
}
