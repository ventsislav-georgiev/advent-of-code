package main

import (
	"io"

	"github.com/ventsislav-georgiev/advent-of-code/golang/pkg/aoc"
)

func main() {
	aoc.Exec(task1, task2)
}

func task1(in io.Reader) {
	panels := solve(in, 0)
	println(len(panels))
}

func task2(in io.Reader) {
	panels := solve(in, 1)

	minx, miny, maxx, maxy := 0, 0, 0, 0
	for p := range panels {
		if p.x < minx {
			minx = p.x
		}
		if p.y < miny {
			miny = p.y
		}
		if p.x > maxx {
			maxx = p.x
		}
		if p.y > maxy {
			maxy = p.y
		}
	}

	for y := miny; y <= maxy; y++ {
		for x := minx; x <= maxx; x++ {
			if panels[Point{x, y}] == 1 {
				print("█")
			} else {
				print(" ")
			}
		}
		println()
	}
}

func solve(in io.Reader, start int) map[Point]int {
	intcode := aoc.ParseIntcode(in)
	go intcode.Run()
	intcode.In <- start

	panels := map[Point]int{}
	pos := Point{0, 0}
	dir := Point{0, -1}

	for {
		newcolor, ok := <-intcode.Out
		if !ok {
			return panels
		}

		panels[pos] = newcolor

		rotation := <-intcode.Out
		if rotation == 0 {
			dir.x, dir.y = dir.y, -dir.x
		} else {
			dir.x, dir.y = -dir.y, dir.x
		}

		pos.x += dir.x
		pos.y += dir.y

		color, ok := panels[pos]
		if !ok {
			color = 0
		}

		intcode.In <- color
	}
}

type Point struct {
	x, y int
}
