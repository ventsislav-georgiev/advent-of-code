package main

import (
	"bufio"
	"fmt"
	"image"
	"io"
	"strings"

	"github.com/ventsislav-georgiev/advent-of-code/golang/pkg/aoc"
)

func main() {
	aoc.Exec(task1, task2)
}

func task1(in io.Reader) {
	in = strings.NewReader(`1,0,1~1,2,1
0,0,2~2,0,2
0,2,3~2,2,3
0,0,4~0,2,4
2,0,5~2,2,5
0,1,6~2,1,6
1,1,8~1,1,9`)

	bricks := [][3]image.Point{}
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "~")
		xyz1 := strings.Split(parts[0], ",")
		xyz2 := strings.Split(parts[1], ",")
		x := image.Point{X: aoc.StrToInt(xyz1[0]), Y: aoc.StrToInt(xyz2[0])}
		y := image.Point{X: aoc.StrToInt(xyz1[1]), Y: aoc.StrToInt(xyz2[1])}
		z := image.Point{X: aoc.StrToInt(xyz1[2]), Y: aoc.StrToInt(xyz2[2])}
		bricks = append(bricks, [3]image.Point{x, y, z})
	}

process:
	for {
		positionMoved := false
		for i, xyz1 := range bricks {
			x1 := xyz1[0]
			y1 := xyz1[1]
			z1 := xyz1[2]

			bricksToMoveDown := [][3]image.Point{}
			for j, xyz2 := range bricks {
				if i == j {
					continue
				}

				z2 := xyz2[2]
				if !aoc.CheckPointsOverlap(z1, z2) {
					bricksToMoveDown = append(bricksToMoveDown, xyz2)
				}
			}

			for _, xyz2 := range bricksToMoveDown {
				x2 := xyz2[0]
				y2 := xyz2[1]

				if aoc.CheckPointsOverlap(x1, x2) || aoc.CheckPointsOverlap(y1, y2) {
					continue
				}

				positionMoved = true
				xyz2[2].X -= 1
				xyz2[2].Y -= 1
			}
		}

		if !positionMoved {
			break process
		}
	}

	for _, xyz1 := range bricks {
		x1 := xyz1[0]
		y1 := xyz1[1]
		z1 := xyz1[2]

		fmt.Printf("Z=%v, X=%v, Y=%v\n", z1, x1, y1)
	}
}

func task2(in io.Reader) {
}
