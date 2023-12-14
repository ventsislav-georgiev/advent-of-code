package aoc

import (
	"bufio"
	"io"
)

func ReadMatrix(in io.Reader) [][]byte {
	scanner := bufio.NewScanner(in)
	matrix := make([][]byte, 0)

	for scanner.Scan() {
		matrix = append(matrix, []byte(string(scanner.Bytes())))
	}

	return matrix
}
