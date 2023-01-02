package main

import (
	"bufio"
	"io"
	"strings"

	"github.com/ventsislav-georgiev/advent-of-code/golang/pkg/aoc"
)

func main() {
	aoc.Exec(task1, task2)
}

func task1(in io.Reader) {
	intcode := parse(in)
	println(solve(intcode, 12, 2))
}

func task2(in io.Reader) {
	intcode := parse(in)
	target := 19690720

	noun := 99
	for ; noun >= 0; noun-- {
		if solve(intcode, noun, 0) <= target {
			break
		}
	}

	for verb := 99; verb >= 0; verb-- {
		out := solve(intcode, noun, verb)
		if out == target {
			println(100*noun + verb)
			return
		}
	}
}

func parse(in io.Reader) []int {
	scanner := bufio.NewScanner(in)
	scanner.Scan()
	strcode := strings.Split(scanner.Text(), ",")
	intcode := make([]int, len(strcode))
	for i, s := range strcode {
		intcode[i] = aoc.StrToInt(s)
	}
	return intcode
}

func solve(intcode []int, noun, verb int) int {
	tmp := make([]int, len(intcode))
	copy(tmp, intcode)
	intcode = tmp

	intcode[1] = noun
	intcode[2] = verb

	for i := 0; i < len(intcode); i += 4 {
		opcode := intcode[i]
		if opcode == 99 {
			break
		}

		in1 := intcode[intcode[i+1]]
		in2 := intcode[intcode[i+2]]
		out := intcode[i+3]

		result := in1
		switch opcode {
		case 1:
			result += in2
		case 2:
			result *= in2
		}

		intcode[out] = result
	}

	return intcode[0]
}
