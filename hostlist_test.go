package hostlist_test

import (
	"reflect"
	"testing"

	"github.com/puttsk/hostlist"
	"github.com/puttsk/hostlist/expand"
)

type ExpandHostlistTestcase struct {
	HostlistExpression string
	ExpectedResult     []string
	ExpectedError      error
}

var ExpandHostlistTestcases = []ExpandHostlistTestcase{
	{
		HostlistExpression: "",
		ExpectedResult:     nil,
		ExpectedError:      expand.ErrEmptyExpression,
	},
	{
		HostlistExpression: "host1",
		ExpectedResult:     []string{"host1"},
		ExpectedError:      nil,
	},
	{
		HostlistExpression: "host[1,2,3]",
		ExpectedResult:     []string{"host1", "host2", "host3"},
		ExpectedError:      nil,
	},
	{
		HostlistExpression: "host[1,,3]",
		ExpectedResult:     []string{"host1", "host", "host3"},
		ExpectedError:      nil,
	},
	{
		HostlistExpression: "host_[1-3]",
		ExpectedResult:     []string{"host_1", "host_2", "host_3"},
		ExpectedError:      nil,
	},

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
		HostlistExpression: "192.168.[0-1].[100-101]",
		ExpectedResult:     []string{"192.168.0.100", "192.168.0.101", "192.168.1.100", "192.168.1.101"},
		ExpectedError:      nil,
	},
	{
		HostlistExpression: "host-[001-004,a],host2-[08-11]",
		ExpectedResult:     []string{"host-001", "host-002", "host-003", "host-004", "host-a", "host2-08", "host2-09", "host2-10", "host2-11"},
		ExpectedError:      nil,
	},
	{
		HostlistExpression: "p[1-2][3-4]s",
		ExpectedResult:     []string{"p13s", "p14s", "p23s", "p24s"},
		ExpectedError:      nil,
	},
	{
		HostlistExpression: "p1,p2,p3",
		ExpectedResult:     []string{"p1", "p2", "p3"},
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
		ExpectedError:      expand.ErrInvalidToken{' ', 7},
	},
	{
		HostlistExpression: "hos]t-[1-4]",
		ExpectedResult:     nil,
		ExpectedError:      expand.ErrInvalidToken{']', 4},
	},
	{
		HostlistExpression: "host-[1-4",
		ExpectedResult:     nil,
		ExpectedError:      expand.ErrExpectedCloseBracket,
	},
	{
		HostlistExpression: "host-[1-4[2-5]]",
		ExpectedResult:     nil,
		ExpectedError:      expand.ErrNestedRangeExpression,
	},
}

// TestExpandHostlist calls hostlist.ExpandHostlist with hostlist expression, checking
// for a valid return value.
func TestExpandHostlist(t *testing.T) {
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
