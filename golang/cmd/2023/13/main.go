package main

import (
	"bufio"
	"fmt"
	"io"

	"github.com/ventsislav-georgiev/advent-of-code/golang/pkg/aoc"
)

func main() {
	aoc.Exec(task1, task2)
}

func task1(in io.Reader) {
	expectedMissmatch := 0
	result := solve(in, expectedMissmatch)
	fmt.Println(result)
}

func task2(in io.Reader) {
	expectedMissmatch := 1
	result := solve(in, expectedMissmatch)
	fmt.Println(result)
}

func solve(in io.Reader, expectedMissmatch int) int {
	scanner := bufio.NewScanner(in)
	matrix := Matrix{}
	var result int

	for scanner.Scan() {
		row := scanner.Text()

		if row == "" {
			result += matrix.Calc(expectedMissmatch)
			matrix = matrix[:0]
			continue
		}

		matrix = append(matrix, row)
	}

	result += matrix.Calc(expectedMissmatch)

	return result
}

type Matrix []string

func (m Matrix) Calc(expectedMissmatch int) (result int) {
	if len(m) == 0 {
		return 0
	}

	result += m.calcCols(expectedMissmatch)
	result += m.calcRows(expectedMissmatch)

	return
}

func (m Matrix) calcCols(expectedMissmatch int) (result int) {
	cols := len(m[0])
	rows := len(m)

	for col := 1; col < cols; col++ {
		var mismatch int

		for row := 0; row < rows; row++ {
			for cur := 0; cur < col; cur++ {
				mirror := col + (col - cur) - 1

				if mirror >= cols {
					continue
				}

				ch1 := m[row][cur]
				ch2 := m[row][mirror]
				if ch1 != ch2 {
					mismatch++
				}
			}
		}

		if mismatch == expectedMissmatch {
			result += col
		}
	}

	return
}

func (m Matrix) calcRows(expectedMissmatch int) (result int) {
	cols := len(m[0])
	rows := len(m)

	for row := 1; row < rows; row++ {
		var mismatch int

		for cur := 0; cur < row; cur++ {
			for col := 0; col < cols; col++ {
				mirror := row + (row - cur) - 1

				if mirror >= rows {
					continue
				}

				ch1 := m[cur][col]
				ch2 := m[mirror][col]
				if ch1 != ch2 {
					mismatch++
				}
			}
		}

		if mismatch == expectedMissmatch {
			result += row * 100
		}
	}

	return
}
