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

	var totalPoints uint
	for scanner.Scan() {
		var ch byte

		for ch != ':' {
			scanner.Scan()
			ch = scanner.Bytes()[0]
		}

		nums1 := parseNumbers(scanner, '|')
		nums2 := parseNumbers(scanner, '\n')

		var points uint
		for n := range nums1 {
			if _, match := nums2[n]; !match {
				continue
			}

			if points == 0 {
				points = 1
			} else {
				points *= 2
			}
		}

		totalPoints += points
	}

	fmt.Println(totalPoints)
}

func task2(in io.Reader) {
	scanner := bufio.NewScanner(in)
	scanner.Split(bufio.ScanBytes)

	var cardIdx uint
	cardsCount := map[uint]uint{}

	for scanner.Scan() {
		var ch byte

		cardIdx++
		cardsCount[cardIdx]++

		for ch != ':' {
			scanner.Scan()
			ch = scanner.Bytes()[0]
		}

		nums1 := parseNumbers(scanner, '|')
		nums2 := parseNumbers(scanner, '\n')

		var cards uint
		for n := range nums1 {
			if _, match := nums2[n]; !match {
				continue
			}

			cards += 1
		}

		for c := cardsCount[cardIdx]; c > 0; c-- {
			for n := cardIdx + 1; n <= cardIdx+cards; n++ {
				cardsCount[n]++
			}
		}
	}

	var sum uint
	for _, count := range cardsCount {
		sum += count
	}

	fmt.Println(sum)
}

func parseNumbers(scanner *bufio.Scanner, term byte) map[uint]struct{} {
	var ch byte
	numbers := map[uint]struct{}{}
	numParts := []byte{}

	for ch != term {
		if !scanner.Scan() {
			break
		}

		ch = scanner.Bytes()[0]

		if ch >= '0' && ch <= '9' {
			numParts = append(numParts, ch)
		}

		if (ch == ' ' || ch == term) && len(numParts) > 0 {
			numbers[aoc.Atoui(numParts)] = struct{}{}
			numParts = numParts[:0]
		}
	}

	return numbers
}
