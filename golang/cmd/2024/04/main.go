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
	matrix := aoc.ReadMatrixAsRunes(in).Rows
	count := 0

	for i := 0; i < len(matrix); i++ {
		for j := 0; j < len(matrix[0]); j++ {
			for _, word := range []string{"XMAS", "SAMX"} {
				if matrix[i][j] != rune(word[0]) {
					continue
				}

				count += checkWord(matrix, i, j, word)
			}
		}
	}

	fmt.Println(count)
}

func task2(in io.Reader) {
	matrix := aoc.ReadMatrixAsRunes(in).Rows
	count := 0
	visited := map[string]struct{}{}

	for i := 0; i < len(matrix); i++ {
		for j := 0; j < len(matrix[0]); j++ {
			for _, word := range []string{"MAS", "SAM"} {
				if matrix[i][j] != rune(word[0]) {
					continue
				}

				if checkWordCross(matrix, i, j, word, visited) {
					count++
				}
			}
		}
	}

	fmt.Println(count)
}

func checkWord(matrix [][]rune, i, j int, word string) int {
	count := 0

	if j+len(word) <= len(matrix[0]) {
		if checkWordHorizontally(matrix, i, j, word) {
			count++
		}
	}

	if i+len(word) <= len(matrix) {
		if checkWordVertically(matrix, i, j, word) {
			count++
		}
	}

	if i+len(word) <= len(matrix) && j+len(word) <= len(matrix[0]) {
		if checkWordDiagonallyDown(matrix, i, j, word) {
			count++
		}
	}

	if i-(len(word)-1) >= 0 && j+len(word) <= len(matrix[0]) {
		if checkWordDiagonallyUp(matrix, i, j, word) {
			count++
		}
	}

	return count
}

func checkWordCross(matrix [][]rune, i, j int, word string, visited map[string]struct{}) (found bool) {
	var loc string
	defer func() {
		if found {
			visited[loc] = struct{}{}
		}
	}()

	if i+len(word) <= len(matrix) && j+len(word) <= len(matrix[0]) {
		if checkWordDiagonallyDown(matrix, i, j, word) {
			loc = fmt.Sprintf("%d-%d", i+1, j+1)
			if _, ok := visited[loc]; ok {
				return false
			}

			if checkWordDiagonallyUp(matrix, i+len(word)-1, j, word) {
				return true
			}
			if checkWordDiagonallyUp(matrix, i+len(word)-1, j, aoc.ReverseString(word)) {
				return true
			}
		}
	}

	if i-(len(word)-1) >= 0 && j+len(word) <= len(matrix[0]) {
		if checkWordDiagonallyUp(matrix, i, j, word) {
			loc = fmt.Sprintf("%d-%d", i-1, j+1)
			if _, ok := visited[loc]; ok {
				return false
			}

			if checkWordDiagonallyDown(matrix, i-(len(word)-1), j, word) {
				return true
			}
			if checkWordDiagonallyDown(matrix, i-(len(word)-1), j, aoc.ReverseString(word)) {
				return true
			}
		}
	}

	return false
}

func checkWordHorizontally(matrix [][]rune, i, j int, word string) bool {
	for k, r := range word {
		if matrix[i][j+k] != r {
			return false
		}
	}

	return true
}

func checkWordVertically(matrix [][]rune, i, j int, word string) bool {
	for k, r := range word {
		if matrix[i+k][j] != r {
			return false
		}
	}

	return true
}

func checkWordDiagonallyDown(matrix [][]rune, i, j int, word string) bool {
	for k, r := range word {
		if matrix[i+k][j+k] != r {
			return false
		}
	}

	return true
}

func checkWordDiagonallyUp(matrix [][]rune, i, j int, word string) bool {
	for k, r := range word {
		if matrix[i-k][j+k] != r {
			return false
		}
	}

	return true
}
