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

	result := 0
	for scanner.Scan() {
		r := scanner.Text()
		halflen := len(r) / 2

		bitmap := toBitmap(r[:halflen]) & toBitmap(r[halflen:])
		for p := 0; p <= 52; p++ {
			if bitmap&(1<<p) != 0 {
				result += p + 1
			}
		}
	}

	fmt.Println(result)
}

func task2(in io.Reader) {
	scanner := bufio.NewScanner(in)

	result := 0
	for scanner.Scan() {
		r1 := scanner.Text()
		scanner.Scan()
		r2 := scanner.Text()
		scanner.Scan()
		r3 := scanner.Text()

		bitmap := toBitmap(r1) & toBitmap(r2) & toBitmap(r3)
		for p := 0; p <= 52; p++ {
			if bitmap&(1<<p) != 0 {
				result += p + 1
				break
			}
		}
	}

	fmt.Println(result)
}

func toBitmap(s string) uint64 {
	var bitmap uint64
	for _, ch := range s {
		if ch >= 'a' {
			bitmap |= 1 << (ch - 'a')
		} else {
			bitmap |= 1 << (ch - 'A' + 26)
		}
	}
	return bitmap
}
