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
	scanner.Scan()

	var result uint
	for {
		if len(scanner.Bytes()) == 0 {
			break
		}

		if !checkNextChars(scanner, []byte("mul(")) {
			continue
		}

		n1, ok := aoc.ParseNumber(scanner, ',')
		if !ok {
			continue
		}

		n2, ok := aoc.ParseNumber(scanner, ')')
		if !ok {
			continue
		}

		result += n1 * n2
	}

	fmt.Println(result)
}

func task2(in io.Reader) {
	scanner := bufio.NewScanner(in)
	scanner.Split(bufio.ScanBytes)
	scanner.Scan()

	var result uint
	enabled := true
	buf := [4]byte{}
	for {
		if len(scanner.Bytes()) == 0 {
			break
		}

		if !checkOneOfNextChars(scanner, []byte{'m', 'd'}) {
			continue
		}
		buf[0] = scanner.Bytes()[0]
		if !checkOneOfNextChars(scanner, []byte{'u', 'o'}) {
			continue
		}
		buf[1] = scanner.Bytes()[0]
		if !checkOneOfNextChars(scanner, []byte{'l', 'n', '('}) {
			continue
		}
		buf[2] = scanner.Bytes()[0]
		if !checkOneOfNextChars(scanner, []byte{'(', '\'', ')'}) {
			continue
		}
		buf[3] = scanner.Bytes()[0]

		if string(buf[:4]) == "do()" {
			enabled = true
			continue
		}

		if string(buf[:4]) == "mul(" {
			n1, ok := aoc.ParseNumber(scanner, ',')
			if !ok {
				continue
			}

			n2, ok := aoc.ParseNumber(scanner, ')')
			if !ok {
				continue
			}

			if enabled {
				result += n1 * n2
			}

			continue
		}

		if string(buf[:4]) == "don'" && checkNextChars(scanner, []byte("t()")) {
			enabled = false
		}
	}

	fmt.Println(result)
}

func checkNextChars(scanner *bufio.Scanner, expected []byte) bool {
	for _, ch := range expected {
		if !scanner.Scan() {
			return false
		}
		if scanner.Bytes()[0] != ch {
			return false
		}
	}
	return true
}

func checkOneOfNextChars(scanner *bufio.Scanner, expectedChars []byte) bool {
	if !scanner.Scan() {
		return false
	}

	for _, ch := range expectedChars {
		if scanner.Bytes()[0] == ch {
			return true
		}
	}

	return false
}
