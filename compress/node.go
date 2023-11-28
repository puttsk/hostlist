package compress

import (
	"fmt"
	"sort"
	"strings"
)

// TokenNode represents a node in an expression tree
// TokenNode can contain multiple tokens, if those tokens can be represented
// in a range expression
type TokenNode struct {
	Token              Token
	Children           []*TokenNode
	ChildredExpression string // Hostlist expression representing the children node.
	Level              int
}

// NewTokenNode initializes TokenNode with a Token t
func NewTokenNode(t Token) *TokenNode {
	return &TokenNode{
		Token:    t,
		Children: []*TokenNode{},
	}
}

func (n TokenNode) String() string {
	return fmt.Sprint(n.Token)
}

// PrintNode returns a string representing the structure of TokenNode n and its children.
// This function traverses the node iteratively in depth-first order
func (n TokenNode) PrintNode() string {

	if len(n.Children) == 0 {
		return n.String()
	}

	builder := strings.Builder{}

	visited := []string{}

	stack := TokenNodePointerStack{}
	stack.Push(TokenNodePointer{Node: &n})

	for len(stack) > 0 {
		ptr, ok := stack.Pop()
		if !ok {
			break
		}

		prefix := ptr.Prefix

		//fmt.Println(stack)

		// n is child node if not visited
		if !ptr.Visited {

			if ptr.Node.Token.Type != RootToken {
				visited = append(visited, ptr.Node.String())
			}

			// n is not leaf node.
			if len(ptr.Node.Children) > 0 {
				// Traverse the tree from the first node
				for i := len(ptr.Node.Children) - 1; i >= 0; i-- {
					// Put a marker when traverse back to parent
					stack.Push(TokenNodePointer{Node: ptr.Node, Visited: true, Prefix: prefix})
					if ptr.Node.Token.Type != RootToken {
						stack.Push(TokenNodePointer{Node: ptr.Node.Children[i], Prefix: prefix + ptr.Node.String() + "->"})
					} else {
						stack.Push(TokenNodePointer{Node: ptr.Node.Children[i], Prefix: prefix})
					}
				}
			}
		} else {
			if len(visited) > 0 {
				prefixStr := prefix + ptr.Node.String() + "->"
				if ptr.Node.Token.Type == RootToken {
					prefixStr = ""
				}
				visitedStr := strings.Join(visited, "->")

				indent := ""
				// Computing indent for branch
				if !strings.Contains(visitedStr, prefixStr) {
					// Check if there is overlap.
					overlapped := 0
					for i := len(visitedStr); i > 0; i-- {
						if strings.Contains(visitedStr, prefixStr[len(prefixStr)-i:]) {
							overlapped = i
							break
						}
					}
					indent = strings.Repeat(" ", len(prefixStr)-overlapped)
				}
				builder.WriteString(indent + visitedStr + "\n")
			}
			visited = []string{}
		}
	}

	return builder.String()
}

// GetExpression returns a hostlist expression representing the TokenNode.
func (n *TokenNode) GetExpression() string {
	builder := strings.Builder{}

	if n.Token.Type != RootToken {
		builder.WriteString(n.Token.Value)
	}

	childExpressions := []string{}

	// A map of a list of number tokens with the same ChildrenExpression.
	// The ChildrenExpression is used as key to group number token together for creating range expression
	numberMaps := map[string][]*TokenNode{}

	for _, c := range n.Children {
		if c.Token.Type == NumberToken {
			c.GetExpression()
			numberMaps[c.ChildredExpression] = append(numberMaps[c.ChildredExpression], c)
		} else {
			childExpressions = append(childExpressions, c.GetExpression())
		}
	}

	for suffix, numbers := range numberMaps {
		if len(numbers) == 1 {
			childExpressions = append(childExpressions, fmt.Sprintf("%s%s", numbers[0].Token.Value, suffix))
			continue
		}

		// List of number and range expressions
		numberExpr := []string{}

		// Sort TokenNodes based on its integer value
		sort.Slice(numbers, func(i int, j int) bool {
			return numbers[i].Token.Int < numbers[j].Token.Int
		})

		isRangeExpr := false
		lb := numbers[0].Token.Value // lower bound of range expression
		ub := ""                     // upper bound of range expression
		for i := 1; i < len(numbers); i++ {
			// Check if stride is 1 and both has the same zero padding length
			if (numbers[i].Token.Int-numbers[i-1].Token.Int == 1) && numbers[i-1].Token.IsNext(numbers[i].Token) {
				//(len(numbers[i].Token.Value) == len(numbers[i-1].Token.Value)) {
				if !isRangeExpr { // Begin list of stride-1
					isRangeExpr = true
					lb = numbers[i-1].Token.Value
					ub = numbers[i].Token.Value
				} else { // Stride-1 continues, update the upper bound
					ub = numbers[i].Token.Value
				}
			} else {
				if isRangeExpr { // Stride-1 streak is broken. Append the current streak as range expression and start a new one.
					numberExpr = append(numberExpr, fmt.Sprintf("%s-%s", lb, ub))
					lb = numbers[i].Token.Value
					isRangeExpr = false
				} else { // There was no stride-1 streak. Just add the expression to children list
					numberExpr = append(numberExpr, lb)
					lb = numbers[i].Token.Value
					isRangeExpr = false
				}
			}
		}

		// Add the last expression to the list
		if isRangeExpr {
			numberExpr = append(numberExpr, fmt.Sprintf("%s-%s", lb, ub))
		} else {
			numberExpr = append(numberExpr, lb)
		}

		childExpressions = append(childExpressions, fmt.Sprintf("[%s]%s", strings.Join(numberExpr, ","), suffix))
	}

	if len(childExpressions) == 1 {
		builder.WriteString(childExpressions[0])
		n.ChildredExpression = childExpressions[0]
	} else if len(childExpressions) > 1 {
		if n.Token.Type != RootToken {
			builder.WriteString(fmt.Sprintf("[%s]", strings.Join(childExpressions, ",")))
			n.ChildredExpression = fmt.Sprintf("[%s]", strings.Join(childExpressions, ","))
		} else {
			builder.WriteString(strings.Join(childExpressions, ","))
			n.ChildredExpression = strings.Join(childExpressions, ",")
		}
	}

	return builder.String()
}

// TokenPointer represents a pointer for traversing ExpressionTree
type TokenNodePointer struct {
	Node    *TokenNode
	Visited bool
	Prefix  string
}

// TraverseTokenNodeStack provides functions for stack of TraverseTokenNode
type TokenNodePointerStack []TokenNodePointer

func (s *TokenNodePointerStack) Push(n TokenNodePointer) {
	*s = append(*s, n)
}

func (s *TokenNodePointerStack) Pop() (TokenNodePointer, bool) {
	if len(*s) == 0 {
		return TokenNodePointer{}, false
	}

	last := len(*s) - 1
	ret := (*s)[last]
	(*s) = (*s)[:last]

	return ret, true
}
