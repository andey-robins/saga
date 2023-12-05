package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"strconv"
	"time"

	"github.com/andey-robins/magical/genetics"
	"github.com/andey-robins/magical/graph"
	"github.com/andey-robins/magical/sequence"
)

func main() {
	flag.Usage = func() {
		fmt.Println("Run with -help for help information.")
	}

	var graphFile, sequenceFile, out, population, epsilon, mutation string
	var help, verify, memory, execution, minimize, verbose bool
	var seed int
	flag.StringVar(&graphFile, "graph", "", "the path to a graph file")
	flag.StringVar(&sequenceFile, "sequence", "", "the path to a sequence file")
	flag.StringVar(&out, "out", "", "the path to an output file")
	flag.BoolVar(&verify, "verify", false, "use to verify that a sequence is valid for a graph")
	flag.BoolVar(&memory, "memory", false, "use to get the memory utilization of a sequence over a graph")
	flag.BoolVar(&execution, "execution", false, "use to construct a valid execution sequence for the graph")
	flag.BoolVar(&minimize, "minimize", false, "use to minimize the memory utilization of a sequence over a graph")
	flag.BoolVar(&verbose, "verbose", false, "use to display verbose output")
	flag.BoolVar(&help, "help", false, "use to display help text")
	flag.StringVar(&population, "population", "400", "the size of the population to use for genetic algorithms")
	flag.StringVar(&epsilon, "epsilon", "100", "the number of generations to keep running without any improvement")
	flag.StringVar(&mutation, "mutation", "0.2", "the chance of a mutation occuring in a sequence [0.0 - 1.0]")
	flag.IntVar(&seed, "seed", 1, "the seed to use for the random number generator")
	flag.Parse()

	if help {
		pad := func() {
			fmt.Printf("\n\n")
		}

		pad()
		fmt.Println(" Welcome to the MAGIC in-memory memory minimizer!")
		fmt.Println(" This code is licensed under GPLv3")
		pad()
		fmt.Println(" Args:")
		fmt.Println("  -graph:      The path to an input graph file")
		fmt.Println("  -sequence:   The path to an input sequence file")
		fmt.Println("  -out:        The path to an output file. Output will be to STDOUT if\n\t\t none is specified")
		pad()
		fmt.Println(" Flags:")
		fmt.Println("  -verify:     Use to verify that a sequence is valid for a graph.\n\t\tRequires graph and sequence arguments")
		fmt.Println("  -memory:     Use to get the memory utilization of a sequence over a\n\t\t graph. Requires graph and sequence arguments")
		fmt.Println("  -execution:  Use to construct a valid execution sequence for the graph\n\t\t Requires graph and sequence arguments")
		fmt.Println("  -minimize:   Use to minimize the memory utilization of a sequence\n\t\t over a graph. Requires graph and sequence arguments")
		fmt.Println("  -verbose:	Use to display verbose output")
		fmt.Println("  -help:       Display this help text :)")
		pad()
		fmt.Println(" Genetics Arguments:")
		fmt.Println("  -population:  The size of the population to use for genetic algorithms")
		fmt.Println("  -epsilon:     The number of generations to keep running without any improvement")
		fmt.Println("  -mutation:    The chance of a mutation occuring in a sequence [0.0 - 1.0]")
		fmt.Println("  -seed:        The seed to use for the random number generator, default 1, set to 0 for random seed")
		pad()
	}

	if !verbose {
		log.SetOutput(io.Discard)
	}

	// driver flags, run the specified portion of the application
	if verify || memory {
		// validate input parameters
		if graphFile == "" {
			fmt.Println("No graph file specified. Run with -help for help information.")
			return
		}
		if sequenceFile == "" {
			fmt.Println("No sequence file specified. Run with -help for help information.")
			return
		}

		if verify {
			verifyDriver(graphFile, sequenceFile)
		}
		if memory {
			memoryDriver(graphFile, sequenceFile)
		}

	} else if execution || minimize {
		// validate input parameters
		if graphFile == "" {
			fmt.Println("No graph file specified. Run with -help for help information.")
			return
		}
		if out == "" {
			fmt.Println("No output file specified. Run with -help for help information.")
			return
		}

		if execution {
			executionDriver(graphFile, out)
		}
		if minimize {
			population, err := strconv.Atoi(population)
			if err != nil {
				fmt.Println("Invalid population size specified. Run with -help for help information.")
			}
			epsilon, err := strconv.Atoi(epsilon)
			if err != nil {
				fmt.Println("Invalid epsilon specified. Run with -help for help information.")
			}
			mutation, err := strconv.ParseFloat(mutation, 64)
			if err != nil {
				fmt.Println("Invalid mutation chance specified. Run with -help for help information.")
			}

			minimizeDriver(graphFile, out, population, epsilon, seed, mutation)
		}

	} else {
		fmt.Println("No valid flags specified. Run with -help for help information.")
	}
}

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

	m := g.SimulateSequence(s)

	fmt.Printf("Maximum memory footprint: %d\n", m.GetMaxUtilization())
}

func executionDriver(graphFpath, outFpath string) {
	g := graph.LoadGraphFromFile(graphFpath)
	s := g.SynthesizeSequence()
	s.WriteToFile(outFpath)

	fmt.Printf("Non minimized execution synthesized to file %s\n", outFpath)
}

func minimizeDriver(graphFpath, seqFpath string, generation, epsilon, seed int, mutation float64) {
	if seed == 0 {
		fmt.Println("Using random seed")
		seed = int(time.Now().UnixNano())
	}
	g := graph.LoadGraphFromFile(graphFpath)
	p := genetics.NewPopulation(generation, epsilon, mutation, g, seed)

	p.Evolve(g)

	fit, seq := p.GetBest(g)

	fmt.Printf("seed=%d\n", seed)
	fmt.Printf("Best fitness: %d\n", fit)

	seq.WriteToFile(seqFpath)
}
