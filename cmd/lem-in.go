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
		fmt.Print("ERROR: File argument wasn't provided")
		return
	}
	// Start to read the file to get the colony and check if it is correct format file
	colony, err := parser.Parse(os.Args[1])
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	// Call the solver.Solve to find the shortest and all other possible paths that dont block each other
	paths, err := solver.Solve(colony)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	// Call the simulator.Simulator to guide the ants with max efficiency
	simulator.Simulator(colony, paths)
}
