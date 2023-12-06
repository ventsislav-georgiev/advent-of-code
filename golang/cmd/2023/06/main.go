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
	scanner := bufio.NewScanner(in)
	scanner.Split(bufio.ScanBytes)

	times := aoc.ParseNumbers(scanner, '\n')
	distances := aoc.ParseNumbers(scanner, '\n')

	total := 1
	for i, time := range times {
		count := 0
		distance := distances[i]

		for n := uint(1); n < time; n++ {
			if (time-n)*n > distance {
				count++
			}
		}

		total *= count
	}

	fmt.Println(total)
}

func task2(in io.Reader) {
	scanner := bufio.NewScanner(in)
	scanner.Scan()
	line := scanner.Bytes()[len("Time: "):]
	time := aoc.Atoui(aoc.RemoveSpaces(line))

	scanner.Scan()
	line = scanner.Bytes()[len("Distance: "):]
	distance := aoc.Atoui(aoc.RemoveSpaces(line))

	count := 0
	for n := uint(1); n < time; n++ {
		if (time-n)*n > distance {
			count++
		}
	}

	fmt.Println(count)
}
