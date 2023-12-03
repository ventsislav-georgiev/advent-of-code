package main

import (
	"bufio"
	"image"
	"io"
	"strings"

	"github.com/ventsislav-georgiev/advent-of-code/golang/pkg/aoc"
)

func main() {
	aoc.Exec(task1, task2)
}

func task1(in io.Reader) {
	in = strings.NewReader(`
         A           
         A           
  #######.#########  
  #######.........#  
  #######.#######.#  
  #######.#######.#  
  #######.#######.#  
  #####  B    ###.#  
BC...##  C    ###.#  
  ##.##       ###.#  
  ##...DE  F  ###.#  
  #####    G  ###.#  
  #########.#####.#  
DE..#######...###.#  
  #.#########.###.#  
FG..#########.....#  
  ###########.#####  
             Z       
             Z       `)
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
				grid.updateLetters(pos1, pos2, pos1, pos1)
			} else {
				grid.updateLetters(pos1, pos2, pos1, pos2)
			}
		}
	}

	for y := 0; y < size.Y; y++ {
		for x := 0; x < size.X; x++ {
			print(grid[image.Pt(x, y)])
		}
		println()
	}
}

func task2(in io.Reader) {
}

type Node struct {
	val   string
	pos   image.Point
	edges []*Node
}

var directions = []image.Point{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

type Grid map[image.Point]string

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

func (g Grid) updateLetters(pos1, pos2, srcPos, destPos image.Point) {
	if pos2.X > pos1.X || pos2.Y > pos1.Y {
		g[destPos] = g[pos1] + g[pos2]
		g[pos2] = " "
	} else {
		g[destPos] = g[pos2] + g[pos1]
		g[pos1] = " "
	}
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
