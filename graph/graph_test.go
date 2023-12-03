package graph

import "testing"

func TestLoadGraph(t *testing.T) {
	tests := []struct {
		graphString string
	}{
		{
			"Inputs 2\n1 2\nOutputs 2\n3 4\nNodes 2\nEdges 2\n1 3\n2 4",
		},
	}

	for _, test := range tests {
		g := LoadGraphFromString(test.graphString)

		if g.ToString() != test.graphString {
			t.Errorf("Graph loaded incorrectly. Expected:\n%s\n\ngot:\n%s", test.graphString, g.ToString())
		}
	}
}
