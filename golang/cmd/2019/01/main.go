package main

import (
	"bufio"
	"io"

	"github.com/ventsislav-georgiev/advent-of-code/golang/pkg/aoc"
)

func main() {
	aoc.Exec(task1, task2)
}

func task1(in io.Reader) {
	solve(in, false)
}

func task2(in io.Reader) {
	solve(in, true)
}

func solve(in io.Reader, additional bool) {
	scanner := bufio.NewScanner(in)

	var sum int
	for scanner.Scan() {
		f := fuel(aoc.Atoi(scanner.Bytes()))

		if !additional {
			sum += f
			continue
		}

		for f >= 0 {
			sum += f
			f = fuel(f)
		}
	}

	println(sum)
}

func fuel(mass int) int {
	return mass/3 - 2
}
