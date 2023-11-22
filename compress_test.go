package hostlist_test

import (
	"slices"
	"strings"
	"testing"

	"github.com/puttsk/hostlist"
)

type CompressHostlistTestcase struct {
	Hostlist       []string
	ExpectedResult string
	ExpectedError  error
}

var PrintTreeTestcases = []CompressHostlistTestcase{
	{
		Hostlist:       []string{},
		ExpectedResult: ``,
		ExpectedError:  nil,
	},
	{
		Hostlist:       []string{"aaaaa"},
		ExpectedResult: `{R:a}->{R:a}->{R:a}->{R:a}->{R:a}`,
		ExpectedError:  nil,
	},
	{
		Hostlist: []string{"aa", "ab"},
		ExpectedResult: `{R:a}->{R:a}
       {R:b}`,
		ExpectedError: nil,
	},
	{
		Hostlist: []string{"01", "02", "90", "10"},
		ExpectedResult: `{D:01}
{D:02}
{D:10}
{D:90}`,
		ExpectedError: nil,
	},
	{
		Hostlist: []string{"192.168.1.1", "192.168.1.2", "192.168.1.120"},
		ExpectedResult: `{D:192}->{R:.}->{D:168}->{R:.}->{D:1}->{R:.}->{D:1}
                                              {D:120}
                                              {D:2}`,
		ExpectedError: nil,
	},
	{
		Hostlist: []string{"192.168.1.1", "192.168.1.2", "192.168.2.1", "192.168.2.2"},
		ExpectedResult: `{D:192}->{R:.}->{D:168}->{R:.}->{D:1}->{R:.}->{D:1}
                                              {D:2}
                                {D:2}->{R:.}->{D:1}
                                              {D:2}`,
		ExpectedError: nil,
	},

	{
		Hostlist: []string{"abcd", "abef", "abeg", "xyz", "x1z", "x2z"},
		ExpectedResult: `{R:a}->{R:b}->{R:c}->{R:d}
              {R:e}->{R:f}
                     {R:g}
{R:x}->{D:1}->{R:z}
       {D:2}->{R:z}
       {R:y}->{R:z}`,
		ExpectedError: nil,
	},
	{
		Hostlist: []string{"host-01", "a", "b", "host-03", "host-02", "10-host-120", "11-host-120", "zz-01-a", "yz-01-b", "yz-02-v", "yz-02x"},
		ExpectedResult: `{D:10}->{R:-}->{R:h}->{R:o}->{R:s}->{R:t}->{R:-}->{D:120}
{D:11}->{R:-}->{R:h}->{R:o}->{R:s}->{R:t}->{R:-}->{D:120}
{R:a}
{R:b}
{R:h}->{R:o}->{R:s}->{R:t}->{R:-}->{D:01}
                                   {D:02}
                                   {D:03}
{R:y}->{R:z}->{R:-}->{D:01}->{R:-}->{R:b}
                     {D:02}->{R:-}->{R:v}
                             {R:x}
{R:z}->{R:z}->{R:-}->{D:01}->{R:-}->{R:a}`,
		ExpectedError: nil,
	},
}

// TestPrintTree tests TokenNode.PrintTree by creating a HostlistExpressionTree without compressing
// and check for return value
func TestPrintTree(t *testing.T) {
	for _, c := range PrintTreeTestcases {
		t.Logf("Testcase: %s\n", strings.Join(c.Hostlist, ","))
		slices.Sort(c.Hostlist)

		tree := hostlist.NewHostlistExpressionTree()
		for _, h := range c.Hostlist {
			tree.AddHost(h)
		}

		result := strings.TrimSpace(tree.Root.PrintTree())

		if result != c.ExpectedResult {
			t.Fatalf("Invalid tree: actual:\n%x\nexpect:\n%x\n", result, c.ExpectedResult)
		}
	}
}
