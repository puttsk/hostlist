package compress

import (
	"fmt"
	"strings"
)

// TokenNode represents a node in an expression tree
// TokenNode can contain multiple tokens, if those tokens can be represented
// in a range expression
type TokenNode struct {
	Tokens   []Token
	Children []*TokenNode
	Level    int
}

// NewTokenNode initializes TokenNode with a Token t
func NewTokenNode(t Token) *TokenNode {
	return &TokenNode{
		Tokens:   []Token{t},
		Children: []*TokenNode{},
	}
}

// GetToken returns the first token of the TokenNode
func (n *TokenNode) GetToken() Token {
	return n.Tokens[0]
}

func (n TokenNode) String() string {
	return fmt.Sprint(n.Tokens)
}

// PrintTree returns a string representing the structure of TokenNode n and its children.
// This function traverses the tree iteratively in depth-first order
func (n TokenNode) PrintTree() string {
	builder := strings.Builder{}

	visited := []string{}

	stack := TraverseTokenNodeStack{}
	stack.Push(TraverseTokenNode{Node: &n})

	for len(stack) > 0 {
		n, ok := stack.Pop()
		if !ok {
			break
		}

		prefix := n.Prefix

		//fmt.Println(stack)

		// n is child node if not visited
		if !n.Visited {

			if n.Node.GetToken().Type != RootToken {
				visited = append(visited, n.Node.GetToken().String())
			}

			// n is not leaf node.
			if len(n.Node.Children) > 0 {
				// Traverse the tree from the first node
				for i := len(n.Node.Children) - 1; i >= 0; i-- {
					// Put a marker when traverse back to parent
					stack.Push(TraverseTokenNode{Node: n.Node, Visited: true, Prefix: prefix})
					if n.Node.GetToken().Type != RootToken {
						stack.Push(TraverseTokenNode{Node: n.Node.Children[i], Prefix: prefix + n.Node.GetToken().String() + "->"})
					} else {
						stack.Push(TraverseTokenNode{Node: n.Node.Children[i], Prefix: prefix})
					}
				}
			}
		} else {
			if len(visited) > 0 {
				prefixStr := prefix + n.Node.GetToken().String() + "->"
				if n.Node.GetToken().Type == RootToken {
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

// TraverseTokenNode represents a TokenNode for tree traveral
type TraverseTokenNode struct {
	Node    *TokenNode
	Visited bool
	Prefix  string
}

// TraverseTokenNodeStack provides functions for stack of TraverseTokenNode
type TraverseTokenNodeStack []TraverseTokenNode

func (s *TraverseTokenNodeStack) Push(n TraverseTokenNode) {
	*s = append(*s, n)
}

func (s *TraverseTokenNodeStack) Pop() (TraverseTokenNode, bool) {
	if len(*s) == 0 {
		return TraverseTokenNode{}, false
	}

	last := len(*s) - 1
	ret := (*s)[last]
	(*s) = (*s)[:last]

	return ret, true
}
