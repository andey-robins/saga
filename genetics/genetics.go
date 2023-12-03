package genetics

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"sync"

	"github.com/andey-robins/magical/graph"
	"github.com/andey-robins/magical/sequence"
)

type Gene struct {
	sequence *sequence.Sequence
	fitness  int
}

type Population struct {
	genes          []*Gene
	bestFitness    int
	bestGene       *Gene
	avgFitness     float64
	epsilon        int // the number of generations we will continue searching without improvements
	generations    int
	size           int
	mutationChance float64
}

func NewPopulation(size int, graph *graph.Graph) *Population {
	genes := make([]*Gene, size)

	totalFitness := 0
	bestFitness := math.MaxInt
	bestGene := &Gene{}

	for i := 0; i < size; i++ {
		seq := graph.SynthesizeRandomValidSequence()
		fitness := graph.SimulateSequence(seq).GetMaxUtilization()

		genes[i] = &Gene{seq, fitness}

		totalFitness += fitness
		if fitness < bestFitness {
			bestFitness = fitness
			bestGene = genes[i]
		}
	}

	return &Population{
		genes:          genes,
		epsilon:        100,
		avgFitness:     float64(totalFitness) / float64(size),
		bestFitness:    bestFitness,
		bestGene:       bestGene,
		generations:    0,
		size:           size,
		mutationChance: 0.2,
	}
}

func (p *Population) Evolve(g *graph.Graph) {
	bestFitness := p.bestFitness
	roundsWithoutImprovement := 0

	reportEpoch := func() {
		fmt.Printf("Epoch %d: Best fitness: %d Avg fitness: %v\n", p.generations, p.bestFitness, p.avgFitness)
	}

	for roundsWithoutImprovement < p.epsilon {
		p.nextEpoch(g)
		if p.bestFitness < bestFitness {
			roundsWithoutImprovement = 0
			bestFitness = p.bestFitness
		} else {
			roundsWithoutImprovement++
		}
		reportEpoch()
	}
}

func (p *Population) nextEpoch(g *graph.Graph) {
	p.evaluation(g)
	p.execute()
	p.crossover(g)
	p.mutate(g)
	p.generations++
}

// select will evaluate all of the genes in the population and update
// the best gene and best fitness values accordingly. this is parallelized
// using waitgroups since we can evaluate each gene independently
func (p *Population) evaluation(g *graph.Graph) {

	var wg sync.WaitGroup
	wg.Add(p.size)

	for _, gene := range p.genes {
		go func(gene *Gene, graph *graph.Graph) {
			gene.fitness = graph.SimulateSequence(gene.sequence).GetMaxUtilization()
			wg.Done()
		}(gene, g)
	}
	wg.Wait()
}

// execute will cull the population down to the top 50% or less of genes breaking
// ties randomly
func (p *Population) execute() {
	p.calculateStats()

	sort.Slice(p.genes, func(i, j int) bool {
		return p.genes[i].fitness < p.genes[j].fitness
	})

	p.bestFitness = p.genes[0].fitness
	p.bestGene = p.genes[0]

	p.genes = p.genes[:p.size/4-2]
	p.genes = append(p.genes, &Gene{sequence.NewSequence(p.bestGene.sequence.GetSequence()), 0}, &Gene{sequence.NewSequence(p.bestGene.sequence.GetSequence()), 0})
}

// crossover will select two genes from the population and combine them
// at a random point to create two new genes. this is repeated until we
// have a new population of the same size as the old population
func (p *Population) crossover(g *graph.Graph) {
	for len(p.genes) < p.size {
		// randomly select two genes and a crossover point
		randGeneOne := p.genes[rand.Intn(len(p.genes))]
		randGeneTwo := p.genes[rand.Intn(len(p.genes))]
		crossoverPoint := rand.Intn(len(randGeneOne.sequence.GetSequence()))

		// create the new genes. we don't evaluate them yet since that'll happen in the next epoch
		newGeneOne := &Gene{
			sequence: sequence.NewSequence(append(randGeneOne.sequence.GetSequence()[:crossoverPoint], randGeneTwo.sequence.GetSequence()[crossoverPoint:]...)),
			fitness:  0,
		}
		newGeneTwo := &Gene{
			sequence: sequence.NewSequence(append(randGeneTwo.sequence.GetSequence()[:crossoverPoint], randGeneOne.sequence.GetSequence()[crossoverPoint:]...)),
			fitness:  0,
		}

		// put them in the population
		p.genes = append(p.genes, newGeneOne, newGeneTwo)
	}

	// if we have too many genes, kill some randomly
	if len(p.genes) > p.size/2 {
		p.genes = p.genes[:p.size/2]
	}

	for len(p.genes) < p.size {
		p.genes = append(p.genes, &Gene{
			sequence: sequence.NewSequence(g.SynthesizeRandomValidSequence().GetSequence()),
			fitness:  0,
		})
	}
}

// mutate will randomly select a gene from the population and randomly
// swap two of the elements in the sequence sometimes
func (p *Population) mutate(g *graph.Graph) {
	var wg sync.WaitGroup
	wg.Add(p.size)
	for _, gene := range p.genes {
		go func(gene *Gene) {
			original := gene.sequence.GetSequence()
			its := 1
			gene.sequence.Mutate(p.mutationChance)

			for !g.IsValidSequence(gene.sequence) && its < 1000 {
				gene.sequence.Mutate(p.mutationChance)
				its++
			}

			if !g.IsValidSequence(gene.sequence) {
				gene.sequence = sequence.NewSequence(original)
			}
			wg.Done()
		}(gene)
	}
	wg.Wait()
}

func (p *Population) GetBest(g *graph.Graph) (int, *sequence.Sequence) {
	for _, gene := range p.genes {
		gene.fitness = g.SimulateSequence(gene.sequence).GetMaxUtilization()
	}

	sort.Slice(p.genes, func(i, j int) bool {
		return p.genes[i].fitness < p.genes[j].fitness && p.genes[i].fitness != 0
	})

	for _, gene := range p.genes {
		if g.IsValidSequence(gene.sequence) && gene.fitness != 0 {
			p.bestFitness = gene.fitness
			p.bestGene = gene
			return p.bestFitness, p.bestGene.sequence
		}
	}
	return 0, p.bestGene.sequence
}

func (p *Population) calculateStats() {
	totalFitness := 0
	for _, gene := range p.genes {
		if gene.fitness != 0 {
			totalFitness += gene.fitness
		}
	}
	p.avgFitness = float64(totalFitness) / float64(p.size)
}
