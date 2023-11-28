package compress

// HostlistExpressionTree represents a syntax tree of a hostlist expression
type HostlistExpressionTree struct {
	Root *TokenNode
	//Leaves [][]*TokenNode // [level][]Token
}

// NewHostlistExpressionTree initializes and return a new HostlistExpressionTree
func NewHostlistExpressionTree() *HostlistExpressionTree {
	return &HostlistExpressionTree{
		Root: NewTokenNode(NewToken(RootToken)),
		//Leaves: [][]*TokenNode{},
	}
}

// AddHost adds a new host to and restructure the HostlistExpressionTree
func (t *HostlistExpressionTree) AddHost(host string) {
	tokens := Tokenize(host)

	head := t.Root
	for _, token := range tokens {
		found := false
		for j, child := range head.Children {
			if child.Token == token {
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
		}
	}
}

func (t HostlistExpressionTree) GetExpression() string {
	return t.Root.GetExpression()
}

func (t HostlistExpressionTree) String() string {
	return t.Root.PrintNode()
}
