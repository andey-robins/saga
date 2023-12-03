package graph

import (
	"math/rand"

	"github.com/andey-robins/magical/sequence"
)

// SynthesizeSequence will return a sequence that is semantically valid
// for the graph g. This sequence will be generated using a simple greedy
// breadth first search. The algorithm will start at the root nodes
// and build up. There are no guarantees about memory usage for this method
// but it will guaranteed produce a solution if one exists.
func (g *Graph) SynthesizeSequence() *sequence.Sequence {
	roots := g.GetRoots()
	seq := make([]int, 0)

	in_seq := func(seq []int, id int) bool {
		for _, val := range seq {
			if val == id {
				return true
			}
		}
		return false
	}

	add_parents_to_seq := func(node *Node, seq []int) []int {
		for _, parent := range node.parents {
			if !in_seq(seq, parent.id) {
				seq = append(seq, parent.id)
			}
		}
		return seq
	}

	for _, root := range roots {
		seq = append(seq, root.id)
		seq = add_parents_to_seq(root, seq)
	}

	return sequence.NewSequence(seq)
}

// SynthesizeRandomValidSequence will return a sequence that is semantically valid
// for the graph g. This sequence will be generated using a random breadth first
// search. It searches from the top down by adding nodes to a front when they have
// all of their predecessors processed and randomly sampling from this front to add
// the next node to the sequence. This method is guaranteed to produce a solution as
// long as one exists. It is used to seed our genetic population.
func (g *Graph) SynthesizeRandomValidSequence() *sequence.Sequence {
	processedNodes := g.GetRoots()
	frontNodes := make([]*Node, 0)
	seq := make([]int, 0)

	// a helper function which determines if a given node is in a sequence
	inSeq := func(seq []int, id int) bool {
		for _, val := range seq {
			if val == id {
				return true
			}
		}
		return false
	}

	// a helper function which checks if all of a node's parents are already processed
	readyToBeProcessed := func(node *Node, seq []int) bool {
		for _, parent := range node.parents {
			if !inSeq(seq, parent.id) {
				return false
			}
		}
		return true
	}

	// a helper function which processes a node and adds its children to the frontNodes list
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
		randomIndex := rand.Intn(len(frontNodes))
		nextNodeToProcess := frontNodes[randomIndex]
		frontNodes = append(frontNodes[:randomIndex], frontNodes[randomIndex+1:]...)
		frontNodes, seq = processNode(nextNodeToProcess, seq, frontNodes)
	}

	// remove roots from seq since they don't need to be processed in actuality,
	// but it made the algorithm simpler to include them and then remove them
	// at the end
	removeVal := func(seq []int, val int) []int {
		for i, v := range seq {
			if v == val {
				return append(seq[:i], seq[i+1:]...)
			}
		}
		return seq
	}

	for _, root := range g.GetRoots() {
		seq = removeVal(seq, root.id)
	}

	return sequence.NewSequence(seq)
}
