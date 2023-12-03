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
	intcode := aoc.ParseIntcode(in)
	go intcode.Run()

	grid := make(map[image.Point]int)
	pointsAffectedByTheTractorBeam := 0

	for y := 0; y < 50; y += 1 {
		for x := 0; x < 50; x += 1 {
			intcode.In <- x
			intcode.In <- y
			result := <-intcode.Out
			if result == 1 {
				pointsAffectedByTheTractorBeam += 1
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
			grid[image.Pt(x, y)] = result

			intcode = intcode.Reset()
			go intcode.Run()
		}
		fmt.Println()
	}

	fmt.Println(pointsAffectedByTheTractorBeam)
}

func task2(in io.Reader) {
	// 	target := 10
	// 	yMin := 0
	// 	yMax := 50
	// 	xMin := 0
	// 	xLen := 50
	// 	xMax := xMin + xLen

	// 	input := `#.......................................
	// .#......................................
	// ..##....................................
	// ...###..................................
	// ....###.................................
	// .....####...............................
	// ......#####.............................
	// ......######............................
	// .......#######..........................
	// ........########........................
	// .........#########......................
	// ..........#########.....................
	// ...........##########...................
	// ...........############.................
	// ............############................
	// .............#############..............
	// ..............##############............
	// ...............###############..........
	// ................###############.........
	// ................#################.......
	// .................##################.....
	// ..................##################....
	// ...................###################..
	// ....................####################
	// .....................###################
	// .....................###################
	// ......................##################
	// .......................#################
	// ........................################
	// .........................###############
	// ..........................##############
	// ..........................##############
	// ...........................#############
	// ............................############
	// .............................###########`

	// 	scanner := bufio.NewScanner(strings.NewReader(input))

	// 	y := 0
	// 	for scanner.Scan() {
	// 		line := scanner.Text()
	// 		for x, cell := range line {
	// 			if cell == '#' {
	// 				grid[image.Pt(x, y)] = 1
	// 			} else if cell == '.' {
	// 				grid[image.Pt(x, y)] = 0
	// 			}
	// 		}
	// 		y += 1
	// 	}

	// 	for y := yMin; y < yMax; y += 1 {
	// 		fmt.Printf("%d: ", y)
	// 		c := 0
	// 		for x := xMin; x < xMax; x += 1 {
	// 			if grid[image.Pt(x, y)] == 0 {
	// 				fmt.Print(".")
	// 			} else {
	// 				fmt.Print("#")
	// 				c += 1
	// 			}
	// 		}
	// 		fmt.Printf(" (%d)\n", c)
	// 	}

	intcode := aoc.ParseIntcode(in)
	go intcode.Run()

	grid := make(map[image.Point]int)

	target := 100
	yMin := 1500
	yMax := 1800
	xMin := 1100
	xLen := 200
	xMax := xMin + xLen
	for y := yMin; y < yMax; y += 1 {
		xMax := xMin + xLen
		for x := xMin; x < xMax; x += 1 {
			intcode.In <- x
			intcode.In <- y
			result := <-intcode.Out
			grid[image.Pt(x, y)] = result
			intcode = intcode.Reset()
			go intcode.Run()
		}
	}

	for y := yMin; y < yMax; y += 1 {
		for x := xMin; x < xMax; x += 1 {
			if grid[image.Pt(x, y)] == 0 {
				continue
			}

			for y_ := y; y_ >= yMin; y_ -= 1 {
				if grid[image.Pt(x, y_)] == 0 {
					break
				}

				if y-(y_-1) != target {
					continue
				}

				for x_ := x; x_ < xMax; x_ += 1 {
					if grid[image.Pt(x_, y_)] == 0 {
						break
					}

					if x_-x+1 != target {
						continue
					}

					fmt.Printf("x: %d, y: %d\n", x, y_)
					fmt.Println(x*10000 + y_)
					return
				}
			}
		}
	}
}
