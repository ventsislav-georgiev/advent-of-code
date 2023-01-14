package main

import (
	"io"

	"github.com/ventsislav-georgiev/advent-of-code/golang/pkg/aoc"
)

func main() {
	aoc.Exec(task1, task2)
}

func task1(in io.Reader) {
	screen := map[Point]int{}
	solve(in, screen, nil)

	blocks := 0
	for _, tile := range screen {
		if tile == block {
			blocks++
		}
	}

	printScreen(screen)
	println("blocks:", blocks)
}

func task2(in io.Reader) {
	screen := map[Point]int{}
	quarters := 2
	score := solve(in, screen, &quarters)

	printScreen(screen)
	println("score:", score)
}

func solve(in io.Reader, screen map[Point]int, quarters *int) int {
	intcode := aoc.ParseIntcode(in)

	if quarters != nil {
		intcode.SetProgram(0, *quarters)
	}

	go intcode.Run()

	pos := Point{0, 0}
	ballpos := Point{0, 0}
	paddlepos := Point{0, 0}
	score := 0

	moveStick := func(tile int) {
		if tile == paddle && paddlepos.x != 0 {
			return
		}

		if ballpos.x < paddlepos.x {
			intcode.In <- -1
		} else if ballpos.x > paddlepos.x {
			intcode.In <- 1
		} else {
			intcode.In <- 0
		}
	}

	for {
		x, ok := <-intcode.Out
		if !ok {
			break
		}

		pos.x = x
		pos.y = <-intcode.Out
		tile := <-intcode.Out

		if pos.x == -1 && pos.y == 0 {
			score = tile
			continue
		}

		screen[pos] = tile

		switch tile {
		case paddle:
			paddlepos = pos
		case ball:
			ballpos = pos
		default:
			continue
		}

		moveStick(tile)
	}

	return score
}

func printScreen(screen map[Point]int) {
	minx, miny := 0, 0
	maxx, maxy := 0, 0
	for pos := range screen {
		if pos.x < minx {
			minx = pos.x
		}
		if pos.x > maxx {
			maxx = pos.x
		}
		if pos.y < miny {
			miny = pos.y
		}
		if pos.y > maxy {
			maxy = pos.y
		}
	}

	for y := miny; y <= maxy; y++ {
		for x := minx; x <= maxx; x++ {
			switch screen[Point{x, y}] {
			case empty:
				print(" ")
			case wall:
				print("█")
			case block:
				print("#")
			case paddle:
				print("▒")
			case ball:
				print("O")
			}
		}
		println()
	}
}

const (
	empty = iota
	wall
	block
	paddle
	ball
)

type Point struct {
	x, y int
}
