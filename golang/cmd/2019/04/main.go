package main

import (
	"bufio"
	"io"
	"strconv"
	"strings"

	"github.com/ventsislav-georgiev/advent-of-code/golang/pkg/aoc"
)

func main() {
	aoc.Exec(task1, task2)
}

func task1(in io.Reader) {
	solve(in, true)
}

func task2(in io.Reader) {
	solve(in, false)
}

func solve(in io.Reader, allowMultiRepeat bool) {
	scanner := bufio.NewScanner(in)
	scanner.Split(bufio.ScanLines)

	var min, max int
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "-")
		min, _ = strconv.Atoi(parts[0])
		max, _ = strconv.Atoi(parts[1])
	}

	var count int
	for i := min; i <= max; i++ {
		if isValid(i, allowMultiRepeat) {
			count++
		}
	}

	println(count)
}

func isValid(n int, allowMultiRepeat bool) bool {
	repetitions := map[int]int{}

	for prev := 10; n > 0; n /= 10 {
		cur := n % 10
		repetitions[cur]++
		if cur > prev {
			return false
		}
		prev = cur
	}

	for _, v := range repetitions {
		if v < 2 {
			continue
		}

		if allowMultiRepeat || v == 2 {
			return true
		}
	}

	return false
}
