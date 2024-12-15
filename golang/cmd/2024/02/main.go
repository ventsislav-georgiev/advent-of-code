package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/ventsislav-georgiev/advent-of-code/golang/pkg/aoc"
)

func main() {
	aoc.Exec(task1, task2)
}

func task1(in io.Reader) {
	scanner := bufio.NewScanner(in)

	safeCount := 0
	for scanner.Scan() {
		numbers := parseNumbers(scanner)
		if isArraySafe(numbers) {
			safeCount += 1
		}
	}

	fmt.Println(safeCount)
}

func task2(in io.Reader) {
	scanner := bufio.NewScanner(in)

	safeCount := 0
	for scanner.Scan() {
		numbers := parseNumbers(scanner)

		if isArraySafe(numbers) {
			safeCount += 1
			continue
		}

		numsCopy := make([]int, len(numbers))
		for i := 0; i < len(numbers); i++ {
			copy(numsCopy, numbers)

			var numsView []int
			if i == 0 {
				numsView = numsCopy[1:]
			} else {
				numsView = append(numsCopy[:i], numsCopy[i+1:]...)
			}

			if isArraySafe(numsView) {
				safeCount += 1
				break
			}
		}
	}

	fmt.Println(safeCount)
}

func parseNumbers(scanner *bufio.Scanner) []int {
	numbersString := strings.Split(scanner.Text(), " ")
	numbers := []int{}
	for _, numStr := range numbersString {
		numbers = append(numbers, aoc.StrToInt(numStr))
	}

	return numbers
}

func isArraySafe(numbers []int) bool {
	isInc := numbers[1]-numbers[0] > 0
	for j := 0; j+1 < len(numbers); j++ {
		if isSafe(numbers[j], numbers[j+1], isInc) {
			return false
		}
	}

	return true
}

func isSafe(n1, n2 int, isInc bool) bool {
	r := n2 - n1
	return (isInc && (r <= 0 || r > 3)) || (!isInc && (r >= 0 || r < -3))
}
