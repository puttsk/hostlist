package expand_test

import (
	"reflect"
	"testing"

	"github.com/puttsk/hostlist/expand"
)

type ExpandHostlistTestcase struct {
	HostlistExpression string
	ExpectedResult     []string
	ExpectedError      error
}

var ExpandRangeExpressionTestcases = []ExpandHostlistTestcase{
	{
		HostlistExpression: "",
		ExpectedResult:     nil,
		ExpectedError:      expand.ErrEmptyExpression,
	},
	{
		HostlistExpression: "1",
		ExpectedResult:     []string{"1"},
		ExpectedError:      nil,
	},
	{
		HostlistExpression: "1,a,3",
		ExpectedResult:     []string{"1", "a", "3"},
		ExpectedError:      nil,
	},
	{
		HostlistExpression: "8-13",
		ExpectedResult:     []string{"8", "9", "10", "11", "12", "13"},
		ExpectedError:      nil,
	},
	{
		HostlistExpression: "01-04",
		ExpectedResult:     []string{"01", "02", "03", "04"},
		ExpectedError:      nil,
	},
	{
		HostlistExpression: "009-011,013,0099-0100",
		ExpectedResult:     []string{"009", "010", "011", "013", "0099", "0100"},
		ExpectedError:      nil,
	},
	{
		HostlistExpression: "009-011,,013,0099-0100",
		ExpectedResult:     []string{"009", "010", "011", "", "013", "0099", "0100"},
		ExpectedError:      nil,
	},
	{
		HostlistExpression: "009-011, ,013,0099-0100",
		ExpectedResult:     []string{"009", "010", "011", " ", "013", "0099", "0100"},
		ExpectedError:      nil,
	},
	{
		HostlistExpression: "100-10",
		ExpectedResult:     nil,
		ExpectedError:      expand.ErrInvalidRange,
	},
	{
		HostlistExpression: "009-011,013,0099-0100,3-1",
		ExpectedResult:     nil,
		ExpectedError:      expand.ErrInvalidRange,
	},
}

// TestExpandRangeExpression calls expand.ExpandRangeExpression with range expression, checking
// for a valid return value.
func TestExpandRangeExpression(t *testing.T) {
	for _, c := range ExpandRangeExpressionTestcases {
		t.Logf("Testcase: %s\n", c.HostlistExpression)
		rangeList, err := expand.ExpandRangeExpression(c.HostlistExpression)
		if err != c.ExpectedError {
			t.Fatalf("Invalid error: actual: %s expected: %s", err, c.ExpectedError)
		}
		if !reflect.DeepEqual(rangeList, c.ExpectedResult) {
			t.Fatalf("Invalid ranges: actual: %+v expect: %+v", rangeList, c.ExpectedResult)
		}
	}
}

var ExpandSingleExpressionTestcases = []ExpandHostlistTestcase{
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
		HostlistExpression: "192.168.[0-1].[100-101]",
		ExpectedResult:     []string{"192.168.0.100", "192.168.0.101", "192.168.1.100", "192.168.1.101"},
		ExpectedError:      nil,
	},
	{
		HostlistExpression: "prefix-[005-010]-suffix",
		ExpectedResult:     []string{"prefix-005-suffix", "prefix-006-suffix", "prefix-007-suffix", "prefix-008-suffix", "prefix-009-suffix", "prefix-010-suffix"},
		ExpectedError:      nil,
	},
	{
		HostlistExpression: "host-1,host-2",
		ExpectedResult:     nil,
		ExpectedError:      expand.ErrNotSingleExpression,
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

// TestExpandExpression calls expand.ExpandSingleExpression with hostlist expression, checking
// for a valid return value.
func TestExpandSingleExpression(t *testing.T) {
	for _, c := range ExpandSingleExpressionTestcases {
		t.Logf("Testcase: %s\n", c.HostlistExpression)
		hostnames, err := expand.ExpandSingleExpression(c.HostlistExpression)
		if err != c.ExpectedError {
			t.Fatalf("Invalid error: actual: %s expected: %s", err, c.ExpectedError)
		}
		if !reflect.DeepEqual(hostnames, c.ExpectedResult) {
			t.Fatalf("Invalid hostnames: actual: %+v expect: %+v", hostnames, c.ExpectedResult)
		}
	}
}
