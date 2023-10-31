package hostlist_test

import (
	"reflect"
	"testing"

	"github.com/puttsk/hostlist"
)

type ExpandHostlistTestcase struct {
	HostlistExpression string
	ExpectedResult     []string
	ExpectedError      error
}

var ExpandHostlistTestcases = []ExpandHostlistTestcase{
	{
		HostlistExpression: "host-[1-4]",
		ExpectedResult:     []string{"host-1", "host-2", "host-3", "host-4"},
		ExpectedError:      nil,
	},
	{
		HostlistExpression: "host-[001-004,a]",
		ExpectedResult:     []string{"host-001", "host-002", "host-003", "host-004", "host-a"},
		ExpectedError:      nil,
	},
	{
		HostlistExpression: "p[1-2][3-4]s",
		ExpectedResult:     []string{"p13s", "p14s", "p23s", "p24s"},
		ExpectedError:      nil,
	},
	{
		HostlistExpression: "p[1-2][3-4]s[01-02]",
		ExpectedResult:     []string{"p13s01", "p13s02", "p14s01", "p14s02", "p23s01", "p23s02", "p24s01", "p24s02"},
		ExpectedError:      nil,
	},
	{
		HostlistExpression: "prefix-[005-010]-suffix",
		ExpectedResult:     []string{"prefix-005-suffix", "prefix-006-suffix", "prefix-007-suffix", "prefix-008-suffix", "prefix-009-suffix", "prefix-010-suffix"},
		ExpectedError:      nil,
	},
	{
		HostlistExpression: "host-[ 001-004,a]",
		ExpectedResult:     nil,
		ExpectedError:      hostlist.ErrInvalidToken{' ', 7},
	},
	{
		HostlistExpression: "hos]t-[1-4]",
		ExpectedResult:     nil,
		ExpectedError:      hostlist.ErrInvalidToken{']', 4},
	},
	{
		HostlistExpression: "host-[1-4",
		ExpectedResult:     nil,
		ExpectedError:      hostlist.ErrExpectedCloseBracket,
	},
	{
		HostlistExpression: "host-[1-4[2-5]]",
		ExpectedResult:     nil,
		ExpectedError:      hostlist.ErrNestedRangeExpression,
	},
}

// TestExpandHosts calls hostlist.ExpandHosts with hostlist expression, checking
// for a valid return value.
func TestExpandHosts(t *testing.T) {
	for _, c := range ExpandHostlistTestcases {
		t.Logf("Testcase: %s\n", c.HostlistExpression)
		hostnames, err := hostlist.ExpandHostlist(c.HostlistExpression)
		if err != c.ExpectedError {
			t.Fatalf("Invalid error: actual: %s expected: %s", err, c.ExpectedError)
		}
		if !reflect.DeepEqual(hostnames, c.ExpectedResult) {
			t.Fatalf("Invalid hostnames: actual: %+v expect: %+v", hostnames, c.ExpectedResult)
		}
	}
}
