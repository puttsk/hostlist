package hostlist_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/puttsk/hostlist"
	"github.com/puttsk/hostlist/expand"
)

type ExpandHostlistTestcase struct {
	HostlistExpression string
	ExpectedResult     []string
	ExpectedError      error
}

type CompressHostlistTestcase struct {
	Hostlist       []string
	ExpectedResult string
	ExpectedError  error
}

// Expand hostlist expression then compress again
type ExpandCompressTestcase struct {
	HostlistExpression    string
	ExpectedExpandError   error
	ExpectedCompressError error
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
		ExpectedError:      expand.ErrInvalidToken{Token: ' ', Position: 7},
	},
	{
		HostlistExpression: "hos]t-[1-4]",
		ExpectedResult:     nil,
		ExpectedError:      expand.ErrInvalidToken{Token: ']', Position: 4},
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

var CompressHostlistTestcases = []CompressHostlistTestcase{
	{
		Hostlist:       []string{},
		ExpectedResult: "",
		ExpectedError:  nil,
	},
	{
		Hostlist:       []string{"aaaaa"},
		ExpectedResult: "aaaaa",
		ExpectedError:  nil,
	},
	{
		Hostlist:       []string{"aa", "ab"},
		ExpectedResult: "a[a,b]",
		ExpectedError:  nil,
	},
	{
		Hostlist:       []string{"7", "8", "9", "10", "11"},
		ExpectedResult: "[7-11]",
		ExpectedError:  nil,
	},
	{
		Hostlist:       []string{"a7", "a8", "a9", "a10", "a11"},
		ExpectedResult: "a[7-11]",
		ExpectedError:  nil,
	},
	{
		Hostlist:       []string{"99b", "98b", "100b", "101b"},
		ExpectedResult: "[98-101]b",
		ExpectedError:  nil,
	},
	{
		Hostlist:       []string{"99b", "98b", "100b", "0101b"},
		ExpectedResult: "[98-100,0101]b",
		ExpectedError:  nil,
	},
	{
		Hostlist:       []string{"7a", "7b", "8a", "8b"},
		ExpectedResult: "[7-8][a,b]",
		ExpectedError:  nil,
	},
	{
		Hostlist:       []string{"01", "02", "90", "10"},
		ExpectedResult: "[01-02,10,90]",
		ExpectedError:  nil,
	},
	{
		Hostlist:       []string{"192.168.1.1", "192.168.1.2", "192.168.1.120"},
		ExpectedResult: "192.168.1.[1-2,120]",
		ExpectedError:  nil,
	},
	{
		Hostlist:       []string{"192.168.1.1", "192.168.1.2", "192.168.2.1", "192.168.2.2"},
		ExpectedResult: "192.168.[1-2].[1-2]",
		ExpectedError:  nil,
	},
	{
		Hostlist:       []string{"1.0.3", "1.0.4", "2.0.3", "2.0.4"},
		ExpectedResult: "[1-2].0.[3-4]",
		ExpectedError:  nil,
	},
	{
		Hostlist:       []string{"abcd", "abef", "abeg", "xyz", "x1z", "x2z"},
		ExpectedResult: "ab[cd,e[f,g]],x[yz,[1-2]z]",
		ExpectedError:  nil,
	},
	{
		Hostlist:       []string{"host-01", "a", "b", "host-03", "host-02", "10-host-120", "11-host-120", "zz-01-a", "yz-01-b", "yz-02-v", "yz-02x"},
		ExpectedResult: "a,b,host-[01-03],yz-[01-b,02[-v,x]],zz-01-a,[10-11]-host-120",
		ExpectedError:  nil,
	},
}

var ExpandCompressHostlistTestcases = []ExpandCompressTestcase{
	{
		HostlistExpression:    "",
		ExpectedExpandError:   expand.ErrEmptyExpression,
		ExpectedCompressError: nil,
	},
	{
		HostlistExpression:    "a",
		ExpectedExpandError:   nil,
		ExpectedCompressError: nil,
	},
	{
		HostlistExpression:    "host1",
		ExpectedExpandError:   nil,
		ExpectedCompressError: nil,
	},
	{
		HostlistExpression:    "host-[1-100]",
		ExpectedExpandError:   nil,
		ExpectedCompressError: nil,
	},
	{
		HostlistExpression:    "p[1-2]_[3-4]s",
		ExpectedExpandError:   nil,
		ExpectedCompressError: nil,
	},
	{
		HostlistExpression:    "prefix-[005-010]-suffix",
		ExpectedExpandError:   nil,
		ExpectedCompressError: nil,
	},
	// {
	// 	HostlistExpression:    "host-[a,001-004],other2-[08-11]",
	// 	ExpectedExpandError:   nil,
	// 	ExpectedCompressError: nil,
	// },
}

// TestExpandHostlist calls hostlist.ExpandHostlist with hostlist expression, checking
// for a valid return value.
func TestExpandHostlist(t *testing.T) {
	for _, c := range ExpandHostlistTestcases {
		t.Logf("Testcase: %s\n", c.HostlistExpression)
		hostnames, err := hostlist.Expand(c.HostlistExpression)
		if err != c.ExpectedError {
			t.Fatalf("Invalid error: actual: %s expected: %s", err, c.ExpectedError)
		}
		if !reflect.DeepEqual(hostnames, c.ExpectedResult) {
			t.Fatalf("Invalid hostnames: actual: %+v expect: %+v", hostnames, c.ExpectedResult)
		}
	}
}

// TestCompressHostlist tests TokenNode.CompressHostlist
func TestCompressHostlist(t *testing.T) {
	for _, c := range CompressHostlistTestcases {
		t.Logf("Testcase: %s\n", strings.Join(c.Hostlist, ","))
		expression, err := hostlist.Compress(c.Hostlist)

		if err != c.ExpectedError {
			t.Fatalf("Invalid error: actual: %s expected: %s", err, c.ExpectedError)
		}
		if expression != c.ExpectedResult {
			t.Fatalf("Invalid expression: actual:\n%s\nexpect:\n%s\n", expression, c.ExpectedResult)
		}
	}
}

// TestCompressHostlist tests TokenNode.CompressHostlist
func TestExpandCompressHostlist(t *testing.T) {
	for _, c := range ExpandCompressHostlistTestcases {
		t.Logf("Testcase: %s\n", c.HostlistExpression)
		hosts, err := hostlist.Expand(c.HostlistExpression)
		if err != c.ExpectedExpandError {
			t.Fatalf("Invalid error: actual: %s expected: %s", err, c.ExpectedExpandError)
		}
		expression, err := hostlist.Compress(hosts)
		if err != c.ExpectedCompressError {
			t.Fatalf("Invalid error: actual: %s expected: %s", err, c.ExpectedCompressError)
		}
		if expression != c.HostlistExpression {
			t.Fatalf("Invalid expression: actual:\n%s\nexpect:\n%s\n", expression, c.HostlistExpression)
		}
	}
}
