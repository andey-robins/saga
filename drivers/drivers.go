package drivers

// Drivers are the main entry points for the application
// They validate their arguments and either exit if they are invalid
// or dispatch work into the API

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/andey-robins/magical/checkpoint"
	"github.com/andey-robins/magical/config"
	"github.com/andey-robins/magical/genetics"
	"github.com/andey-robins/magical/graph"
	"github.com/andey-robins/magical/parsers/blif"
	"github.com/andey-robins/magical/sequence"
	"github.com/andey-robins/magical/validation"
)

// VerifyDriver validates a sequence against a graph
func VerifyDriver(graphFpath, seqFpath string) {
	v := validation.NewValidator(validation.Rules{
		validation.ValidateNonEmpty("graph", graphFpath),
		validation.ValidateNonEmpty("sequence", seqFpath),
	})
	v.MustValidate()

	g := graph.LoadGraphFromFile(graphFpath)
	s := sequence.LoadSequenceFromFile(seqFpath)

	isValid := g.IsValidSequence(s)

	if isValid {
		fmt.Println("The execution sequence is valid!")
	} else {
		fmt.Println("The execution sequence is invalid!")
	}
}

// MinimizeDriver uses genetic algorithms to minimize the memory utilization of a sequence over a graph
func MinimizeDriver(graphFpath, seqFpath string, popSize, epsilon, seed int, mutation float64, checkpointFreq int, chkpath string) {
	v := validation.NewValidator(validation.Rules{
		validation.ValidateNonEmpty("graph", graphFpath),
		validation.ValidateNonEmpty("sequence", seqFpath),
		validation.ValidateRangeInt(4, 10_000, popSize),
		validation.ValidateRangeInt(0, 1_000_000, epsilon),
		validation.ValidateRangeFloat(0.0, 1.0, mutation),
		validation.ValidateRangeInt(0, math.MaxInt64, checkpointFreq),
	})
	v.MustValidate()

	if seed == 0 {
		log.Println("Using random seed")
		seed = int(time.Now().UnixNano())
	}

	g := loadGraphByFileType(graphFpath)
	p := genetics.NewGA(popSize, epsilon, mutation, g, seed, checkpointFreq, chkpath)

	p.Evolve(g)

	fit, seq := p.GetBest(g)

	fmt.Printf("seed=%d\n", p.Seed)
	fmt.Printf("Best fitness: %d\n", fit)

	seq.WriteToFile(seqFpath)
}

// MemoryDriver calculates the maximum memory utilization of a sequence over a graph
func MemoryDriver(graphFpath, seqFpath string) {
	v := validation.NewValidator(validation.Rules{
		validation.ValidateNonEmpty("graph", graphFpath),
		validation.ValidateNonEmpty("sequence", seqFpath),
	})
	v.MustValidate()

	g := graph.LoadGraphFromFile(graphFpath)
	s := sequence.LoadSequenceFromFile(seqFpath)

	m, err := g.SimulateSequence(s)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Maximum memory footprint: %d\n", m.GetMaxUtilization())
}

// ResumeDriver resumes a genetic algorithm from a checkpoint
func ResumeDriver(checkpointFpath, graphFile, outFile string) {
	v := validation.NewValidator(validation.Rules{
		validation.ValidateNonEmpty("checkpoint", checkpointFpath),
		validation.ValidateNonEmpty("graph", graphFile),
		validation.ValidateNonEmpty("out", outFile),
	})
	v.MustValidate()

	p := &genetics.GA{}
	checkpoint.Load(checkpointFpath, p)
	p.SynchronizeRNG()

	g := loadGraphByFileType(graphFile)
	p.Evolve(g)

	fit, seq := p.GetBest(g)

	fmt.Printf("seed=%d\n", p.Seed)
	fmt.Printf("Best fitness: %d\n", fit)

	seq.WriteToFile(outFile)
}

func ConfigDriver(configFile string) {
	v := validation.NewValidator(validation.Rules{
		validation.ValidateNonEmpty("config", configFile),
	})
	v.MustValidate()

	cfg := config.ParseConfig(configFile)
	if err := cfg.Validate(); err != nil {
		log.Fatalf("invalid config file, err: %v\n", err)
		return
	}

	GAs := make(map[string]*config.Population)
	for _, pop := range cfg.Populations {
		GAs[pop.Name] = pop
	}

	for _, job := range cfg.Jobs {
		pop, ok := GAs[job.Population]
		if !ok {
			log.Fatalf("invalid population name: %s\n", job.Population)
			return
		}
		g := loadGraphByFileType(job.GraphFile)

		p := genetics.NewGA(pop.Population, pop.Epsilon, pop.MutationRate, g, pop.Seed, pop.CheckpointFreq, pop.CheckpointPath)

		fmt.Println(job.GraphFile)
		p.Evolve(g)

		fit, seq := p.GetBest(g)

		fmt.Printf("seed=%d\n", p.Seed)
		fmt.Printf("Best fitness: %d\n", fit)

		seq.WriteToFile(fmt.Sprintf("%s/%s", job.OutputDir, "final.seq"))
	}
}

func loadGraphByFileType(graphFpath string) *graph.Graph {
	if graphFpath[len(graphFpath)-5:] == ".blif" {
		return blif.LoadBlifAsGraph(graphFpath)
	}

	return graph.LoadGraphFromFile(graphFpath)
}
