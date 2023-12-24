package graph

import (
	"math/rand"

	"github.com/andey-robins/magical/sequence"
)

// SynthesizeRandomValidSequence will return a sequence that is semantically valid
// for the graph g. This sequence will be generated using a random breadth first
// search. It searches from the top down by adding nodes to a front when they have
// all of their predecessors processed and randomly sampling from this front to add
// the next node to the sequence. This method is guaranteed to produce a solution as
// long as one exists. It is used to seed our genetic population.
func (g *Graph) SynthesizeRandomValidSequence(seed int) *sequence.Sequence {
	rng := rand.New(rand.NewSource(int64(seed)))
	processedNodes := g.GetInputNodes()
	frontNodes := make([]*Node, 0)
	seq := make([]int, 0)

	// determines if a given node `id` is in a sequence `seq`
	inSeq := func(seq []int, id int) bool {
		for _, val := range seq {
			if val == id {
				return true
			}
		}
		return false
	}

	// checks if all of a node `node`s parents are already in the list
	// of processed nodes `seq`
	readyToBeProcessed := func(node *Node, seq []int) bool {
		for _, parent := range node.parents {
			if !inSeq(seq, parent.id) {
				return false
			}
		}
		return true
	}

	// processes a node and adds its children to the front list
	// if they have all of their parents processed
	processNode := func(node *Node, seq []int, front []*Node) ([]*Node, []int) {
		if inSeq(seq, node.id) {
			return front, seq
		}

		seq = append(seq, node.id)
		for _, child := range node.children {
			if readyToBeProcessed(child, seq) {
				front = append(front, child)
			}
		}
		return front, seq
	}

	// process all of the nodes which will begin in memory
	for _, node := range processedNodes {
		seq = append(seq, node.id)
	}

	// identify all nodes one step away from the processed nodes
	for _, node := range processedNodes {
		for _, child := range node.children {
			if readyToBeProcessed(child, seq) {
				frontNodes = append(frontNodes, child)
			}
		}
	}

	for len(frontNodes) > 0 {
		randomIndex := rng.Intn(len(frontNodes))
		nextNodeToProcess := frontNodes[randomIndex]
		frontNodes = append(frontNodes[:randomIndex], frontNodes[randomIndex+1:]...)
		frontNodes, seq = processNode(nextNodeToProcess, seq, frontNodes)
	}

	// removes the value `val` from `seq` it it's there. returns
	// `seq` unchanged if `val` is not in `seq`
	removeVal := func(seq []int, val int) []int {
		for i, v := range seq {
			if v == val {
				return append(seq[:i], seq[i+1:]...)
			}
		}
		return seq
	}

	// remove roots from seq since they don't need to be processed in actuality,
	// but it made the algorithm simpler to include them and then remove them
	// at the end
	for _, root := range g.GetInputNodes() {
		seq = removeVal(seq, root.id)
	}

	return sequence.NewSequence(seq)
}

// SmartMutate will take a sequence `s` and mutate it by swapping a node
// if it is possible to do so without invalidating the sequence. It finds nodes
// which are peers as discussed in the paper and swaps the node at the mutation
// point with one of its peers.
func (g *Graph) SmartMutate(s *sequence.Sequence, seed int) *sequence.Sequence {
	nodes := append(make([]*Node, 0), g.nodes...)

	remove := func(seq []*Node, val *Node) []*Node {
		for i, v := range seq {
			if v.id == val.id {
				return append(seq[:i], seq[i+1:]...)
			}
		}
		return seq
	}

	swapByVal := func(seq []int, val1, val2 int) []int {
		for i, v := range seq {
			if v == val1 {
				seq[i] = val2
			} else if v == val2 {
				seq[i] = val1
			}
		}
		return seq
	}

	// remove all root nodes, since they can't be re-ordered
	for _, r := range g.GetInputNodes() {
		nodes = remove(nodes, r)
	}

	// select a mutation point
	rng := rand.New(rand.NewSource(int64(seed)))
	mutationPoint := rng.Intn(len(s.GetSequence()))
	mutationNodeId := s.GetSequence()[mutationPoint]
	mutationNode, err := g.GetNodeById(mutationNodeId)
	if err != nil {
		panic(err)
	}

	// remove all children nodes of the mutation point from the list of available nodes
	nodes = removeChildrenFromOptions(nodes, mutationNode.children)

	// remove all parent nodes of the mutation point from the list of available nodes
	nodes = removeParentsFromOptions(nodes, mutationNode.parents)

	// select a new node to swap with
	swapCandidateNodeId := nodes[rng.Intn(len(nodes))].id
	newSequence := sequence.NewSequence(swapByVal(s.GetSequence(), mutationNodeId, swapCandidateNodeId))

	for !g.IsValidSequence(newSequence) {
		swapCandidateNodeId = nodes[rng.Intn(len(nodes))].id
		newSequence = sequence.NewSequence(swapByVal(s.GetSequence(), mutationNodeId, swapCandidateNodeId))
	}

	return newSequence
}

func removeChildrenFromOptions(options []*Node, children []*Node) []*Node {
	remove := func(seq []*Node, val *Node) []*Node {
		for i, v := range seq {
			if v.id == val.id {
				return append(seq[:i], seq[i+1:]...)
			}
		}
		return seq
	}

	for _, child := range children {
		options = remove(options, child)
		options = removeChildrenFromOptions(options, child.children)
	}

	return options
}

func removeParentsFromOptions(options []*Node, parents []*Node) []*Node {
	remove := func(seq []*Node, val *Node) []*Node {
		for i, v := range seq {
			if v.id == val.id {
				return append(seq[:i], seq[i+1:]...)
			}
		}
		return seq
	}

	for _, parent := range parents {
		options = remove(options, parent)
		options = removeParentsFromOptions(options, parent.parents)
	}

	return options
}
