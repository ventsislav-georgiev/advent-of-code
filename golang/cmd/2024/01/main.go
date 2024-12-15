package main

import (
	"bufio"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/ventsislav-georgiev/advent-of-code/golang/pkg/aoc"
)

func main() {
	aoc.Exec(task1, task2)
}

func task1(in io.Reader) {
	scanner := bufio.NewScanner(in)
	leftList := []int{}
	rightList := []int{}

	for scanner.Scan() {
		n1, n2 := parseLine(scanner)
		leftList = append(leftList, n1)
		rightList = append(rightList, n2)
	}

	sort.IntSlice(leftList).Sort()
	sort.IntSlice(rightList).Sort()

	dist := 0
	for i := 0; i < len(leftList); i++ {
		dist += aoc.Abs(leftList[i] - rightList[i])
	}

	fmt.Println(dist)
}

func task2(in io.Reader) {
	scanner := bufio.NewScanner(in)
	leftList := []int{}
	rightList := map[int]int{}

	for scanner.Scan() {
		n1, n2 := parseLine(scanner)
		leftList = append(leftList, n1)
		rightList[n2] += 1
	}

	similarity := 0
	for _, n := range leftList {
		count := rightList[n]
		similarity += n * count
	}

	fmt.Println(similarity)
}

func parseLine(scanner *bufio.Scanner) (int, int) {
	numbersString := strings.Split(scanner.Text(), "   ")
	return aoc.StrToInt(numbersString[0]), aoc.StrToInt(numbersString[1])
}
