package graph

// Constructors and methods for the Graph data structure

import (
	"fmt"
	"log"

	"github.com/andey-robins/magical/memory"
	"github.com/andey-robins/magical/sequence"
)

type Graph struct {
	nodes []*Node
	edges int
}

// NewGraph will return a pointer to a new Graph object
// with the given nodes and edges. We count the edges not
// because we couldn't later, but to simplify the printing
// of the graph to a string since the information is provided
// in the input file.
func NewGraph(nodes []*Node, edges int) *Graph {
	return &Graph{nodes, edges}
}

// GetLeaves will return a list of pointers to the leaf nodes in g
// (nodes which are output nodes)
func (g *Graph) GetLeaves() []*Node {
	leaves := make([]*Node, 0)

	for _, node := range g.nodes {
		if !node.HasChildren() {
			leaves = append(leaves, node)
		}
	}

	return leaves
}

// GetRoots will return a list of pointers to the root nodes in g
// (nodes which are input nodes)
func (g *Graph) GetRoots() []*Node {
	roots := make([]*Node, 0)

	for _, node := range g.nodes {
		if !node.HasParents() {
			roots = append(roots, node)
		}
	}

	return roots
}

// GetNodeById will return a pointer to the node in g with the given id
// If no node can be found with the given id, it will return an error.
func (g *Graph) GetNodeById(id int) (*Node, error) {
	for _, node := range g.nodes {
		if node.id == id {
			return node, nil
		}
	}

	return nil, fmt.Errorf("no node with id %d found in graph", id)
}

// IsValidSequence will determine if a given sequence is valid for
// the graph G. It will return true if s can be successfully
// executed for g and false otherwise.
func (g *Graph) IsValidSequence(s *sequence.Sequence) bool {
	// A node in an computational graph can only be executed if all
	// of the prececessor nodes in the graph have been processed.
	// In other words, for a node v in s, all ancestors of v must
	// have already been processed.

	processedNodes := make(map[int]bool)
	for _, node := range g.nodes {
		processedNodes[node.id] = false
	}

	// mark all leaf nodes as processed
	for _, leaf := range g.GetRoots() {
		processedNodes[leaf.id] = true
	}

	log.Println("Initialized processedNodes map and marked leaves processed")

	// helper function which will return true if all ancestors of
	// the parameter have already been processed
	allAncestorsProcessed := func(node *Node, processed *map[int]bool) bool {
		// memoize our own lookup
		if (*processed)[node.id] {
			return true
		}

		// check all parents of the node
		for _, parent := range node.parents {
			if !(*processed)[parent.id] {
				return false
			}
		}
		return true
	}

	// go through each node in the sequence and see if all predecessors
	// have been processed by the time we reach each node
	log.Println("Beginning sequence valditity check")
	for _, nodeId := range s.GetSequence() {
		log.Printf("Next node in sequence: %d\n", nodeId)
		log.Printf("Processed Nodes: %v\n", processedNodes)
		node, err := g.GetNodeById(nodeId)
		if err != nil {
			fmt.Println(err)
			return false
		}

		if !allAncestorsProcessed(node, &processedNodes) {
			return false
		}
		processedNodes[nodeId] = true
		log.Printf("Successfully processed node %d\n", nodeId)
	}

	return true
}

// SimulateSequence takes a sequence s and simulates the execution
// of that sequence over the graph. It returns a pointer to a Memory
// object that had the simulation performed in it. This memory object
// can be used to determine the maximum memory footprint of the sequence
func (g *Graph) SimulateSequence(s *sequence.Sequence) *memory.Memory {
	mem := memory.NewMemory()

	if !g.IsValidSequence(s) {
		return mem
	}

	// load in initial memory to begin simulation
	rootNodes := g.GetRoots()
	for _, node := range rootNodes {
		mem.ProcessNode(node.id, node.GetChildCount(), []int{})
	}

	// step through sequence, processing each node one by one
	for _, nodeId := range s.GetSequence() {
		node, err := g.GetNodeById(nodeId)
		// we check because an invalid sequence is reason to panic
		check(err)

		dependentNodeIds := node.GetParentIds()
		mem.ProcessNode(nodeId, node.GetChildCount(), dependentNodeIds)
	}

	return mem
}

// ToString will return a string representation of g
// that is semantically equivalent to the input string
// used to create the graph object
func (g *Graph) ToString() string {
	roots := g.GetRoots()
	leaves := g.GetLeaves()
	graphString := ""

	// we use this anonymous function to take a list of nodes and
	// get them listed as their ids in a string
	nodeListToIdString := func(nodes []*Node) string {
		nodeString := ""
		for _, node := range nodes {
			nodeString = fmt.Sprintf("%s%d ", nodeString, node.id)
		}
		return nodeString[:len(nodeString)-1]
	}

	// header information
	graphString = fmt.Sprintf("Inputs %d", len(roots))
	graphString = fmt.Sprintf("%s\n%s", graphString, nodeListToIdString(roots))
	graphString = fmt.Sprintf("%s\nOutputs %d", graphString, len(leaves))
	graphString = fmt.Sprintf("%s\n%s", graphString, nodeListToIdString(leaves))
	graphString = fmt.Sprintf("%s\nNodes %d", graphString, len(g.nodes)-len(roots))
	graphString = fmt.Sprintf("%s\nEdges %d", graphString, g.edges)

	// edge information
	for _, node := range g.nodes {
		for _, child := range node.children {
			graphString = fmt.Sprintf("%s\n%d %d", graphString, node.id, child.id)
		}
	}

	return graphString
}
