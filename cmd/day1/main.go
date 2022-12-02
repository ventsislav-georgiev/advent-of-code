package main

import (
	"bufio"
	"fmt"
	"io"
	"sort"
	"strconv"

	"github.com/ventsislav-georgiev/advent-of-code-22/pkg/aoc"
)

func main() {
	task := aoc.GetTask(task1, task2)
	in := aoc.GetInput(1)
	defer in.Close()
	task(in)
}

func task1(in io.ReadCloser) {
	scanner := bufio.NewScanner(in)

	max := 0
	cur := 0
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			cur = 0
		}

		val, _ := strconv.Atoi(line)
		cur += val

		if cur > max {
			max = cur
		}
	}

	fmt.Println(max)
}

func task2(in io.ReadCloser) {
	scanner := bufio.NewScanner(in)

	totals := make([]int, 0, 10)

	max := 0
	cur := 0
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			totals = append(totals, cur)
			cur = 0
		}

		val, _ := strconv.Atoi(line)
		cur += val

		if cur > max {
			max = cur
		}
	}

	totals = append(totals, cur)

	sort.Ints(totals)
	len := len(totals)
	fmt.Println(totals[len-1] + totals[len-2] + totals[len-3])
}
