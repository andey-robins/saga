# MAGICAL

Memristor Aided Genetic Intelligent Computation Algorithms -- is a footprint reduction tool for in-memory computation with memristors. 

This code is open source and licensed under the GPLv3 license. See `LICENSE` for complete licensing terms.

_Current Version:_ `0.1.1`

[![Go](https://github.com/andey-robins/magical/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/andey-robins/magical/actions/workflows/go.yml)

## Table of Contents
- [MAGICAL](#magical)
  - [Table of Contents](#table-of-contents)
  - [Getting Started Guide](#getting-started-guide)
  - [Running](#running)
    - [Verification Mode](#verification-mode)
    - [Memory Footprint Mode](#memory-footprint-mode)
    - [Minimization Mode](#minimization-mode)
  - [Building](#building)
  - [Papers](#papers)
    - [MAGICAL](#magical-1)
  - [Changelog](#changelog)
    - [0.1.1](#011)
    - [0.1.0](#010)
  - [Roadmap](#roadmap)
    - [0.1.2](#012)
    - [0.2.0](#020)
    - [0.3.0](#030)
    - [0.4.0](#040)
    - [0.5.0](#050)
    - [1.0.0](#100)


## Getting Started Guide

A number of operating modes are made available within the MAGICAL utility. For synthesizing new execution sequences using genetic evolution, see the section on the Minimization Mode below.

Clone this repository and either run the command as detailed below or build the utility into an executable binary using the process in the section titled "Building"

## Running

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

## Building

This project is built with go version `1.21.5`.

Run `go build -o magical` to build from source.


## Papers

### MAGICAL

Genetic algorithms for more efficient in-memory computation through applied graph analysis is a report which details some initial motivation and refers to the literature which introduced this problem. It presents initial performance of the v0.1.1 version of the software and describes some insights which were used to justify the mapping of this problem to a genetic algorithm.

[Read the paper here.](https://github.com/andey-robins/magical/docs/paper/main.pdf)

## Changelog

This section details changes between revisions of this utility.

### 0.1.1

- Finalize project migration to public namespace. Currently documentation and interaction are sub-optimal, so we'll do one big pass for the 0.1.2 publication before we do the major overhaul on configuration and execution that'll be numbered 0.2.0

### 0.1.0

- Added initial project version

## Roadmap

These are planned features and the releases they are expected with. This should not be seen as a firm commitment but a clear signpost of what is to come provided I remain the sole developer of this tool.

### 0.1.2

- Major updates to documentation
- Simplified CLI flags and arguments

### 0.2.0

- Experimental configuration files
- Checkpointing
- Intermediary saving

### 0.3.0

- External API

### 0.4.0

- Performance analysis
- Additional genetic algorithms

### 0.5.0

- Configurable genetic algorithms and pipelines

### 1.0.0

- Configurable, distributed synthesis