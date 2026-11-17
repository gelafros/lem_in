package solver

import (
	"fmt"

	"lem-in/internal/parser"
)

// Solve computes a set of vertex-disjoint paths from the colony's start room to its end room
// using the Edmonds-Karp max-flow algorithm on a vertex-split directed graph.
func Solve(colony *parser.Colony) ([][]string, error) {
	var paths [][]string
	// capacities tracks the residual capacity of each directed edge (1 = free, 0 = saturated).
	capacities := make(map[string]int)
	graph := make(map[string][]string)

	for roomName := range colony.Rooms {
		// Start and End rooms can hold multiple ants simultaneously, so their capacities are not restricted.
		if roomName == colony.StartRoom || roomName == colony.EndRoom {
			continue
		}
		// Split each intermediate room into "in" and "out" virtual nodes with capacity 1.
		key := roomName + ":in->" + roomName + ":out"
		capacities[key] = 1

		graph[roomName+":in"] = append(graph[roomName+":in"], roomName+":out")
		graph[roomName+":out"] = append(graph[roomName+":out"], roomName+":in")
	}

	outNode := func(r string) string {
		if r == colony.StartRoom || r == colony.EndRoom {
			return r
		}
		return r + ":out"
	}
	inNode := func(r string) string {
		if r == colony.StartRoom || r == colony.EndRoom {
			return r
		}
		return r + ":in"
	}

	for room1, linkedRooms := range colony.Links {
		for _, room2 := range linkedRooms {
			key := outNode(room1) + "->" + inNode(room2)
			capacities[key] = 1

			graph[outNode(room1)] = append(graph[outNode(room1)], inNode(room2))
			graph[inNode(room2)] = append(graph[inNode(room2)], outNode(room1))
		}
	}

	// Repeatedly search for shortest augmenting paths via BFS.
	for {
		queue := []string{colony.StartRoom}
		visited := make(map[string]bool)
		visited[colony.StartRoom] = true
		parent := make(map[string]string)

		for len(queue) > 0 {
			current := queue[0]
			queue = queue[1:]

			if current == colony.EndRoom {
				break
			}

			for _, neighbor := range graph[current] {
				if capacities[current+"->"+neighbor] > 0 && !visited[neighbor] {
					visited[neighbor] = true
					parent[neighbor] = current
					queue = append(queue, neighbor)
				}
			}
		}

		// Stop once no more augmenting paths exist to EndRoom.
		if !visited[colony.EndRoom] {
			break
		}

		// Backtrack along the parent pointers and update residual capacities.
		curr := colony.EndRoom
		for curr != colony.StartRoom {
			prev := parent[curr]
			capacities[prev+"->"+curr]--
			capacities[curr+"->"+prev]++
			curr = prev
		}
	}

	// Extract the vertex-disjoint paths by tracing saturated edges from StartRoom to EndRoom.
	for _, firstRoom := range colony.Links[colony.StartRoom] {
		if capacities[colony.StartRoom+"->"+inNode(firstRoom)] == 0 {
			path := []string{colony.StartRoom, firstRoom}
			current := firstRoom

			for current != colony.EndRoom {
				for _, nextRoom := range colony.Links[current] {
					if capacities[outNode(current)+"->"+inNode(nextRoom)] == 0 {
						path = append(path, nextRoom)
						current = nextRoom
						break
					}
				}
			}
			paths = append(paths, path)
		}
	}

	if len(paths) == 0 {
		return nil, fmt.Errorf("ERROR: invalid data format, no path found")
	}

	return paths, nil
}
