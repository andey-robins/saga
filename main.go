package main

import (
	"flag"
	"fmt"
	"io"
	"log"

	"github.com/andey-robins/magical/drivers"
)

func main() {
	flag.Usage = func() {
		fmt.Println("Run with -help for help information.")
	}

	var graphFile, sequenceFile, out, resume, chkpath string
	var help, verify, memory, evolve, verbose bool
	var seed, population, epsilon, checkpointFreq int
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

	flag.IntVar(&checkpointFreq, "chkfreq", 1, "the number of generations between checkpoints")
	flag.StringVar(&chkpath, "chkpath", "./checkpoints", "the path to a directory to save checkpoint files to")

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
		fmt.Println("  -chkfreq:    The number of generations between checkpoints (default 1)")
		fmt.Println("  -chkpath:    The path to a directory to save checkpoints to (default ./checkpoints)")
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

	if resume != "" {
		drivers.ResumeDriver(resume, graphFile, out)

	} else if verify {
		drivers.VerifyDriver(graphFile, sequenceFile)

	} else if memory {
		drivers.MemoryDriver(graphFile, sequenceFile)

	} else if evolve {
		drivers.MinimizeDriver(graphFile, out, population, epsilon, seed, mutation, checkpointFreq, chkpath)

	} else {
		fmt.Println("No valid flags specified. Run with -help for help information.")
	}
}
