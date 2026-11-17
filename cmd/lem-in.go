package main

import (
	"fmt"
	"os"

	"lem-in/internal/parser"
	"lem-in/internal/simulator"
	"lem-in/internal/solver"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("ERROR: File argument wasn't provided")
		os.Exit(1)
	}
	// Read and validate the colony file format.
	colony, err := parser.Parse(args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Find the optimal set of vertex-disjoint paths from start to end.
	paths, err := solver.Solve(colony)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Simulate and print turn-by-turn ant movements across the paths.
	if err := simulator.Simulator(colony, paths); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
