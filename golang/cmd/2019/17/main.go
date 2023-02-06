package main

import (
	"fmt"
	"image"
	"io"
	"time"

	"github.com/ventsislav-georgiev/advent-of-code/golang/pkg/aoc"
)

func main() {
	aoc.Exec(task1, task2)
}

func task1(in io.Reader) {
	intcode := aoc.ParseIntcode(in)
	go intcode.Run()

	grid := make(map[image.Point]byte)
	row, col := 0, 0

	for out := range intcode.Out {
		fmt.Printf("%c", out)

		if out == 10 {
			row++
			col = 0
			continue
		}

		grid[image.Pt(col, row)] = byte(out)
		col++
	}

	sum := 0
	directions := []image.Point{{1, 0}, {0, 1}, {-1, 0}, {0, -1}}

	for p := range grid {
		if grid[p] != '#' {
			continue
		}

		for _, direction := range directions {
			if grid[p.Add(direction)] != '#' {
				goto skip
			}
		}

		sum += p.X * p.Y

	skip:
	}

	fmt.Println(sum)
}

func task2(in io.Reader) {
	intcode := aoc.ParseIntcode(in)
	intcode.Set(0, 2)

	go intcode.Run()

	go func() {
		for out := range intcode.Out {
			if out > 255 {
				fmt.Println(out)
			} else {
				fmt.Printf("%c", out)
			}
		}
	}()

	op := []string{"A", "A", "B", "C", "B", "C", "B", "C", "B", "A", "\n"}
	a := []string{"R", "6", "L", "12", "R", "6", "\n"}
	b := []string{"L", "12", "R", "6", "L", "8", "L", "12", "\n"}
	c := []string{"R", "12", "L", "10", "L", "10", "\n"}
	n := []string{"n", "\n"}

	time.Sleep(100 * time.Millisecond)
	for _, cmdset := range [][]string{op, a, b, c, n} {
		for idx, cmd := range cmdset {
			fmt.Printf("%s", cmd)

			if cmd == "\n" {
				intcode.In <- '\n'
				continue
			}

			for _, ch := range cmd {
				intcode.In <- int(ch)
			}

			if idx < len(cmdset)-2 {
				fmt.Print(",")
				intcode.In <- ','
			}
		}

		time.Sleep(100 * time.Millisecond)
	}

	intcode.WaitHalt()
}
