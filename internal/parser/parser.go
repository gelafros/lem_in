package parser

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Colony represents a parsed ant colony topology, including rooms, tunnels, and ant count.
type Colony struct {
	AntsCount  int
	Rooms      map[string]Room
	Links      map[string][]string
	StartRoom  string
	EndRoom    string
	InputLines []string // Raw lines read from the input file
}

// Room represents a named node in the colony with 2D coordinates.
type Room struct {
	Name string
	X, Y int
}

// Parse reads and validates a colony description file at filePath and returns a Colony.
func Parse(filePath string) (*Colony, error) {
	rooms := make(map[string]Room)
	links := make(map[string][]string)
	coordinatesRooms := make(map[string]Room)
	linkedRooms := make(map[string]bool)
	var inputLines []string

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		return nil, fmt.Errorf("ERROR: invalid data format, empty file")
	}

	antLine := scanner.Text()
	inputLines = append(inputLines, antLine)

	ants, err := strconv.Atoi(antLine)
	if err != nil {
		return nil, fmt.Errorf("ERROR: invalid data format, could not read the number of ants")
	} else if ants <= 0 {
		return nil, fmt.Errorf("ERROR: invalid data format, invalid number of ants")
	}

	var startRoom string
	var endRoom string
	isStart := false
	isEnd := false

	for scanner.Scan() {
		line := scanner.Text()
		inputLines = append(inputLines, line)
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// Handle ##start and ##end directives and skip standard comments (#).
		if strings.HasPrefix(line, "#") {
			if line == "##start" {
				if startRoom != "" || isStart || isEnd {
					return nil, fmt.Errorf("ERROR: invalid data format, ##start duplicate or misplaced")
				}
				isStart = true
				continue
			}
			if line == "##end" {
				if endRoom != "" || isEnd || isStart {
					return nil, fmt.Errorf("ERROR: invalid data format, ##end duplicate or misplaced")
				}
				isEnd = true
				continue
			}
			continue
		}

		// Split the line to determine whether it defines a room or a tunnel link.
		roomParts := strings.Fields(line)
		linkParts := strings.Split(line, "-")

		if len(roomParts) == 3 {
			roomName := roomParts[0]

			if strings.HasPrefix(roomName, "L") || strings.HasPrefix(roomName, "#") {
				return nil, fmt.Errorf("ERROR: invalid data format, room name is invalid")
			}

			if isStart {
				startRoom = roomName
				isStart = false
			}
			if isEnd {
				endRoom = roomName
				isEnd = false
			}

			X, err := strconv.Atoi(roomParts[1])
			if err != nil {
				return nil, fmt.Errorf("ERROR: invalid data format, could not read X coordinate")
			}
			Y, err := strconv.Atoi(roomParts[2])
			if err != nil {
				return nil, fmt.Errorf("ERROR: invalid data format, could not read Y coordinate")
			}
			// Reject rooms with duplicate coordinates.
			coordKey := fmt.Sprintf("%d,%d", X, Y)
			if _, exists := coordinatesRooms[coordKey]; exists {
				return nil, fmt.Errorf("ERROR: invalid data format, duplicate coordinates for room")
			}
			coordinatesRooms[coordKey] = Room{
				X: X,
				Y: Y,
			}

			// Reject duplicate room names.
			if _, exists := rooms[roomName]; exists {
				return nil, fmt.Errorf("ERROR: invalid data format, duplicate room name")
			}
			rooms[roomName] = Room{
				Name: roomName,
				X:    X,
				Y:    Y,
			}

		} else if len(linkParts) == 2 && !strings.Contains(line, " ") {
			room1 := linkParts[0]
			room2 := linkParts[1]
			_, exists1 := rooms[room1]
			_, exists2 := rooms[room2]
			if !exists1 || !exists2 {
				return nil, fmt.Errorf("ERROR: invalid data format, link to unknown room")
			} else if room1 == room2 {
				return nil, fmt.Errorf("ERROR: invalid data format, self-linking room")
			} else {
				// Canonicalize the undirected edge key alphabetically to detect duplicate tunnels.
				var linkKey string
				if room1 < room2 {
					linkKey = room1 + "-" + room2
				} else {
					linkKey = room2 + "-" + room1
				}
				if linkedRooms[linkKey] {
					return nil, fmt.Errorf("ERROR: invalid data format, duplicate link")
				}
				linkedRooms[linkKey] = true
				links[room1] = append(links[room1], room2)
				links[room2] = append(links[room2], room1)
			}
		} else {
			return nil, fmt.Errorf("ERROR: invalid data format, malformed line encountered")
		}

	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("ERROR: invalid data format, scanner error: %w", err)
	}
	if startRoom == "" || endRoom == "" {
		return nil, fmt.Errorf("ERROR: invalid data format, invalid colony format")
	}
	if startRoom == endRoom {
		return nil, fmt.Errorf("ERROR: invalid data format, ##start is the same as ##end")
	}

	return &Colony{
		AntsCount:  ants,
		Rooms:      rooms,
		Links:      links,
		StartRoom:  startRoom,
		EndRoom:    endRoom,
		InputLines: inputLines,
	}, nil
}
