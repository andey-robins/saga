package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/andey-robins/magical/checkpoint"
	"github.com/andey-robins/magical/genetics"
	"github.com/andey-robins/magical/graph"
	"github.com/andey-robins/magical/parsers/blif"
	"github.com/andey-robins/magical/sequence"
)

func main() {
	flag.Usage = func() {
		fmt.Println("Run with -help for help information.")
	}

	var graphFile, sequenceFile, out, resume string
	var help, verify, memory, evolve, verbose bool
	var seed, population, epsilon int
	var mutation float64
	flag.StringVar(&graphFile, "graph", "", "the path to a graph file")
	flag.StringVar(&sequenceFile, "sequence", "", "the path to a sequence file")
	flag.StringVar(&out, "out", "", "the path to an output file")
	flag.StringVar(&resume, "resume", "", "use to resume from a checkpoint file")

	flag.BoolVar(&verify, "verify", false, "use to verify that a sequence is valid for a graph")
	flag.BoolVar(&memory, "memory", false, "use to get the memory utilization of a sequence over a graph")
	flag.BoolVar(&evolve, "evolve", false, "use to minimize the memory utilization of a sequence over a graph with genetic evolution")
	flag.BoolVar(&verbose, "verbose", false, "use to display verbose output")
	flag.BoolVar(&help, "help", false, "use to display help text")

	flag.IntVar(&population, "pop", 400, "the size of the population to use for genetic algorithms")
	flag.IntVar(&epsilon, "epsilon", 100, "the number of generations to keep running without any improvement")
	flag.Float64Var(&mutation, "mutation", 0.2, "the chance of a mutation occuring in a sequence [0.0 - 1.0]")
	flag.IntVar(&seed, "seed", 1, "the seed to use for the random number generator")

	flag.Parse()

	if help {
		pad := func() {
			fmt.Printf("\n\n")
		}

		pad()
		fmt.Println(" Welcome to the SAGA CLI utility (v0.2.0)")
		fmt.Println(" This code is licensed under GPLv3. Source on GitHub.")
		pad()
		fmt.Println(" Args:")
		fmt.Println("  -graph:      The path to an input graph file")
		fmt.Println("  -sequence:   The path to an input sequence file")
		fmt.Println("  -out:        The path to an output file. Output will be to STDOUT if\n\t\t none is specified")
		fmt.Println("  -resume:     The path to a checkpoint file to resume from. NOTE: This will override any other flags.")
		pad()
		fmt.Println(" Flags:")
		fmt.Println("  -verify:     Use to verify that a sequence is valid for a graph.\n\t\tRequires graph and sequence arguments")
		fmt.Println("  -memory:     Use to get the memory utilization of a sequence over a\n\t\t graph. Requires graph and sequence arguments")
		fmt.Println("  -evolve:     Use to minimize the memory utilization of a sequence\n\t\t over a graph. Requires graph and sequence arguments")
		fmt.Println("  -verbose:	Use to display verbose output")
		fmt.Println("  -help:       Display this help text :)")
		pad()
		fmt.Println(" Genetics Arguments:")
		fmt.Println("  -pop:         The size of the population to use for genetic algorithms (default 400)")
		fmt.Println("  -epsilon:     The number of generations to keep running without any improvement (default 100)")
		fmt.Println("  -mutation:    The chance of a mutation occuring in a sequence [0.0 - 1.0] (default 0.2)")
		fmt.Println("  -seed:        The seed to use for the random number generator, set to 0 for random seed (default 1)")
		pad()
		return
	}

	if !verbose {
		log.SetOutput(io.Discard)
	}

	if graphFile == "" {
		fmt.Println("No graph file specified. Run with -help for help information.")
		return
	}

	if resume != "" {
		if out == "" {
			fmt.Println("No output file specified. Run with -help for help information.")
			return
		}

		fmt.Println("Resuming from checkpoint file:", resume)
		p := &genetics.Population{}
		checkpoint.Load(resume, p)
		p.SynchronizeRNG()

		g := loadGraphByFileType(graphFile)
		p.Evolve(g)

		fit, seq := p.GetBest(g)

		fmt.Printf("seed=%d\n", p.Seed)
		fmt.Printf("Best fitness: %d\n", fit)

		seq.WriteToFile(out)

	} else if verify || memory {

		if sequenceFile == "" {
			fmt.Println("No sequence file specified. Run with -help for help information.")
			return
		}

		if verify {
			verifyDriver(graphFile, sequenceFile)
		} else if memory {
			memoryDriver(graphFile, sequenceFile)
		}

	} else if evolve {

		if out == "" {
			fmt.Println("No output file specified. Run with -help for help information.")
			return
		}

		// parameter validation
		if population < 0 {
			fmt.Println("Population must be greater than 0. Run with -help for help information.")
			return
		}
		if epsilon < 0 {
			fmt.Println("Epsilon must be greater than 0. Run with -help for help information.")
			return
		}
		if mutation < 0.0 || mutation > 1.0 {
			fmt.Println("Mutation must be between 0.0 and 1.0. Run with -help for help information.")
			return
		}

		minimizeDriver(graphFile, out, population, epsilon, seed, mutation)

	} else {
		fmt.Println("No valid flags specified. Run with -help for help information.")
	}
}

// Drivers are the main entry points for the application
// They assume that their inputs have already been validated and
// directly dispatch work into the API as appropriate
func verifyDriver(graphFpath, seqFpath string) {
	g := graph.LoadGraphFromFile(graphFpath)
	s := sequence.LoadSequenceFromFile(seqFpath)

	isValid := g.IsValidSequence(s)

	if isValid {
		fmt.Println("The execution sequence is valid!")
	} else {
		fmt.Println("The execution sequence is invalid!")
	}
}

func memoryDriver(graphFpath, seqFpath string) {
	g := graph.LoadGraphFromFile(graphFpath)
	s := sequence.LoadSequenceFromFile(seqFpath)

	m, err := g.SimulateSequence(s)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Maximum memory footprint: %d\n", m.GetMaxUtilization())
}

func minimizeDriver(graphFpath, seqFpath string, generation, epsilon, seed int, mutation float64) {
	if seed == 0 {
		log.Println("Using random seed")
		seed = int(time.Now().UnixNano())
	}
	g := loadGraphByFileType(graphFpath)
	p := genetics.NewPopulation(generation, epsilon, mutation, g, seed)

	p.Evolve(g)

	fit, seq := p.GetBest(g)

	fmt.Printf("seed=%d\n", p.Seed)
	fmt.Printf("Best fitness: %d\n", fit)

	seq.WriteToFile(seqFpath)
}

func loadGraphByFileType(graphFpath string) *graph.Graph {
	if graphFpath[len(graphFpath)-5:] == ".blif" {
		return blif.LoadBlifAsGraph(graphFpath)
	}

	return graph.LoadGraphFromFile(graphFpath)
}
