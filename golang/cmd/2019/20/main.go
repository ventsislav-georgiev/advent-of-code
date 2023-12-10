package main

import (
	"bufio"
	"fmt"
	"image"
	"io"

	"github.com/ventsislav-georgiev/advent-of-code/golang/pkg/aoc"
)

func main() {
	aoc.Exec(task1, task2)
}

func task1(in io.Reader) {
	grid, size := parse(in)

	for y := 0; y < size.Y; y++ {
		for x := 0; x < size.X; x++ {
			pos1 := image.Pt(x, y)
			cell := grid[pos1]
			if cell == " " || cell == "#" || cell == "." {
				continue
			}

			pos2, found := grid.firstImmediateNeighbourLetter(pos1)
			if !found {
				continue
			}

			if _, ok := grid.firstImmediateNeighbourPath(pos1); ok {
				grid.updateLetters(pos1, pos2)
			}
		}
	}

	coords := CordsMap{}
	for y := 0; y < size.Y; y++ {
		for x := 0; x < size.X; x++ {
			p := image.Pt(x, y)
			id := grid[p]

			if len(id) < 2 {
				continue
			}

			if _, ok := coords[id]; ok {
				id += "*"
				grid[p] = id
			}

			coords[id] = p
		}
	}

	start := coords["AA"]
	end := coords["ZZ"]
	queue := []image.Point{start}
	distances := map[image.Point]int{start: -1}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		if cur == end {
			continue
		}

		for _, dir := range directions {
			next := cur.Add(dir)
			if d, ok := distances[next]; ok && d <= distances[cur]+1 {
				continue
			}

			id := grid[next]
			if id == "#" || id == " " {
				continue
			}

			var portalID string
			if len(id) == 3 {
				portalID = id[:2]
			} else if len(id) == 2 {
				portalID = id + "*"
			}

			coord, foundPortal := coords[portalID]
			distCount := 1
			if foundPortal {
				next = coord
				distCount = 0
			}

			queue = append(queue, next)
			distances[next] = distances[cur] + distCount
		}
	}

	fmt.Println(distances[end] - 1)
}

func task2(in io.Reader) {
}

var directions = []image.Point{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

type Grid map[image.Point]string
type CordsMap map[string]image.Point

func (g Grid) firstImmediateNeighbourLetter(src image.Point) (image.Point, bool) {
	for _, dir := range directions {
		if g[src.Add(dir)] >= "A" && g[src.Add(dir)] <= "Z" {
			return src.Add(dir), true
		}
	}
	return image.Point{}, false
}

func (g Grid) firstImmediateNeighbourPath(src image.Point) (image.Point, bool) {
	for _, dir := range directions {
		if g[src.Add(dir)] == "." {
			return src.Add(dir), true
		}
	}
	return image.Point{}, false
}

func (g Grid) updateLetters(pos1, pos2 image.Point) {
	if pos1.X > pos2.X || pos1.Y > pos2.Y {
		g[pos1] = g[pos2] + g[pos1]
	} else {
		g[pos1] = g[pos1] + g[pos2]
	}

	g[pos2] = " "
}

func parse(in io.Reader) (Grid, image.Point) {
	grid := Grid{}
	scanner := bufio.NewScanner(in)

	x := 0
	y := 0
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}

		x = 0
		for _, ch := range line {
			grid[image.Pt(x, y)] = string(ch)
			x++
		}
		y++
	}

	return grid, image.Pt(x, y)
}
