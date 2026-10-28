package solver

import (
	"testing"

	"lem-in/internal/parser"
)

func TestSolveValidColonies(t *testing.T) {
	files := []struct {
		path      string
		minPaths  int
	}{
		{"../../tests/example00.txt", 1},
		{"../../tests/example01.txt", 2},
		{"../../tests/example02.txt", 2},
		{"../../tests/example03.txt", 1},
		{"../../tests/example04.txt", 1},
		{"../../tests/example05.txt", 2},
	}

	for _, tc := range files {
		colony, err := parser.Parse(tc.path)
		if err != nil {
			t.Fatalf("Parse(%s) failed: %v", tc.path, err)
		}
		paths, err := Solve(colony)
		if err != nil {
			t.Fatalf("Solve(%s) failed: %v", tc.path, err)
		}
		if len(paths) < tc.minPaths {
			t.Errorf("Solve(%s): got %d paths, want at least %d", tc.path, len(paths), tc.minPaths)
		}

		// Verify all paths are vertex-disjoint on intermediate rooms
		seen := make(map[string]bool)
		for _, p := range paths {
			if p[0] != colony.StartRoom || p[len(p)-1] != colony.EndRoom {
				t.Errorf("path %v does not start at %s and end at %s", p, colony.StartRoom, colony.EndRoom)
			}
			for _, room := range p[1 : len(p)-1] {
				if seen[room] {
					t.Errorf("intermediate room %q used in multiple paths: %v", room, paths)
				}
				seen[room] = true
			}
		}
	}
}

func TestSolveDisconnectedGraph(t *testing.T) {
	for _, f := range []string{"../../tests/badexample00.txt", "../../tests/badexample01.txt"} {
		colony, err := parser.Parse(f)
		if err != nil {
			continue
		}
		if _, err := Solve(colony); err == nil {
			t.Errorf("Solve(%s) expected error for disconnected graph, got nil", f)
		}
	}
}
