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

func (n *Node) GetId() int {
	return n.id
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

func (n *Node) HasParents() bool {
	return len(n.parents) > 0
}

func (n *Node) HasParent(parentId int) bool {
	for _, parent := range n.parents {
		if parent.id == parentId {
			return true
		}
	}

	return false
}

func (n *Node) HasChildren() bool {
	return len(n.children) > 0
}

func (n *Node) HasChild(childId int) bool {
	for _, child := range n.children {
		if child.id == childId {
			return true
		}
	}

	return false
}
