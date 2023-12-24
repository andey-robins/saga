package genetics

import (
	"log"
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
	rng            *rand.Rand
}

func NewPopulation(size, e int, mut float64, graph *graph.Graph, seed int) *Population {
	genes := make([]*Gene, size)

	totalFitness := 0
	bestFitness := math.MaxInt
	bestGene := &Gene{}

	rng := rand.New(rand.NewSource(int64(seed)))

	for i := 0; i < size; i++ {
		seq := graph.SynthesizeRandomValidSequence(rng.Int())
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
		epsilon:        e,
		avgFitness:     float64(totalFitness) / float64(size),
		bestFitness:    bestFitness,
		bestGene:       bestGene,
		generations:    0,
		size:           size,
		mutationChance: mut,
		rng:            rng,
	}
}

func (p *Population) Evolve(g *graph.Graph) {
	bestFitness := p.bestFitness
	roundsWithoutImprovement := 0

	reportEpoch := func() {
		log.Printf("Epoch %d: Best fitness: %d Avg fitness: %v\n", p.generations, p.bestFitness, p.avgFitness)
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
//
// This function is deterministic
func (p *Population) evaluation(g *graph.Graph) {

	var wg sync.WaitGroup
	wg.Add(len(p.genes))
	for _, gene := range p.genes {
		go func(gene *Gene, graph *graph.Graph) {
			if !g.IsValidSequence(gene.sequence) {
				panic("Invalid sequence")
			}
			gene.fitness = graph.SimulateSequence(gene.sequence).GetMaxUtilization()
			wg.Done()
		}(gene, g)
	}
	wg.Wait()
}

// execute will cull the population down to the top 50% or less of genes breaking
// ties randomly
//
// This function is deterministic
func (p *Population) execute() {
	p.calculateStats()

	sort.Slice(p.genes, func(i, j int) bool {
		return p.genes[i].fitness < p.genes[j].fitness
	})

	p.bestFitness = p.genes[0].fitness
	p.bestGene = p.genes[0]

	p.genes = p.genes[:p.size/4]
}

// crossover will select two genes from the population and combine them
// at a random point to create two new genes. this is repeated until we
// have a new population of the same size as the old population
//
// This function uses random numbers, but pulls from p.rng which is seeded
// deterministically and doesn't spawn any go-routines
func (p *Population) crossover(g *graph.Graph) {
	for len(p.genes) < p.size {
		// randomly select two genes and a crossover point
		randGeneOne := p.genes[p.rng.Intn(len(p.genes))]
		randGeneTwo := p.genes[p.rng.Intn(len(p.genes))]
		crossoverPoint := p.rng.Intn(len(randGeneOne.sequence.GetSequence()))

		crossover := func(g1, g2 *Gene, pt int) *Gene {

			genes := make([]int, 0)
			genes = append(genes, g1.sequence.GetSequence()[:pt]...)

			in := func(ll []int, v int) bool {
				for _, vv := range ll {
					if vv == v {
						return true
					}
				}
				return false
			}

			for _, v := range g2.sequence.GetSequence() {
				if !in(genes, v) {
					genes = append(genes, v)
				}
			}

			return &Gene{
				sequence: sequence.NewSequence(genes),
				fitness:  0,
			}
		}

		// create the new genes. we don't evaluate them yet since that'll happen in the next epoch
		newGeneOne := crossover(randGeneOne, randGeneTwo, crossoverPoint)
		newGeneTwo := crossover(randGeneTwo, randGeneOne, crossoverPoint)

		// put them in the population
		p.genes = append(p.genes, newGeneOne, newGeneTwo)
	}

	// if we have too many genes, cull
	if len(p.genes) > p.size {
		p.genes = p.genes[:p.size]
	}
}

// mutate will randomly select a gene from the population and randomly
// swap two of the elements in the sequence sometimes
//
// This function generates random numbers at the time we invoke
// each go-routine and then each routine seeds itself. This prevents
// a race condition preventing determinism that was present in an earlier
// version of this method
func (p *Population) mutate(g *graph.Graph) {
	var wg sync.WaitGroup
	wg.Add(len(p.genes))
	for _, gene := range p.genes {
		go func(gene *Gene, seed int, g *graph.Graph) {
			gene.sequence = g.SmartMutate(gene.sequence, seed)

			wg.Done()
		}(gene, p.rng.Int(), g)
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
