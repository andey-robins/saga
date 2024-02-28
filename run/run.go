package run

import (
	"github.com/andey-robins/magical/genetics"
	"github.com/andey-robins/magical/graph"
)

type Run struct {
	GA    *genetics.GA
	Graph *graph.Graph
}

func NewRun(pop *genetics.GA, g *graph.Graph) *Run {
	return &Run{pop, g}
}

func (r *Run) Evaluate() {
	r.GA.Evolve(r.Graph)
}
