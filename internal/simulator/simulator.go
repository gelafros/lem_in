package simulator

import (
	"fmt"

	"lem-in/internal/parser"
)

func Simulator(colony *parser.Colony, paths [][]string) error {
	for _, v := range colony.InputLines {
		fmt.Println(v)
	}
	return nil
}
