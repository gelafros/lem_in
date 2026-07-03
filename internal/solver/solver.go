package solver

import (
	"fmt"

	"lem-in/internal/parser"
)

func Solve(colony *parser.Colony) ([][]string, error) {
	// This is the var that we will return
	var paths [][]string
	// This is the var to keep track of the capacity of each tunnel(with 0 meaning used and 1 meaning free) so we know which room is already used
	capacities := make(map[string]int)
	graph := make(map[string][]string)
	// We loop through all the rooms in Colony struct
	for roomName := range colony.Rooms {
		// Start and End rooms can hold multiple ants so we dont limit their capacity
		if roomName == colony.StartRoom || roomName == colony.EndRoom {
			continue
		}
		// Create a virtual tunnel for the room e.g.  "RoomA:in->RoomA:out"
		key := roomName + ":in->" + roomName + ":out"
		// Set the capacity to 1 meaning only 1 ant can use this room at a time
		capacities[key] = 1

		graph[roomName+":in"] = append(graph[roomName+":in"], roomName+":out")
		graph[roomName+":out"] = append(graph[roomName+":out"], roomName+":in")
	}

	// Helper to get the correct exit node name for a room
	outNode := func(r string) string {
		if r == colony.StartRoom || r == colony.EndRoom {
			return r
		}
		return r + ":out"
	}
	// Helper to get the correct entrance node name for a room
	inNode := func(r string) string {
		if r == colony.StartRoom || r == colony.EndRoom {
			return r
		}
		return r + ":in"
	}
	// Loop through all the connections (tunnels) between rooms
	for room1, linkedRooms := range colony.Links {
		for _, room2 := range linkedRooms {
			// Connect the exit of room1 to the entrance of room2 with a capacity of 1
			key := outNode(room1) + "->" + inNode(room2)
			capacities[key] = 1

			graph[outNode(room1)] = append(graph[outNode(room1)], inNode(room2))
			graph[inNode(room2)] = append(graph[inNode(room2)], outNode(room1))
		}
	}

	// This infinite loop keeps searching for paths until no more can be found
	for {
		// queue holds the virtual nodes(rooms) we need to visit next in our BFS
		queue := []string{colony.StartRoom}

		// visited keeps track of nodes we have seen so we dont go in circles
		visited := make(map[string]bool)
		visited[colony.StartRoom] = true

		// parent remembers where we came from, so we can trace the path backward later
		parent := make(map[string]string)

		// Breadth-First Search (BFS) loop
		for len(queue) > 0 {
			// Take the first node out of the queue
			current := queue[0]
			queue = queue[1:]

			// If we reached the end room, we successfully found a path so we stop searching
			if current == colony.EndRoom {
				break
			}

			// Look at all the neighbors connected to our current node in the virtual graph
			for _, neighbor := range graph[current] {
				// We can only move to the neighbor if the tunnel has capacity 1
				// AND we havent visited it yet in this search
				if capacities[current+"->"+neighbor] > 0 && !visited[neighbor] {
					visited[neighbor] = true        // Mark as visited
					parent[neighbor] = current      // Remember that we came from 'current'
					queue = append(queue, neighbor) // Add it to the queue to explore later
				}
			}
		}

		// If the BFS finished but never reached the EndRoom, it means the graph is
		// fully blocked and no more paths exist. We break the infinite loop to finish
		if !visited[colony.EndRoom] {
			break
		}

		// We reached the EndRoom and now we backtrack using the 'parent' map to update capacities
		curr := colony.EndRoom
		for curr != colony.StartRoom {
			prev := parent[curr]

			// We consume the forward capacity (it becomes 0, so it's blocked for the next search)
			capacities[prev+"->"+curr]--
			// We INCREASE the reverse capacity (it becomes 1). This is the "undo" magic
			// that allows a future search to push flow backward and reroute this path!
			capacities[curr+"->"+prev]++

			// Move one step backward toward the StartRoom
			curr = prev
		}
	} // <-- The infinite loop properly closes here

	// --- Phase 2: Extracting the actual paths from the graph ---

	// We look at all rooms connected directly to the StartRoom
	for _, firstRoom := range colony.Links[colony.StartRoom] {
		// If the capacity from Start to this room is 0, it means flow (an ant) went this way!
		if capacities[colony.StartRoom+"->"+inNode(firstRoom)] == 0 {

			// Start building our path slice
			path := []string{colony.StartRoom, firstRoom}
			current := firstRoom

			// Keep moving forward until we reach the EndRoom
			for current != colony.EndRoom {
				// Check all neighbors of the current room
				for _, nextRoom := range colony.Links[current] {
					// If the capacity from the exit of current room to entrance of next room is 0,
					// it means the flow continued this way!
					if capacities[outNode(current)+"->"+inNode(nextRoom)] == 0 {
						path = append(path, nextRoom) // Add to our path
						current = nextRoom            // Move forward
						break                         // Found the correct next room, stop checking neighbors
					}
				}
			}
			// We finished tracing a full path to the end, save it to our list of paths!
			paths = append(paths, path)
		}
	}

	// If we couldn't extract any paths at all, return an error
	if len(paths) == 0 {
		return nil, fmt.Errorf("ERROR: invalid data format, no path found")
	}

	// Return the final list of completely not conncecting and  optimal paths!
	return paths, nil
}
