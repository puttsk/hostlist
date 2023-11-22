package compress

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
