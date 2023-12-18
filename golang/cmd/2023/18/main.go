package main

import (
	"bufio"
	"fmt"
	"image"
	"io"
	"strconv"
	"strings"

	"github.com/ventsislav-georgiev/advent-of-code/golang/pkg/aoc"
)

func main() {
	aoc.Exec(task1, task2)
}

func task1(in io.Reader) {
	useHexForDist := false
	result := solve(in, useHexForDist)
	fmt.Println(result)
}

func task2(in io.Reader) {
	useHexForDist := true
	result := solve(in, useHexForDist)
	fmt.Println(result)
}

func solve(in io.Reader, useHexForDist bool) int {
	pos := image.Point{0, 0}
	loop := []image.Point{pos}
	var perimeter int

	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := scanner.Text()
		dirParts := strings.Split(line, " ")
		dir := dirParts[0]
		dist, _ := strconv.Atoi(dirParts[1])

		code := dirParts[2][2:8]

		if useHexForDist {
			dist = 0
			dirID := code[len(code)-1] - '0'
			switch dirID {
			case 0:
				dir = "R"
			case 1:
				dir = "D"
			case 2:
				dir = "L"
			case 3:
				dir = "U"
			}

			distHex := code[:len(code)-1]

			for _, c := range distHex {
				dist *= 16
				if c >= 'a' {
					dist += int(c - 'a' + 10)
				} else {
					dist += int(c - '0')
				}
			}

		}

		perimeter += dist

		var moveVec image.Point
		switch dir {
		case "R":
			moveVec = image.Pt(dist, 0)
		case "L":
			moveVec = image.Pt(-dist, 0)
		case "U":
			moveVec = image.Pt(0, -dist)
		case "D":
			moveVec = image.Pt(0, dist)
		}
		pos = pos.Add(moveVec)
		loop = append(loop, pos)
	}

	area := int(aoc.CalcSimplePolygonArea(loop))

	// From A = B/2 + I - 1
	// ref: https://en.wikipedia.org/wiki/Pick%27s_theorem
	latticePoints := area + 1 - perimeter/2

	return latticePoints + perimeter
}
