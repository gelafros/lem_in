package parser

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Colony struct {
	AntsCount  int
	Rooms      map[string]Room
	Links      map[string][]string
	StartRoom  string
	EndRoom    string
	InputLines []string // the lines we got from the txt
}

type Room struct {
	Name string
	X, Y int
}

func Parse(filePath string) (*Colony, error) {
	// To store all the rooms as and their links as we parse them, we can initialize a map at the beginning
	rooms := make(map[string]Room)
	links := make(map[string][]string)
	coordinatesRooms := make(map[string]Room)
	linkedRooms := make(map[string]bool)
	// To store the lines we got from the txt we initialize this var
	var InputLines []string

	// Open the file(test.txt)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	// Use defer to close the file before the return
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// Check if the file is empty
	if !scanner.Scan() {
		return nil, fmt.Errorf("ERROR: invalid data format, empty file")
	}
	// Check the firstline (antLine) and use Atoi so ants are a valid number
	antLine := scanner.Text()
	// We append the firstline (antLine) in InputLines to return it later
	InputLines = append(InputLines, antLine)

	ants, err := strconv.Atoi(antLine)
	if err != nil {
		return nil, fmt.Errorf("ERROR: invalid data format, couldnt read the number of ants")
	} else if ants <= 0 {
		return nil, fmt.Errorf("ERROR: invalid data format, invalid number of ants")
	}

	// Initialize var and flag to keep track of whether the next room we parse is the start room or the end room
	var startRoom string
	var endRoom string
	isStart := false
	isEnd := false

	// Start to read all lines
	for scanner.Scan() {
		line := scanner.Text()
		// We append each line if InputLines to return it later
		InputLines = append(InputLines, line)
		// Check that the line has # and not ##(so its a comment) and we skip it (continue)
		if strings.HasPrefix(line, "#") { //&& !strings.HasPrefix(line, "##") {
			if line == "##start" {
				if startRoom != "" {
					return nil, fmt.Errorf("ERROR: invalid data format, ##start duplicate")
				}
				isStart = true
				continue
			}
			if line == "##end" {
				if endRoom != "" {
					return nil, fmt.Errorf("ERROR: invalid data format, ##end duplicate")
				}
				isEnd = true
				continue
			}
			continue
		}

		// Split the line to check if its a room or a link
		roomParts := strings.Fields(line)
		linkParts := strings.Split(line, "-")

		if len(roomParts) == 3 {
			// its a Room

			roomName := roomParts[0]

			if strings.HasPrefix(roomName, "L") || strings.HasPrefix(roomName, "#") {
				return nil, fmt.Errorf("ERROR: invalid data format, roomname is incorrect")
			}

			if isStart {
				startRoom = roomName
				isStart = false
			}
			if isEnd {
				endRoom = roomName
				isEnd = false
			}

			// Keep X and Y as coordinates.
			// {Example} after the line ##start we have a line 1 2 3. If it is a room
			// roomParts(0) = 1 , roomParts(1)= 2, roomParts(2) = 3) etc

			X, err := strconv.Atoi(roomParts[1])
			if err != nil {
				return nil, fmt.Errorf("ERROR: invalid format, couldnt read the X coordinate")
			}
			Y, err := strconv.Atoi(roomParts[2])
			if err != nil {
				return nil, fmt.Errorf("ERROR: invalid format, couldnt read the Y coordinate")
			}
			// Check if there is a duplicate room with the same coordinates
			coordKey := fmt.Sprintf("%d,%d", X, Y)
			_, exists := coordinatesRooms[coordKey]
			if exists {
				return nil, fmt.Errorf("ERROR: invalid data format, duplicate coordinates for room")
			} else {
				coordinatesRooms[coordKey] = Room{
					X: X,
					Y: Y,
				}
			}
			// Check if this room already exists in Room struct
			_, exists = rooms[roomName]
			if exists {
				return nil, fmt.Errorf("ERROR: invalid data format, duplicate room name")
				// If this room doesnt exist add it to Room struct
			} else {
				rooms[roomName] = Room{
					Name: roomName,
					X:    X,
					Y:    Y,
				}
			}

		} else if len(linkParts) == 2 && !strings.Contains(line, " ") {
			// Its a link(tunel)

			// Check if both rooms exist
			room1 := linkParts[0]
			room2 := linkParts[1]
			_, exists1 := rooms[room1]
			_, exists2 := rooms[room2]
			// if they both EXIST we update links(tunels between the rooms)
			if !exists1 || !exists2 {
				return nil, fmt.Errorf("ERROR: invalid data format, link to unknown room")
			} else if room1 == room2 {
				return nil, fmt.Errorf("ERROR: invalid data format, self-linking room")
			} else {
				// Check alphabetically if the link A-B is there and check so we dont add B-A
				var linkKey string
				if room1 < room2 {
					linkKey = room1 + "-" + room2
				} else {
					linkKey = room2 + "-" + room1
				}
				if linkedRooms[linkKey] {
					return nil, fmt.Errorf("ERROR: invalid data format duplicate link")
				}
				linkedRooms[linkKey] = true
				links[room1] = append(links[room1], room2)
				links[room2] = append(links[room2], room1)
			}
		} else {
			return nil, fmt.Errorf("ERROR: invalid data format, malformed line encountered")
		}

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
		InputLines: InputLines,
	}, nil
}
