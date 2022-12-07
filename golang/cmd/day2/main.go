package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/ventsislav-georgiev/advent-of-code-22/golang/pkg/aoc"
)

func main() {
	aoc.Exec(task1, task2)
}

var points = map[string]int{
	"A": 0,
	"B": 1,
	"C": 2,
	"X": 0,
	"Y": 1,
	"Z": 2,
}

func task1(in io.Reader) {
	scanner := bufio.NewScanner(in)

	result := 0
	for scanner.Scan() {
		f := strings.Split(scanner.Text(), " ")
		p1 := points[f[0]]
		p2 := points[f[1]]

		winner := (3 + p1 - p2) % 3
		switch winner {
		case 0:
			result += 3 + p2 + 1
		case 1:
			result += 0 + p2 + 1
		case 2:
			result += 6 + p2 + 1
		}
	}

	fmt.Println(result)
}

func task2(in io.Reader) {
	scanner := bufio.NewScanner(in)

	result := 0
	for scanner.Scan() {
		f := strings.Split(scanner.Text(), " ")
		p1 := points[f[0]]
		p2 := points[f[1]]

		switch p2 {
		case 0:
			p2 = p1 - 1
		case 1:
			p2 = p1
		case 2:
			p2 = p1 + 1
		}

		if p2 < 0 {
			p2 = 2
		} else if p2 > 2 {
			p2 = 0
		}

		winner := (3 + p1 - p2) % 3
		switch winner {
		case 0:
			result += 3 + p2 + 1
		case 1:
			result += 0 + p2 + 1
		case 2:
			result += 6 + p2 + 1
		}
	}

	fmt.Println(result)
}
