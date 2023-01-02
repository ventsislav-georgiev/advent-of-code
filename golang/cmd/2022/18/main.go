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
	cords, grid, _ := parse(in)
	exposed := 0
	for _, c := range cords {
		x, y, z := c[0], c[1], c[2]
		left := grid[toKey(x+1, y, z)]
		right := grid[toKey(x-1, y, z)]
		top := grid[toKey(x, y+1, z)]
		bottom := grid[toKey(x, y-1, z)]
		front := grid[toKey(x, y, z+1)]
		back := grid[toKey(x, y, z-1)]
		if !left {
			exposed++
		}
		if !right {
			exposed++
		}
		if !top {
			exposed++
		}
		if !bottom {
			exposed++
		}
		if !front {
			exposed++
		}
		if !back {
			exposed++
		}
	}

	println(exposed)
}

func task2(in io.Reader) {
	_, grid, bounds := parse(in)

	visited := map[uint64]bool{}
	q := [][3]int16{{-1, -1, -1}}
	exposed := 0

	for len(q) > 0 {
		x, y, z := q[0][0], q[0][1], q[0][2]
		key := toKey(x, y, z)
		q = q[1:]

		if visited[key] {
			continue
		}
		if x < -1 || x > bounds[0] || y < -1 || y > bounds[1] || z < -1 || z > bounds[2] {
			continue
		}
		if grid[key] {
			exposed++
			continue
		}

		visited[key] = true
		q = append(q, [][3]int16{
			{x - 1, y, z},
			{x + 1, y, z},
			{x, y - 1, z},
			{x, y + 1, z},
			{x, y, z - 1},
			{x, y, z + 1},
		}...)
	}

	println(exposed)
}

func parse(in io.Reader) ([][3]int16, map[uint64]bool, [3]int16) {
	scanner := bufio.NewScanner(in)
	cords := make([][3]int16, 0)
	grid := map[uint64]bool{}
	bounds := [3]int16{0, 0, 0}

	for scanner.Scan() {
		var x, y, z int16
		fmt.Sscanf(scanner.Text(), "%d,%d,%d", &x, &y, &z)
		cords = append(cords, [3]int16{x, y, z})
		grid[toKey(x, y, z)] = true

		if x > bounds[0] {
			bounds[0] = x
		}
		if y > bounds[1] {
			bounds[1] = y
		}
		if z > bounds[2] {
			bounds[2] = z
		}
	}

	bounds[0]++
	bounds[1]++
	bounds[2]++

	return cords, grid, bounds
}

func toKey(x, y, z int16) uint64 {
	return uint64(x)<<32 + uint64(y)<<16 + uint64(z)
}
