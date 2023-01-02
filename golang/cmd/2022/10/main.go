package main

import (
	"bufio"
	"fmt"
	"io"
	"math"

	"github.com/ventsislav-georgiev/advent-of-code/golang/pkg/aoc"
)

func main() {
	aoc.Exec(solve, solve)
}

var silent bool

func solve(in io.Reader) {
	register := 1
	signal := 0
	clock := 20
	cycle := 1
	crt := [6][40]rune{}

	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		op := 1
		n := 0
		cmd := scanner.Bytes()
		if cmd[0] == 'a' {
			op = 2
			n = aoc.Atoi(cmd[5:])
		}

		for ; op > 0; op-- {
			if cycle == clock {
				clock += 40
				signal += cycle * register
			}
			draw(&crt, cycle, register)
			cycle++
		}

		register += n
	}

	if silent {
		return
	}

	for _, row := range crt {
		for _, col := range row {
			fmt.Print(string(col))
		}
		fmt.Println()
	}

	fmt.Println(signal)
}

func draw(crt *[6][40]rune, cycle, register int) {
	row := int(math.Ceil(float64(cycle)/40)) - 1
	col := (cycle - 1) % 40
	char := ' '
	if col == register || col == register+1 || col == register-1 {
		char = '█'
	}
	crt[row][col] = char
}
