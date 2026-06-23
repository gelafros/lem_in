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
		return nil, err
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
		if strings.HasPrefix(line, "#") && !strings.HasPrefix(line, "##") {
			continue
		}
		if line == "##start" {
			isStart = true
			continue
		}
		if line == "##end" {
			isEnd = true
			continue
		}
		// Split the line to check if its a room or a link
		roomParts := strings.Fields(line)
		linkParts := strings.Split(line, "-")

		if len(roomParts) == 3 {
			// its a Room

			roomName := roomParts[0]

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
				return nil, err
			}
			Y, err := strconv.Atoi(roomParts[2])
			if err != nil {
				return nil, err
			}
			// Check if this room already exists in Room struct
			_, exists := rooms[roomName]
			if exists {
				return nil, fmt.Errorf("ERROR: invalid data format, duplicate room")
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
				links[room1] = append(links[room1], room2)
				links[room2] = append(links[room2], room1)
			}
		}

	}
	if startRoom == "" || endRoom == "" {
		return nil, fmt.Errorf("ERROR: invalid data format, invalid colony format")
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
