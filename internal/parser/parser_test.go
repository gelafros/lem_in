package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempColony(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "colony.txt")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp colony file: %v", err)
	}
	return p
}

func TestParseValidColonies(t *testing.T) {
	files := []string{
		"../../tests/example00.txt",
		"../../tests/example01.txt",
		"../../tests/example02.txt",
		"../../tests/example03.txt",
		"../../tests/example04.txt",
		"../../tests/example05.txt",
		"../../tests/example06.txt",
		"../../tests/example07.txt",
	}

	for _, f := range files {
		c, err := Parse(f)
		if err != nil {
			t.Errorf("Parse(%s) unexpected error: %v", f, err)
			continue
		}
		if c.AntsCount <= 0 {
			t.Errorf("Parse(%s) expected AntsCount > 0, got %d", f, c.AntsCount)
		}
		if c.StartRoom == "" || c.EndRoom == "" {
			t.Errorf("Parse(%s) missing start or end room", f)
		}
	}
}

func TestParseErrors(t *testing.T) {
	tests := []struct {
		name    string
		content string
	}{
		{"empty_file", ""},
		{"non_numeric_ants", "abc\n##start\n0 0 0\n##end\n1 1 1\n0-1\n"},
		{"zero_ants", "0\n##start\n0 0 0\n##end\n1 1 1\n0-1\n"},
		{"negative_ants", "-5\n##start\n0 0 0\n##end\n1 1 1\n0-1\n"},
		{"duplicate_start", "3\n##start\n0 0 0\n##start\n2 2 2\n##end\n1 1 1\n0-1\n"},
		{"duplicate_end", "3\n##start\n0 0 0\n##end\n1 1 1\n##end\n2 2 2\n0-1\n"},
		{"room_starts_with_L", "3\n##start\nLroom 0 0\n##end\n1 1 1\nLroom-1\n"},
		{"invalid_x_coord", "3\n##start\n0 x 0\n##end\n1 1 1\n0-1\n"},
		{"invalid_y_coord", "3\n##start\n0 0 y\n##end\n1 1 1\n0-1\n"},
		{"duplicate_coordinates", "3\n##start\n0 1 1\n##end\n1 1 1\n0-1\n"},
		{"duplicate_room_name", "3\n##start\n0 0 0\n0 2 2\n##end\n1 1 1\n0-1\n"},
		{"link_unknown_room", "3\n##start\n0 0 0\n##end\n1 1 1\n0-99\n"},
		{"self_linking_room", "3\n##start\n0 0 0\n##end\n1 1 1\n0-0\n"},
		{"duplicate_link", "3\n##start\n0 0 0\n##end\n1 1 1\n0-1\n1-0\n"},
		{"malformed_line", "3\n##start\n0 0 0\n##end\n1 1 1\nmalformed line with four tokens\n"},
		{"missing_start_or_end", "3\n0 0 0\n1 1 1\n0-1\n"},
		{"start_equals_end", "3\n##start\n##end\n0 0 0\n1 1 1\n0-1\n"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := writeTempColony(t, tc.content)
			if _, err := Parse(p); err == nil {
				t.Errorf("expected error for %s, got nil", tc.name)
			}
		})
	}

	if _, err := Parse("non_existent_file.txt"); err == nil {
		t.Errorf("expected error for non-existent file")
	}
}
