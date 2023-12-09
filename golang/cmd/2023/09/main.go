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
	reverse := false
	series := parse(in, reverse)
	fmt.Println(getNextValueSum(series))
}

func task2(in io.Reader) {
	reverse := true
	series := parse(in, reverse)
	fmt.Println(getNextValueSum(series))
}

func getNextValueSum(series [][]int) int {
	var nextValuesSum int
	for _, numbers := range series {
		diffs := calculateDiffsUntilZero(numbers)

		for i := aoc.LastIdx(diffs); i-1 >= 0; i-- {
			diff := aoc.LastElement(diffs[i]) + aoc.LastElement(diffs[i-1])
			diffs[i-1] = append(diffs[i-1], diff)
		}

		nextValuesSum += aoc.LastElement(numbers) + aoc.LastElement(diffs[0])
	}

	return nextValuesSum
}

func calculateDiffsUntilZero(numbers []int) [][]int {
	diffs := [][]int{}

	diffs = append(diffs, getDiff(numbers))
	for !allZeros(aoc.LastElement(diffs)) {
		diffs = append(diffs, getDiff(aoc.LastElement(diffs)))
	}

	return diffs
}

func getDiff(numbers []int) []int {
	diffs := make([]int, len(numbers)-1)

	for i := 1; i < len(numbers); i++ {
		diffs[i-1] = numbers[i] - numbers[i-1]
	}

	return diffs
}

func allZeros(n []int) bool {
	for _, num := range n {
		if num != 0 {
			return false
		}
	}

	return true
}

func parse(in io.Reader, reverse bool) [][]int {
	series := [][]int{}
	scanner := bufio.NewScanner(in)

	for scanner.Scan() {
		numbersRaw := aoc.SplitBytes(scanner.Bytes(), ' ')
		numbers := make([]int, 0, len(numbersRaw))

		if reverse {
			aoc.Reverse(numbersRaw)
		}

		for _, numRaw := range numbersRaw {
			num := aoc.Atoi(numRaw)
			numbers = append(numbers, num)
		}

		series = append(series, numbers)
	}

	return series
}
