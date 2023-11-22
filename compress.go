package hostlist

import (
	"fmt"
	"strings"
	"unicode"
)

type TokenType int16

const (
	RootToken   TokenType = iota // Root node
	RuneToken                    // Rune token
	NumberToken                  // Number token
)

func (t TokenType) String() string {
	switch t {
	case RootToken:
		return "*"
	case RuneToken:
		return "R"
	case NumberToken:
		return "D"
	}
	return "unknown"
}

type Token struct {
	Value string
	Type  TokenType
}

func (t Token) String() string {
	return fmt.Sprintf("{%s:%s}", t.Type, t.Value)
}

// NewToken initializes a Token of type `t` with the providing `args`
func NewToken(t TokenType, args ...string) Token {
	tok := Token{Type: t}
	switch t {
	case RootToken:
		tok.Value = "*"
	case RuneToken:
		tok.Value = args[0][:1]
	case NumberToken:
		tok.Value = args[0]
	}
	return tok
}

// Tokenize converts string to a list of tokens for hostlist expression.
// Token can be either a rune token, containing single character, or
// a number token, containing an integer.
func Tokenize(str string) []Token {
	result := []Token{}
	hasDigit := false
	builder := strings.Builder{}
	builder.WriteByte(str[0])
	if unicode.IsDigit(rune(str[0])) {
		hasDigit = true
	}

	for _, s := range str[1:] {
		if !unicode.IsDigit(s) {
			tok := RuneToken
			if hasDigit {
				tok = NumberToken
			}
			result = append(result, NewToken(tok, builder.String()))
			hasDigit = false
			builder.Reset()
		} else {
			if !hasDigit {
				tok := RuneToken
				if hasDigit {
					tok = NumberToken
				}
				result = append(result, NewToken(tok, builder.String()))
				hasDigit = true
				builder.Reset()
			}
		}
		builder.WriteRune(s)
	}
	if builder.Len() > 0 {
		tok := RuneToken
		if hasDigit {
			tok = NumberToken
		}
		result = append(result, NewToken(tok, builder.String()))
	}
	return result
}

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

// // [NOT USED]PrintTreeRecursive returns a string representing the structure of TokenNode n and its children.
// // This function traverses the tree recursively in depth-first order
//
//	func (n TokenNode) PrintTreeRecursive(indent int) string {
//		s := make([]string, len(n.Children))
//		tok := n.GetToken().String()
//
//		if len(n.Children) == 0 {
//			return tok
//		}
//
//		for i, c := range n.Children {
//			s[i] = c.PrintTreeRecursive(indent + len(tok+"->"))
//		}
//
//		if n.GetToken().Type == RootToken {
//			return strings.Join(s, "\n"+strings.Repeat(" ", indent))
//		}
//		return tok + "->" + strings.Join(s, "\n"+strings.Repeat(" ", indent))
//	}

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

// HostlistExpressionTree represents a syntax tree of a hostlist expression
type HostlistExpressionTree struct {
	Root   *TokenNode
	Leaves [][]*TokenNode // [level][]Token
}

// NewHostlistExpressionTree initializes and return a new HostlistExpressionTree
func NewHostlistExpressionTree() *HostlistExpressionTree {
	return &HostlistExpressionTree{
		Root:   NewTokenNode(NewToken(RootToken)),
		Leaves: [][]*TokenNode{},
	}
}

// AddHost adds a new host to and restructure the HostlistExpressionTree
func (t *HostlistExpressionTree) AddHost(host string) {
	tokens := Tokenize(host)

	head := t.Root
	for i, token := range tokens {
		found := false
		for j, child := range head.Children {
			if child.GetToken() == token {
				found = true
				head = head.Children[j]
			}
		}

		if !found {
			// Create a new Character node
			node := NewTokenNode(token)
			node.Level = head.Level + 1
			head.Children = append(head.Children, node)
			head = head.Children[len(head.Children)-1]
			if len(t.Leaves) < (head.Level + 1) {
				t.Leaves = append(t.Leaves, []*TokenNode{})
			}
			t.Leaves[i] = append(t.Leaves[i], head)
		}
	}
}

func (t HostlistExpressionTree) String() string {
	return t.Root.PrintTree()
}
