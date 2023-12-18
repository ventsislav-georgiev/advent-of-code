package main

import (
	"fmt"
	"image"
	"io"

	"github.com/ventsislav-georgiev/advent-of-code/golang/pkg/aoc"
)

func main() {
	aoc.Exec(task1, task2)
}

func task1(in io.Reader) {
	expandRate := 2
	_, sum := getDistances(in, expandRate)
	fmt.Println(sum)
}

func task2(in io.Reader) {
	expandRate := 1000000
	_, sum := getDistances(in, expandRate)
	fmt.Println(sum)
}

func getDistances(in io.Reader, expandRate int) (map[[2]image.Point]int, int) {
	matrix := aoc.ReadMatrixAsBytes(in).Rows
	points := getPoints(matrix)
	points = shiftPoints(matrix, points, expandRate-1)

	results := map[[2]image.Point]int{}

	for i := 0; i < len(points)-1; i++ {
		start := points[i]
		for j := i + 1; j < len(points); j++ {
			end := points[j]
			resKey := [2]image.Point{*start, *end}
			distance := aoc.ManhattanDistance(*start, *end)
			results[resKey] = distance
		}
	}

	var sum int
	for _, v := range results {
		sum += v
	}

	return results, sum
}

func shiftPoints(matrix [][]byte, points []*image.Point, dist int) []*image.Point {
	rowsToInsert := map[int]struct{}{}
	for y := 0; y < len(matrix); y++ {
		for x := 0; x < len(matrix[y]); x++ {
			if matrix[y][x] == '#' {
				goto skip
			}
		}

		rowsToInsert[y] = struct{}{}

	skip:
	}

	colsToInsert := map[int]struct{}{}
	for x := 0; x < len(matrix[0]); x++ {
		for y := 0; y < len(matrix); y++ {
			if matrix[y][x] == '#' {
				goto skip2
			}
		}

		colsToInsert[x] = struct{}{}

	skip2:
	}

	yShift := 0
	for y := 0; y < len(matrix)+yShift; y++ {
		if _, ok := rowsToInsert[y]; ok {
			yIdx := y + yShift
			for _, p := range points {
				if p.Y < yIdx {
					continue
				}

				p.Y += dist
			}
			yShift += dist
		}
	}

	xShift := 0
	for x := 0; x < len(matrix[0])+xShift; x++ {
		if _, ok := colsToInsert[x]; ok {
			xIdx := x + xShift
			for _, p := range points {
				if p.X < xIdx {
					continue
				}

				p.X += dist
			}
			xShift += dist
		}
	}

	return points
}

func getPoints(matrix [][]byte) []*image.Point {
	points := []*image.Point{}

	for y, row := range matrix {
		for x, cell := range row {
			if cell == '#' {
				points = append(points, &image.Point{x, y})
			}
		}
	}

	return points
}
