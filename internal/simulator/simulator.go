package simulator

import (
	"fmt"
	"math"

	"lem-in/internal/parser"
)

// Simulator schedules ants across the provided disjoint paths and prints
// the original colony input followed by the turn-by-turn ant movements.
func Simulator(colony *parser.Colony, paths [][]string) error {
	for _, v := range colony.InputLines {
		fmt.Println(v)
	}
	fmt.Println()

	// antsOnPath stores a list of ant IDs assigned to each path index.
	// For example, antsOnPath[0] will hold all ant IDs using path 0 (e.g. [1, 2, 4])
	antsOnPath := make([][]int, len(paths))

	// Greedily assign each ant to the path that will allow it to reach the end room earliest
	for i := 1; i <= colony.AntsCount; i++ {
		minArrival := math.MaxInt
		bestPath := -1

		// Evaluate each path to see which one gets ant 'i' to the end room first
		for v := 0; v < len(paths); v++ {
			// arrival = path length (in rooms) + number of ants already on this path
			arrival := len(paths[v]) + len(antsOnPath[v])
			if arrival < minArrival {
				minArrival = arrival
				bestPath = v
			}
		}

		// Record that ant 'i' is assigned to bestPath
		antsOnPath[bestPath] = append(antsOnPath[bestPath], i)
	}
	// occupied tracks which ant ID is in which room of each path
	occupied := make([][]int, len(paths))
	for i := 0; i < len(paths); i++ {
		occupied[i] = make([]int, len(paths[i]))
	}

	// releasedCount tracks how many ants have been sent out on each path
	releasedCount := make([]int, len(paths))

	// antsReachedEnd tracks how many ants have reached the ##end room
	antsReachedEnd := 0

	for antsReachedEnd < colony.AntsCount {
		var moves []string
		// Loop through each path to process movement
		for v := 0; v < len(paths); v++ {
			// Part A: Move existing ants on path 'v' forward (in reverse order)
			for j := len(paths[v]) - 2; j >= 1; j-- {
				if occupied[v][j] > 0 {
					antID := occupied[v][j]
					// Move the ant forward
					occupied[v][j+1] = antID
					occupied[v][j] = 0
					// Record this move
					moves = append(moves, fmt.Sprintf("L%d-%s", antID, paths[v][j+1]))
					// If the ant reached the end room, increment our counter
					if j+1 == len(paths[v])-1 {
						antsReachedEnd++
					}
				}
			}
			if releasedCount[v] < len(antsOnPath[v]) {
				// The first room on the path is paths[v][1]
				if 1 == len(paths[v])-1 || occupied[v][1] == 0 {
					antID := antsOnPath[v][releasedCount[v]]
					occupied[v][1] = antID

					moves = append(moves, fmt.Sprintf("L%d-%s", antID, paths[v][1]))
					releasedCount[v]++

					if 1 == len(paths[v])-1 {
						antsReachedEnd++
					}
				}
			}
		}

		// Print all moves made in this round
		if len(moves) > 0 {
			for idx, move := range moves {
				if idx > 0 {
					fmt.Print(" ")
				}
				fmt.Print(move)
			}
			fmt.Println()
		}
	}

	return nil
}
