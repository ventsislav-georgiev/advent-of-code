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
	var sum uint

	for scanner.Scan() {
		line := scanner.Bytes()
		var firstDigit, lastDigit uint

		for _, ch := range line {
			if n, ok := isNum(ch); ok {
				firstDigit = n
				break
			}
		}

		for i := len(line) - 1; i >= 0; i-- {
			ch := line[i]
			if n, ok := isNum(ch); ok {
				lastDigit = n
				break
			}
		}

		sum += firstDigit*10 + lastDigit
	}

	fmt.Println(sum)
}

func task2(in io.Reader) {
	scanner := bufio.NewScanner(in)

	digitsMap := map[string]uint{
		"one":   1,
		"two":   2,
		"three": 3,
		"four":  4,
		"five":  5,
		"six":   6,
		"seven": 7,
		"eight": 8,
		"nine":  9,
	}

	var sum uint

	for scanner.Scan() {
		line := scanner.Bytes()
		var firstDigit, lastDigit uint

		for i, ch := range line {
			if n, ok := isNum(ch); ok {
				firstDigit = n
				break
			}

			for digit, num := range digitsMap {
				if len(line) <= i+len(digit) {
					continue
				}

				endIndex := i + len(digit)
				possibleDigit := line[i:endIndex]

				if compareSlice(possibleDigit, []byte(digit)) {
					firstDigit = num
					goto last_digit
				}
			}
		}

	last_digit:
		for i := len(line) - 1; i >= 0; i-- {
			ch := line[i]

			if n, ok := isNum(ch); ok {
				lastDigit = n
				break
			}

			endIndex := i + 1
			for digit, num := range digitsMap {
				if i-len(digit)+1 < 0 {
					continue
				}

				beginIndex := i - len(digit) + 1
				possibleDigit := line[beginIndex:endIndex]

				if compareSlice(possibleDigit, []byte(digit)) {
					lastDigit = num
					goto sum
				}
			}
		}

	sum:
		sum += firstDigit*10 + lastDigit
	}

	fmt.Println(sum)
}

func isNum(ch byte) (uint, bool) {
	if ch > 48 && ch < 58 {
		return uint(ch - 48), true
	}

	return 0, false
}

func compareSlice(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}

	for i, ch := range a {
		if ch != b[i] {
			return false
		}
	}

	return true
}
