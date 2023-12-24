package graph

type Node struct {
	id       int
	parents  []*Node
	children []*Node
}

func NewNode(id int) *Node {
	return &Node{id, make([]*Node, 0), make([]*Node, 0)}
}

func (n *Node) AddParent(parent *Node) {
	n.parents = append(n.parents, parent)
}

func (n *Node) AddChild(child *Node) {
	n.children = append(n.children, child)
}

func (n *Node) GetParentIds() []int {
	ids := make([]int, 0)

	for _, parent := range n.parents {
		ids = append(ids, parent.id)
	}

	return ids
}

func (n *Node) GetChildCount() int {
	return len(n.children)
}

func (n *Node) HasAnyParents() bool {
	return len(n.parents) > 0
}

func (n *Node) HasAnyChildren() bool {
	return len(n.children) > 0
}
