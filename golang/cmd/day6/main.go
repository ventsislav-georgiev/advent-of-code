package main

import (
	"bufio"
	"fmt"
	"io"

	"github.com/ventsislav-georgiev/advent-of-code-22/golang/pkg/aoc"
)

func main() {
	aoc.Exec(task1, task2)
}

func task1(in io.ReadCloser) {
	r := bufio.NewReader(in)
	ch1, _, _ := r.ReadRune()
	ch2, _, _ := r.ReadRune()
	ch3, _, _ := r.ReadRune()
	i := 3

	for {
		ch4, _, err := r.ReadRune()
		i += 1
		if err != nil {
			break
		}

		if ch1 != ch2 && ch1 != ch3 && ch1 != ch4 && ch2 != ch3 && ch2 != ch4 && ch3 != ch4 {
			fmt.Print(i)
			return
		} else {
			ch1 = ch2
			ch2 = ch3
			ch3 = ch4
		}
	}
}

func task2(in io.ReadCloser) {
	r := bufio.NewReader(in)
	marker := make([]rune, 0)
	pos := 0
	for {
		ch, _, err := r.ReadRune()
		if err != nil {
			break
		}

		pos += 1
		for i, mch := range marker {
			if ch == mch {
				marker = marker[i+1:]
				break
			}
		}

		marker = append(marker, ch)

		if len(marker) == 14 {
			fmt.Print(pos)
			return
		}
	}
}
