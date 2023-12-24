package graph

// Constructors and methods for the Graph data structure

import (
	"fmt"

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

// GetOutputNodes will return a list of pointers to the output nodes in g
func (g *Graph) GetOutputNodes() []*Node {
	outputNodes := make([]*Node, 0)

	for _, node := range g.nodes {
		if !node.HasAnyChildren() {
			outputNodes = append(outputNodes, node)
		}
	}

	return outputNodes
}

// GetInputNodes will return a list of pointers to the input nodes in g
func (g *Graph) GetInputNodes() []*Node {
	inputNodes := make([]*Node, 0)

	for _, node := range g.nodes {
		if !node.HasAnyParents() {
			inputNodes = append(inputNodes, node)
		}
	}

	return inputNodes
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
	for _, leaf := range g.GetInputNodes() {
		processedNodes[leaf.id] = true
	}

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
	for _, nodeId := range s.GetSequence() {
		node, err := g.GetNodeById(nodeId)
		if err != nil {
			fmt.Println(err)
			return false
		}

		if !allAncestorsProcessed(node, &processedNodes) {
			return false
		}
		processedNodes[nodeId] = true
	}

	return true
}

// SimulateSequence takes a sequence s and simulates the execution
// of that sequence over the graph. It returns a pointer to a Memory
// object that had the simulation performed in it. This memory object
// can be used to determine the maximum memory footprint of the sequence.
// If the sequence is invalid, it returns an empty memory object.
func (g *Graph) SimulateSequence(s *sequence.Sequence) (*memory.Memory, error) {
	mem := memory.NewMemory()

	if !g.IsValidSequence(s) {
		return nil, fmt.Errorf("sequence is invalid for graph")
	}

	// load in initial memory to begin simulation
	rootNodes := g.GetInputNodes()
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

	return mem, nil
}

// ToString will return a string representation of g
// that is semantically equivalent to the input string
// used to create the graph object. It may not be identical
// if the input wasn't produced by this tool, but will be
// isomorphic to the input graph.
func (g *Graph) ToString() string {
	roots := g.GetInputNodes()
	leaves := g.GetOutputNodes()
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
