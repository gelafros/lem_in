package simulator

import (
	"fmt"

	"lem-in/internal/parser"
)

func Simulator(colony *parser.Colony, paths [][]string) error {
	for _, v := range colony.InputLines {
		fmt.Println(v)
	}
	assignedCounts := make([]int, len(paths))
	antPaths := make([]int, colony.AntsCount+1) // to save which path each ant takes

	for i := 1; i <= colony.AntsCount; i++ {
		minArrival := 999999
		bestPath := -1

		for v := 0; v < len(paths); v++ {
			arrival := len(paths[v]) + assignedCounts[v]
			if arrival < minArrival {
				minArrival = arrival
				bestPath = v
			}
		}

		// Save the best path for this ant
		antPaths[i] = bestPath
		// Update the count of ants on this path
		assignedCounts[bestPath]++
	}
	return nil
}
