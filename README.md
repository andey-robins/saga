# SAGA

Synthesis Augmentation with Genetic Algorithms, SAGA -- is a footprint reduction tool for in-memory computation with MAGIC (memristor aided logic). 

This code is open source and licensed under the GPLv3 license. See `LICENSE` for complete licensing terms.

_Current Version:_ `0.1.2`

[![Go](https://github.com/andey-robins/saga/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/andey-robins/saga/actions/workflows/go.yml)

## Table of Contents
- [SAGA](#saga)
  - [Table of Contents](#table-of-contents)
  - [Getting Started Guide](#getting-started-guide)
  - [Execution](#execution)
    - [Verification Mode](#verification-mode)
    - [Memory Footprint Mode](#memory-footprint-mode)
    - [Minimization Mode](#minimization-mode)
  - [API Usage](#api-usage)
  - [Building](#building)
  - [What is MAGIC?](#what-is-magic)


## Getting Started Guide

A number of operating modes are made available within the SAGA utility. For synthesizing new execution sequences using genetic evolution, see the section on the Minimization Mode below.

Clone this repository and either run the command as detailed below or build the utility into an executable binary using the process in the section titled "Building"

## Execution

The project can be run simply using `go run main.go`. Using that command will provide help information which details CLI arguments and flags. Each operating mode currently supported is enumerated with an example below.

### Verification Mode

This operating mode will verify that a sequence is a semantically correct execution sequence for a given graph. It requires both the sequence and graph arguments.

`go run main.go -verify -graph ./docs/graphs/adder2.graph -sequence ./docs/sequences/adder2.seq`

> ```bash
> The execution sequence is valid!
> ```

### Memory Footprint Mode

This operating mode will report on the memory footprint used by the sequence for the given graph. Similar to *verification mode*, this requires both the sequence and graph as arguments.

`go run main.go -memory -graph ./docs/graphs/adder2.graph -sequence ./docs/sequences/adder2.seq`

> ```bash
> Maximum memory footprint: 9
> ```

### Minimization Mode

This operating mode is the one which applies the genetic algorithms for which this package is named. Additional command line arguments are optional, but allow for configuration of the evolution environment. It requires specifying both a graph and an output file. Another optional argument of `seed` may be specified to create deterministic behavior.

`go run main.go -evolve -graph ./docs/graphs/adder2.graph -pop 300 -epsilon 10 -mutation 0.2 -out ./docs/sequences/synth2.seq -seed 2`

> ```bash
> seed=2
> Best fitness: 7
> ```

## API Usage

A more robust public API is forth-coming in subsequent versions; however, replicating the workflow for minimization and evaluation of a graph can be performed manually through evaluation of the following methods:

1. Create a population with the `genetics.NewPopulation(...)` function.
    - This method takes as input configurations for the population such as mutation rate, maximum population size, etc. and a graph and produces population object which can be evaluated for more efficient solutions.
2. Call the `(* population).Evolve(...)` method on the population.  
   - This takes as an argument the graph. Both references passed to the population (for evolve and for NewPopulation) are immutable references. This will conceivably allow for multiple evolution pipelines to be run over a single graph object in future iterations, but for now is done to parameterize the behavior rather than including the graph as a part of the population.
3. Retrieve the best performance from `(* population).GetBest(...)` which returns both the fitness (memory cost) of the solution and the solution sequence.

## Building

This project is built with go version `1.22.0`.

Run `go build -o saga` to build from source.

## What is MAGIC?

MAGIC, also known as Memristor-Aided Logic, is an emerging computer paradigm that performs computation in-memory rather than on a CPU or other processing device. Values of digital logic can be stored in special memristor-based memory cells. Computation is then performed by issuing read commands to specific addresses and the binary computation is the value available at that address. 

More information on this technique can be found in the foundational paper [here](https://ieeexplore.ieee.org/abstract/document/6895258/).
