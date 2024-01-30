package blif

import (
	"fmt"
	"os"
	"strconv"

	"github.com/andey-robins/magical/graph"
)

func check(e error) {
	if e != nil {
		// panicing is fine here because we can't know how
		// to recover from file errors
		panic(e)
	}
}

// LoadBlifAsGraph will attempt to parse a blif file
// and convert it into a graph. Panics if it can't.
// The parameter fpath is the path of the blif file.
func LoadBlifAsGraph(fpath string) *graph.Graph {
	f, err := os.ReadFile(fpath)
	check(err)

	parser, err := NewBlifParser()
	if err != nil {
		panic(fmt.Sprintf("Error creating blif parser: %s", err))
	}

	blifFile, err := parser.ParseString("", string(f))
	if err != nil {
		panic(fmt.Sprintf("Error parsing blif file: %s", err))
	}

	return blifFile.BlifToGraph()
}

// BlifToGraph will take generate the internal graph object
// from a parsed blif file.
func (b *BlifFile) BlifToGraph() *graph.Graph {
	// Since we already have the logic to convert a graph text string
	// to a graph, we can quickly decode the blif file to that
	// representation and then use the same logic rather than
	// directly building it from the blif file. This could be
	// changed in the future if the conversion becomes a bottleneck
	headerString := ""
	edgesString := ""
	edgeCount := 0

	blifLabelToGraphId := make(map[string]int)
	getGraphId := func(label string) int {
		if id, ok := blifLabelToGraphId[label]; ok {
			return id
		}
		id := len(blifLabelToGraphId) + 1
		blifLabelToGraphId[label] = id
		return id
	}

	headerString += "Inputs " + strconv.Itoa(len(b.Header.InputList.Nodes)) + "\n"
	for _, label := range b.Header.InputList.Nodes {
		headerString += strconv.Itoa(getGraphId(label)) + " "
	}
	headerString = headerString[:len(headerString)-1]
	headerString += "\n"

	headerString += "Outputs " + strconv.Itoa(len(b.Header.OutputList.Nodes)) + "\n"
	for _, label := range b.Header.OutputList.Nodes {
		headerString += strconv.Itoa(getGraphId(label)) + " "
	}
	headerString = headerString[:len(headerString)-1]
	headerString += "\n"

	for _, gate := range b.Gates {
		switch {
		case gate.NotGate != nil:
			edgesString += strconv.Itoa(getGraphId(gate.NotGate.InputName)) + " " + strconv.Itoa(getGraphId(gate.NotGate.OutputName)) + "\n"
			edgeCount++
		case gate.NorGate != nil:
			edgesString += strconv.Itoa(getGraphId(gate.NorGate.InputOneName)) + " " + strconv.Itoa(getGraphId(gate.NorGate.OutputName)) + "\n"
			edgesString += strconv.Itoa(getGraphId(gate.NorGate.InputTwoName)) + " " + strconv.Itoa(getGraphId(gate.NorGate.OutputName)) + "\n"
			edgeCount += 2
		}
	}

	// we calculate nodes and edges last since they go together and
	// we need to write the edge list before counting edges
	headerString += "Nodes " + strconv.Itoa(len(b.Gates)) + "\n"
	headerString += "Edges " + strconv.Itoa(edgeCount) + "\n"

	graphString := headerString + edgesString[:len(edgesString)-1]

	return graph.LoadGraphFromString(graphString)
}
