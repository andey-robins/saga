package main

import (
	"testing"

	"github.com/andey-robins/magical/graph"
	"github.com/andey-robins/magical/sequence"
)

func TestIsValidSequence(t *testing.T) {
	tests := []struct {
		graphString      string
		sequenceString   string
		expectedValidity bool
		maxMemoryUtil    int
	}{
		{
			"Inputs 3\n1 2 3\nOutputs 1\n4\nNodes 4\nEdges 6\n1 5\n2 7\n3 6\n5 7\n6 4\n7 4",
			"Operations 4\n5\n6\n7\n4",
			true,
			4,
		},
	}

	for _, test := range tests {
		g := graph.LoadGraphFromString(test.graphString)
		s := sequence.LoadSequenceFromString(test.sequenceString)

		if g.IsValidSequence(s) != test.expectedValidity {
			t.Errorf("Sequence validity check failed. Expected %t, got %t", test.expectedValidity, g.IsValidSequence(s))
		}

		mem := g.SimulateSequence(s)

		if mem.GetMaxUtilization() != test.maxMemoryUtil {
			t.Errorf("Max memory utilization check failed. Expected %d, got %d", test.maxMemoryUtil, mem.GetMaxUtilization())
		}
	}
}
