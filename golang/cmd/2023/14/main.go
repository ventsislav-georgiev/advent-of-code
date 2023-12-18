package main

import (
	"fmt"
	"io"

	"github.com/ventsislav-georgiev/advent-of-code/golang/pkg/aoc"
)

func main() {
	aoc.Exec(task1, task2)
}

func task1(in io.Reader) {
	matrix := aoc.ReadMatrixAsBytes(in).Rows
	tiltNorth(matrix)
	fmt.Println(calcLoad(matrix))
}

func task2(in io.Reader) {
	matrix := aoc.ReadMatrixAsBytes(in).Rows

	totalCycles := 1_000_000_000
	visited := map[string]int{}
	var curCycle, patternLen int

	for ; totalCycles > 0; curCycle++ {
		hash := hashMatrix(matrix)

		prevCycle, found := visited[hash]
		if found {
			patternLen = curCycle - prevCycle
			break
		}

		visited[hash] = curCycle
		runCycle(matrix)
	}

	remainingCycles := (totalCycles - curCycle) % patternLen

	for ; remainingCycles > 0; remainingCycles-- {
		runCycle(matrix)
	}

	fmt.Println(calcLoad(matrix))
}

func hashMatrix(matrix [][]byte) string {
	var contents string
	for _, row := range matrix {
		contents += string(row)
	}

	hash := aoc.Hash(contents)
	return hash
}

func calcLoad(matrix [][]byte) int {
	var load int
	for y, row := range matrix {
		for _, ch := range row {
			if ch == 'O' {
				load += len(matrix) - y
			}
		}
	}

	return load
}

func runCycle(matrix [][]byte) {
	tiltNorth(matrix)
	tiltWest(matrix)
	tiltSouth(matrix)
	tiltEast(matrix)
}

func tiltNorth(matrix [][]byte) {
	for {
		var moved int

		for col := 0; col < len(matrix[0]); col++ {
			for row := 1; row < len(matrix); row++ {
				cur := matrix[row][col]
				colToNorth := matrix[row-1][col]

				if cur == 'O' && colToNorth == '.' {
					matrix[row][col] = '.'
					matrix[row-1][col] = 'O'
					moved++
				}
			}
		}

		if moved == 0 {
			break
		}
	}
}

func tiltSouth(matrix [][]byte) {
	for {
		var moved int

		for col := 0; col < len(matrix[0]); col++ {
			for row := len(matrix) - 2; row >= 0; row-- {
				cur := matrix[row][col]
				colToSouth := matrix[row+1][col]

				if cur == 'O' && colToSouth == '.' {
					matrix[row][col] = '.'
					matrix[row+1][col] = 'O'
					moved++
				}
			}
		}

		if moved == 0 {
			break
		}
	}
}

func tiltEast(matrix [][]byte) {
	for {
		var moved int

		for row := 0; row < len(matrix); row++ {
			for col := len(matrix[0]) - 2; col >= 0; col-- {
				cur := matrix[row][col]
				colToEast := matrix[row][col+1]

				if cur == 'O' && colToEast == '.' {
					matrix[row][col] = '.'
					matrix[row][col+1] = 'O'
					moved++
				}
			}
		}

		if moved == 0 {
			break
		}
	}
}

func tiltWest(matrix [][]byte) {
	for {
		var moved int

		for row := 0; row < len(matrix); row++ {
			for col := 1; col < len(matrix[0]); col++ {
				cur := matrix[row][col]
				colToWest := matrix[row][col-1]

				if cur == 'O' && colToWest == '.' {
					matrix[row][col] = '.'
					matrix[row][col-1] = 'O'
					moved++
				}
			}
		}

		if moved == 0 {
			break
		}
	}
}
