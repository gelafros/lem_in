package simulator

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"lem-in/internal/parser"
	"lem-in/internal/solver"
)

func TestSimulatorOutput(t *testing.T) {
	colony, err := parser.Parse("../../tests/example00.txt")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	paths, err := solver.Solve(colony)
	if err != nil {
		t.Fatalf("Solve failed: %v", err)
	}

	origStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Pipe failed: %v", err)
	}
	os.Stdout = w

	simErr := Simulator(colony, paths)
	w.Close()
	os.Stdout = origStdout

	if simErr != nil {
		t.Fatalf("Simulator returned error: %v", simErr)
	}

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	out := buf.String()

	if !strings.Contains(out, "L1-") || !strings.Contains(out, "L4-"+colony.EndRoom) {
		t.Errorf("Simulator output missing expected ant movements:\n%s", out)
	}
}

func TestSimulatorDirectStartEndLink(t *testing.T) {
	colony := &parser.Colony{
		AntsCount:  3,
		StartRoom:  "start",
		EndRoom:    "end",
		InputLines: []string{"3", "##start", "start 0 0", "##end", "end 1 1", "start-end"},
	}
	paths := [][]string{{"start", "end"}}

	origStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	_ = Simulator(colony, paths)
	w.Close()
	os.Stdout = origStdout

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	out := buf.String()
	if !strings.Contains(out, "L1-end") || !strings.Contains(out, "L3-end") {
		t.Errorf("expected direct start->end moves for all 3 ants, got:\n%s", out)
	}
}
