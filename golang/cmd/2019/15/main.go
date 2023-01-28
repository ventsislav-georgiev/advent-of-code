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

const (
	wall   = 0
	oxygen = 2
)

const (
	north = 1
	south = 2
	west  = 3
	east  = 4
)

var directions = map[int]image.Point{
	north: {0, -1},
	south: {0, 1},
	west:  {-1, 0},
	east:  {1, 0},
}

func task1(in io.Reader) {
	intcode := aoc.ParseIntcode(in)
	go intcode.Run()
	_, _, stepsToOxy := exploreMaze(intcode, true)
	fmt.Println(stepsToOxy)
}

func task2(in io.Reader) {
	intcode := aoc.ParseIntcode(in)
	go intcode.Run()
	maze, oxyPos, _ := exploreMaze(intcode, false)

	visited := map[image.Point]int{oxyPos: 0}
	queue := []image.Point{oxyPos}
	var pos image.Point

	for len(queue) > 0 {
		pos = queue[0]
		queue = queue[1:]

		for _, dirPoint := range directions {
			nextPos := pos.Add(dirPoint)
			if _, ok := visited[nextPos]; ok {
				continue
			}

			if maze[nextPos] == wall {
				continue
			}

			dist := visited[pos] + 1
			visited[nextPos] = dist
			queue = append(queue, nextPos)
		}
	}

	fmt.Println(visited[pos])
}

func exploreMaze(intcode *aoc.Intcode, exitOnOxy bool) (maze map[image.Point]int, oxyPos image.Point, stepsToOxy int) {
	maze = map[image.Point]int{}
	pos := image.Point{0, 0}
	var path []int

explore:
	for {
		for dirOp, dirPoint := range directions {
			nextPos := pos.Add(dirPoint)
			if _, ok := maze[nextPos]; ok {
				continue
			}

			intcode.In <- dirOp
			maze[nextPos] = <-intcode.Out

			if maze[nextPos] == oxygen {
				oxyPos = nextPos
				stepsToOxy = len(path) + 1

				if exitOnOxy {
					return nil, oxyPos, stepsToOxy
				}
			}

			if maze[nextPos] == wall {
				continue
			}

			pos = nextPos
			oppositeDir := ((dirOp - 1) ^ 1) + 1
			path = append(path, oppositeDir)
			continue explore
		}

		if len(path) == 0 {
			break
		}

		prevDir := path[len(path)-1]
		intcode.In <- prevDir
		<-intcode.Out

		pos = pos.Add(directions[prevDir])
		path = path[:len(path)-1]
	}

	return maze, oxyPos, stepsToOxy
}
