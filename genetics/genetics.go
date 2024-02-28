package genetics

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"sort"
	"sync"

	"github.com/andey-robins/magical/checkpoint"
	"github.com/andey-robins/magical/graph"
	"github.com/andey-robins/magical/sequence"
)

type Gene struct {
	Sequence *sequence.Sequence `json:"sequence"`
	Fitness  int                `json:"fitness"`
}

type GA struct {
	Genes          []*Gene `json:"genes"`
	BestFitness    int     `json:"bestFitness"`
	BestGene       *Gene   `json:"bestGene"`
	AvgFitness     float64 `json:"avgFitness"`
	Generations    int     `json:"generations"`
	Size           int     `json:"size"`
	MutationChance float64 `json:"mutationChance"`
	Seed           int     `json:"seed"`
	rng            *rand.Rand

	// the number of generations we will continue searching without improvements
	Epsilon        int    `json:"epsilon"`
	CheckpointFreq int    `json:"checkpointFreq"` // set to 0 to disable checkpoints
	CheckpointPath string `json:"checkpointPath"`
}

// NewGA will create a new population of population size `size` with a mutation
// chance of `mut` and epsilon `e`. It will seed the population with random valid sequences
// and evaluate them.
func NewGA(size, e int, mut float64, graph *graph.Graph, seed int, checkpointFreq int, chkpath string) *GA {
	genes := make([]*Gene, size)

	totalFitness := 0
	bestFitness := math.MaxInt
	bestGene := &Gene{}

	rng := rand.New(rand.NewSource(int64(seed)))

	for i := 0; i < size; i++ {
		seq := graph.SynthesizeRandomValidSequence(rng.Int())
		mem, err := graph.SimulateSequence(seq)
		if err != nil {
			panic(err)
		}
		fitness := mem.GetMaxUtilization()

		genes[i] = &Gene{seq, fitness}

		totalFitness += fitness
		if fitness < bestFitness {
			bestFitness = fitness
			bestGene = genes[i]
		}
	}

	return &GA{
		Genes:          genes,
		Epsilon:        e,
		AvgFitness:     float64(totalFitness) / float64(size),
		BestFitness:    bestFitness,
		BestGene:       bestGene,
		Generations:    0,
		Size:           size,
		MutationChance: mut,
		Seed:           seed,
		rng:            rng,
		CheckpointFreq: checkpointFreq,
		CheckpointPath: chkpath,
	}
}

// Evolve will evolve the population until we have gone `epsilon` generations without
// improving the best fitness
func (p *GA) Evolve(g *graph.Graph) {
	if p.rng == nil {
		p.SynchronizeRNG()
	}

	bestFitness := p.BestFitness
	roundsWithoutImprovement := 0

	reportEpoch := func() {
		log.Printf("Epoch %d: Best fitness: %d Avg fitness: %v\n", p.Generations, p.BestFitness, p.AvgFitness)
	}

	checkpointFilename := func(p *GA) string {
		return fmt.Sprintf("%s/%d.json", p.CheckpointPath, p.Generations)
	}

	if p.CheckpointFreq > 0 {
		os.MkdirAll(p.CheckpointPath, 0755)
		checkpoint.Save(checkpointFilename(p), p)
	}

	for roundsWithoutImprovement < p.Epsilon {
		p.nextEpoch(g)

		if p.CheckpointFreq > 0 && p.Generations%p.CheckpointFreq == 0 {
			checkpoint.Save(checkpointFilename(p), p)
		}

		if p.BestFitness < bestFitness {
			roundsWithoutImprovement = 0
			bestFitness = p.BestFitness
		} else {
			roundsWithoutImprovement++
		}
		reportEpoch()
	}
}

func (p *GA) nextEpoch(g *graph.Graph) {
	p.evaluation(g)
	p.execute()
	p.crossover(g)
	p.mutate(g)
	p.Generations++
}

// select will evaluate all of the genes in the population and update
// the best gene and best fitness values accordingly. this is parallelized
// using waitgroups since we can evaluate each gene independently
//
// This function is deterministic
func (p *GA) evaluation(g *graph.Graph) {

	var wg sync.WaitGroup
	wg.Add(len(p.Genes))
	for _, gene := range p.Genes {
		go func(gene *Gene, graph *graph.Graph) {
			if !g.IsValidSequence(gene.Sequence) {
				panic("Invalid sequence")
			}
			mem, err := graph.SimulateSequence(gene.Sequence)
			if err != nil {
				panic(err)
			}
			gene.Fitness = mem.GetMaxUtilization()
			wg.Done()
		}(gene, g)
	}
	wg.Wait()
}

// execute will cull the population down to the top 50% or less of genes breaking
// ties randomly
//
// This function is deterministic
func (p *GA) execute() {
	p.calculateStats()

	sort.Slice(p.Genes, func(i, j int) bool {
		return p.Genes[i].Fitness < p.Genes[j].Fitness
	})

	p.BestFitness = p.Genes[0].Fitness
	p.BestGene = p.Genes[0]

	p.Genes = p.Genes[:p.Size/4]
}

// crossover will select two genes from the population and combine them
// at a random point to create two new genes. this is repeated until we
// have a new population of the same size as the old population
//
// This function uses random numbers, but pulls from p.rng which is seeded
// deterministically and doesn't spawn any go-routines
func (p *GA) crossover(g *graph.Graph) {
	for len(p.Genes) < p.Size {
		// randomly select two genes and a crossover point
		randGeneOne := p.Genes[p.rng.Intn(len(p.Genes))]
		randGeneTwo := p.Genes[p.rng.Intn(len(p.Genes))]
		crossoverPoint := p.rng.Intn(len(randGeneOne.Sequence.GetSequence()))

		crossover := func(g1, g2 *Gene, pt int) *Gene {

			genes := make([]int, 0)
			genes = append(genes, g1.Sequence.GetSequence()[:pt]...)

			in := func(ll []int, v int) bool {
				for _, vv := range ll {
					if vv == v {
						return true
					}
				}
				return false
			}

			for _, v := range g2.Sequence.GetSequence() {
				if !in(genes, v) {
					genes = append(genes, v)
				}
			}

			return &Gene{
				Sequence: sequence.NewSequence(genes),
				Fitness:  0,
			}
		}

		// create the new genes. we don't evaluate them yet since that'll happen in the next epoch
		newGeneOne := crossover(randGeneOne, randGeneTwo, crossoverPoint)
		newGeneTwo := crossover(randGeneTwo, randGeneOne, crossoverPoint)

		// put them in the population
		p.Genes = append(p.Genes, newGeneOne, newGeneTwo)
	}

	// if we have too many genes, cull
	if len(p.Genes) > p.Size {
		p.Genes = p.Genes[:p.Size]
	}
}

// mutate will randomly select a gene from the population and randomly
// swap two of the elements in the sequence sometimes
//
// This function generates random numbers at the time we invoke
// each go-routine and then each routine seeds itself. This prevents
// a race condition preventing determinism that was present in an earlier
// version of this method
func (p *GA) mutate(g *graph.Graph) {
	var wg sync.WaitGroup
	wg.Add(len(p.Genes))
	for _, gene := range p.Genes {
		go func(gene *Gene, seed int, g *graph.Graph) {
			gene.Sequence = g.SmartMutate(gene.Sequence, seed)

			wg.Done()
		}(gene, p.rng.Int(), g)
	}
	wg.Wait()
}

// GetBest will return the best gene in the population. If there are no valid
// genes in the population, it will return a score of 0. Otherwise, it also returns
// the best fitness and the best sequence
func (p *GA) GetBest(g *graph.Graph) (int, *sequence.Sequence) {
	for _, gene := range p.Genes {
		mem, err := g.SimulateSequence(gene.Sequence)
		if err != nil {
			panic(err)
		}
		gene.Fitness = mem.GetMaxUtilization()
	}

	sort.Slice(p.Genes, func(i, j int) bool {
		return p.Genes[i].Fitness < p.Genes[j].Fitness && p.Genes[i].Fitness != 0
	})

	for _, gene := range p.Genes {
		if g.IsValidSequence(gene.Sequence) && gene.Fitness != 0 {
			p.BestFitness = gene.Fitness
			p.BestGene = gene
			return p.BestFitness, p.BestGene.Sequence
		}
	}
	return 0, p.BestGene.Sequence
}

func (p *GA) calculateStats() {
	totalFitness := 0
	for _, gene := range p.Genes {
		if gene.Fitness != 0 {
			totalFitness += gene.Fitness
		}
	}
	p.AvgFitness = float64(totalFitness) / float64(p.Size)
}

// SynchronizeRNG will reseed the random number generator for the population
// so that we can reproduce the same results. This *MUST* be used when restarting
// from a saved checkpoint
func (p *GA) SynchronizeRNG() {
	p.rng = rand.New(rand.NewSource(int64(p.Seed)))
}
