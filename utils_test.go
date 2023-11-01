package hostlist_test

import (
	"fmt"
	"math/rand"
	"reflect"
	"strings"
	"testing"

	"github.com/puttsk/hostlist"
)

type CartesianProductTestcase struct {
	InputList      [][]any
	ExpectedResult [][]any
	ExpectedError  error
}

var CartesianProductTestcases = []CartesianProductTestcase{
	{
		InputList:      [][]any{{1, 2}, {3, 4}},
		ExpectedResult: [][]any{{1, 3}, {1, 4}, {2, 3}, {2, 4}},
	},
	{
		InputList:      [][]any{{1, 2}, {"a", "b"}},
		ExpectedResult: [][]any{{1, "a"}, {1, "b"}, {2, "a"}, {2, "b"}},
	},
	{
		InputList:      [][]any{{1, 2}, {"a", "b"}, {1.1, 5.2}},
		ExpectedResult: [][]any{{1, "a", 1.1}, {1, "a", 5.2}, {1, "b", 1.1}, {1, "b", 5.2}, {2, "a", 1.1}, {2, "a", 5.2}, {2, "b", 1.1}, {2, "b", 5.2}},
	},
}

// TestExpandHosts calls hostlist.ExpandHosts with hostlist expression, checking
// for a valid return value.
func TestCartesianProduct(t *testing.T) {
	for _, c := range CartesianProductTestcases {
		t.Logf("Testcase: %+v\n", c.InputList)
		product := hostlist.CartesianProduct(c.InputList)
		if !reflect.DeepEqual(product, c.ExpectedResult) {
			t.Fatalf("Invalid hostnames: actual: %+v expect: %+v", product, c.ExpectedResult)
		}
	}
}

var CartesianProductBenchmarks = [][][]int{
	{rand.Perm(2), rand.Perm(2)},
	{rand.Perm(100), rand.Perm(100)},
	{rand.Perm(10), rand.Perm(10), rand.Perm(10)},
	{rand.Perm(100), rand.Perm(100), rand.Perm(100)},
	{rand.Perm(10), rand.Perm(10), rand.Perm(10), rand.Perm(10)},
	{rand.Perm(100), rand.Perm(100), rand.Perm(100), rand.Perm(100)},
}

func BenchmarkCartesianProduct(t *testing.B) {
	for _, c := range CartesianProductBenchmarks {
		inputSize := []string{}
		for _, a := range c {
			inputSize = append(inputSize, fmt.Sprint(len(a)))
		}
		t.Run(strings.Join(inputSize, "x"), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				hostlist.CartesianProduct(c)
			}
		})
	}
}
