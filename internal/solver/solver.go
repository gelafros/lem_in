package solver

import (
	"fmt"

	"lem-in/internal/parser"
)

func Solve(colony *parser.Colony) ([][]string, error) {
	var paths [][]string
	blocked := make(map[string]bool)

	for {
		// queue will hold paths (slices of room names) we are exploring and has the first arg as colony.StartRoom(the starting point)
		queue := [][]string{{colony.StartRoom}}

		// visited tracks which rooms we have already visited to avoid loops
		visited := make(map[string]bool)
		// Mark all blocked rooms from previous paths as visited so BFS ignores them.
		for room := range blocked {
			visited[room] = true
		}

		// set the first arg of the map visited as colony.StartRoom(the start point)
		visited[colony.StartRoom] = true

		var foundPath []string

		// Keep exploring paths while there are still paths in the queue
		for len(queue) > 0 {
			// Get and remove the first path from the queue
			path := queue[0]
			queue = queue[1:]

			// If we reached the end room, record the path and break BFS
			if path[len(path)-1] == colony.EndRoom {
				foundPath = path
				break
			}

			// Explore neighbors of the last room in the current path
			lastRoom := path[len(path)-1]
			// Loop through all rooms connected to the last room (its neighbors)
			for _, neighbor := range colony.Links[lastRoom] {
				if !visited[neighbor] {
					visited[neighbor] = true

					// Create a new path by copying the current path and appending the neighbor
					newPath := make([]string, len(path))
					copy(newPath, path)
					newPath = append(newPath, neighbor)

					// Add the new path to the back of the queue so its neighbors will be explored in turn (level-by-level)
					queue = append(queue, newPath)
				}
			}
		}

		// If no path was found in this iteration, we have found all possible disjoint paths
		if len(foundPath) == 0 {
			break
		}

		// Save the found path
		paths = append(paths, foundPath)

		// Block the intermediate rooms (excluding start and end)
		// of the foundPath so the next iteration cannot use them.
		for i := 1; i < len(foundPath)-1; i++ {
			blocked[foundPath[i]] = true
		}
	}

	// If no paths at all were found, return an error
	if len(paths) == 0 {
		return nil, fmt.Errorf("ERROR: invalid data format, no path found")
	}
	return paths, nil
}
