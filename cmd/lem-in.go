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
	// Start to read the file to get the colony and check if it is correct file
	colony, err := parser.Parse(os.Args[1])
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	paths, err := solver.Solve(colony)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	simulator.Simulator(colony, paths)
}
