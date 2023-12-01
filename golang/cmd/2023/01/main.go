package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"

	"github.com/ventsislav-georgiev/advent-of-code/golang/pkg/aoc"
)

func main() {
	aoc.Exec(task1, task2)
}

func task1(in io.Reader) {
	scanner := bufio.NewScanner(in)
	scanner.Split(bufio.ScanLines)

	sum := 0

	for scanner.Scan() {
		line := scanner.Text()
		firstDigit := 0
		lastDigit := 0

		for _, c := range line {
			if c > 48 && c < 58 {
				firstDigit = int(c) - 48
				break
			}
		}

		for i := len(line) - 1; i >= 0; i-- {
			c := line[i]
			if c > 48 && c < 58 {
				lastDigit = int(c) - 48
				break
			}
		}

		num, _ := strconv.Atoi(strconv.Itoa(firstDigit) + strconv.Itoa(lastDigit))
		sum += num
	}

	fmt.Println(sum)
}

func task2(in io.Reader) {
	scanner := bufio.NewScanner(in)
	scanner.Split(bufio.ScanLines)

	numMap := map[string]int{
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

	sum := 0

	for scanner.Scan() {
		line := scanner.Text()
		firstDigit := 0
		lastDigit := 0

		for i, c := range line {
			if c > 48 && c < 58 {
				firstDigit = int(c) - 48
				goto next
			} else {
				for word, num := range numMap {
					if len(line) > i+len(word) && line[i:i+len(word)] == word {
						firstDigit = num
						goto next
					}
				}
			}
		}

	next:
		for i := len(line) - 1; i >= 0; i-- {
			c := line[i]
			if c > 48 && c < 58 {
				lastDigit = int(c) - 48
				goto sum
			} else {
				for word, num := range numMap {
					if i-len(word)+1 > -1 && line[i-len(word)+1:i+1] == word {
						lastDigit = num
						goto sum
					}
				}
			}
		}

	sum:
		num, _ := strconv.Atoi(strconv.Itoa(firstDigit) + strconv.Itoa(lastDigit))
		sum += num
	}

	fmt.Println(sum)
}
