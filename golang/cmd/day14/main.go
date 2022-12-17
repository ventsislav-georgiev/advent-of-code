package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math"

	"github.com/ventsislav-georgiev/advent-of-code-22/golang/pkg/aoc"
)

func main() {
	aoc.Exec(task1, task2)
}

func task1(in io.Reader) {
	cave, bounds := parse(in)
	solve(cave, bounds, false)
	printCave(cave, bounds)
}

func task2(in io.Reader) {
	cave, bounds := parse(in)

	bounds[1][1] += 2
	bounds[0][0] = bounds[0][0] - bounds[1][1]*2
	bounds[1][0] = bounds[1][0] + bounds[1][1]*2
	for x := bounds[0][0]; x <= bounds[1][0]; x++ {
		cave[toUint32(x, bounds[1][1])] = rock
	}

	solve(cave, bounds, true)
}

const (
	empty uint8 = 0
	rock  uint8 = 254
	sand  uint8 = 255
)

func solve(cave map[uint32]uint8, bounds [2][2]int16, toSource bool) {
	x, y := int16(500), int16(-1)
	result := 0

	for {
		if cave[toUint32(x, y+1)] != empty {
			// move left
			if cave[toUint32(x-1, y+1)] == empty {
				x -= 1
				y += 1
				continue
			}

			// move right
			if cave[toUint32(x+1, y+1)] == empty {
				x += 1
				y += 1
				continue
			}

			result += 1
			cave[toUint32(x, y)] = sand
			if toSource && x == 500 && y == 0 {
				break
			}
			x, y = int16(500), int16(0)
		} else {
			// move down
			y += 1
			if !toSource && y > bounds[1][1] {
				break
			}
		}
	}

	fmt.Println(result)
}

func parse(in io.Reader) (map[uint32]uint8, [2][2]int16) {
	cave := map[uint32]uint8{}
	bounds := [2][2]int16{{math.MaxInt16, math.MaxInt16}, {0, 0}}
	scanner := bufio.NewScanner(in)

	for scanner.Scan() {
		var prev [2]int16
		cords := bufio.NewScanner(bytes.NewBuffer(scanner.Bytes()))
		cords.Split(bufio.ScanWords)

		for cords.Scan() {
			xy := cords.Bytes()
			if xy[0] == '-' {
				continue
			}

			sep := bytes.Index(xy, []byte(","))
			x := int16(aoc.Atoui(xy[:sep]))
			y := int16(aoc.Atoui(xy[sep+1:]))
			cave[toUint32(x, y)] = rock

			if prev[0] == 0 {
				prev[0] = x
				prev[1] = y
				continue
			}

			if x > prev[0] {
				for px := prev[0] + 1; px < x; px++ {
					cave[toUint32(px, y)] = rock
				}
			} else if x < prev[0] {
				for px := x + 1; px < prev[0]; px++ {
					cave[toUint32(px, y)] = rock
				}
			}

			if y > prev[1] {
				for py := prev[1] + 1; py < y; py++ {
					cave[toUint32(x, py)] = rock
				}
			} else if y < prev[1] {
				for py := y + 1; py < prev[1]; py++ {
					cave[toUint32(x, py)] = rock
				}
			}

			prev[0] = x
			prev[1] = y

			bounds[0][0] = int16(math.Min(float64(bounds[0][0]), float64(x)))
			bounds[0][1] = int16(math.Min(float64(bounds[0][1]), float64(y)))
			bounds[1][0] = int16(math.Max(float64(bounds[1][0]), float64(x)))
			bounds[1][1] = int16(math.Max(float64(bounds[1][1]), float64(y)))
		}
	}

	return cave, bounds
}

func printCave(cave map[uint32]uint8, bounds [2][2]int16) {
	for y := int16(0); y <= bounds[1][1]; y++ {
		for x := bounds[0][0]; x <= bounds[1][0]; x++ {
			el, ok := cave[toUint32(int16(x), int16(y))]
			if !ok {
				fmt.Print(".")
				continue
			}

			switch el {
			case rock:
				fmt.Print("#")
			case sand:
				fmt.Print("O")
			case empty:
				if x == 500 && y == 0 {
					fmt.Print("+")
				}
			}
		}
		fmt.Println()
	}
}

func toUint32(x, y int16) uint32 {
	return uint32(x)<<16 | uint32(y)
}
