package main

import (
	"bufio"
	"io"
	"strings"

	"github.com/ventsislav-georgiev/advent-of-code-22/golang/pkg/aoc"
)

var testInput = `..............
..............
.......#......
.....###.#....
...#...#.#....
....#...##....
...#.###......
...##.#.##....
....#..#......
..............
..............
..............`

var testMode = false

func main() {
	aoc.Exec(task1, task2)
}

func task1(in io.Reader) {
	if testMode {
		in = strings.NewReader(testInput)
	}

	field := parse(in)
	if testMode {
		scanField(field, 0)
	}

	for i := 0; i < 10; i++ {
		field.Round(i)
		if testMode && i+1 < 10 {
			scanField(field, i+1)
		}
	}

	empty := scanField(field, 10)
	println(empty)
}

func task2(in io.Reader) {
	if testMode {
		in = strings.NewReader(testInput)
	}

	field := parse(in)
	if testMode {
		scanField(field, 0)
	}

	moved := true
	round := 0
	for moved {
		moved = field.Round(round)
		if testMode && round+1 < 10 {
			scanField(field, round+1)
		}
		round++
	}

	println(round)
}

var (
	north = [2]int{0, -1}
	south = [2]int{0, 1}
	west  = [2]int{-1, 0}
	east  = [2]int{1, 0}
)

type Field struct {
	elfs       []*Elf
	grid       map[uint64]struct{}
	bounds     FieldBounds
	directions [4][2]int
}

func (f *Field) Round(roundIdx int) bool {
	dirIndex := roundIdx % 4
	moved := false

	// phase 1
	for _, elf := range f.elfs {
		x, y := elf.xy[0], elf.xy[1]

		nx := x - 1
		sx := x + 1
		wy := y - 1
		ey := y + 1

		_, n := f.grid[toKey(nx, y)]
		_, nw := f.grid[toKey(nx, wy)]
		_, ne := f.grid[toKey(nx, ey)]
		_, s := f.grid[toKey(sx, y)]
		_, sw := f.grid[toKey(sx, wy)]
		_, se := f.grid[toKey(sx, ey)]
		_, w := f.grid[toKey(x, wy)]
		_, e := f.grid[toKey(x, ey)]
		if !n && !nw && !ne && !s && !sw && !se && !w && !e {
			continue
		}

		for i := 0; i < 4; i++ {
			dir := f.directions[(dirIndex+i)%4]
			nextXY := [2]int{x + dir[0], y + dir[1]}
			x, y := nextXY[0], nextXY[1]

			if _, ok := f.grid[toKey(x, y)]; ok {
				continue
			}

			if dir == north || dir == south {
				_, w := f.grid[toKey(x-1, y)]
				_, e := f.grid[toKey(x+1, y)]
				if w || e {
					continue
				}
			} else {
				_, n := f.grid[toKey(x, y-1)]
				_, s := f.grid[toKey(x, y+1)]
				if n || s {
					continue
				}
			}

			elf.nextXY = nextXY
			elf.hasNext = true
			break
		}
	}

	// phase 2
	for _, elf := range f.elfs {
		if !elf.hasNext {
			continue
		}

		skip := false
		for _, other := range f.elfs {
			if other == elf || !other.hasNext || other.nextXY[0] != elf.nextXY[0] || other.nextXY[1] != elf.nextXY[1] {
				continue
			}

			skip = true
			other.hasNext = false
		}

		if skip {
			elf.hasNext = false
			continue
		}

		moved = true
		delete(f.grid, toKeyXY(elf.xy))
		elf.xy = elf.nextXY
		elf.hasNext = false
		f.grid[toKeyXY(elf.xy)] = struct{}{}
		updateBounds(f, elf)
	}

	return moved
}

type FieldBounds struct {
	minX, maxX, minY, maxY *Elf
}

type Elf struct {
	nextXY  [2]int
	hasNext bool
	xy      [2]int
}

func parse(in io.Reader) *Field {
	scanner := bufio.NewScanner(in)
	field := Field{
		elfs:       []*Elf{},
		grid:       map[uint64]struct{}{},
		bounds:     FieldBounds{},
		directions: [4][2]int{north, south, west, east},
	}

	var y int
	for scanner.Scan() {
		line := scanner.Bytes()
		for x, c := range line {
			switch c {
			case '#':
				elf := Elf{
					xy: [2]int{x, y},
				}
				field.grid[toKeyXY(elf.xy)] = struct{}{}
				field.elfs = append(field.elfs, &elf)
				updateBounds(&field, &elf)
			}
		}
		y++
	}

	return &field
}

func updateBounds(field *Field, elf *Elf) {
	if field.bounds.minX == nil || field.bounds.minX.xy[0] > elf.xy[0] {
		field.bounds.minX = elf
	}
	if field.bounds.maxX == nil || field.bounds.maxX.xy[0] < elf.xy[0] {
		field.bounds.maxX = elf
	}
	if field.bounds.minY == nil || field.bounds.minY.xy[1] > elf.xy[1] {
		field.bounds.minY = elf
	}
	if field.bounds.maxY == nil || field.bounds.maxY.xy[1] < elf.xy[1] {
		field.bounds.maxY = elf
	}
}

func toKeyXY(xy [2]int) uint64 {
	return toKey(xy[0], xy[1])
}

func toKey(x, y int) uint64 {
	return uint64(x)<<16 + uint64(y)
}

func scanField(field *Field, round int) int {
	empty := 0

	if testMode {
		println("Round:", round)
	}

	for y := field.bounds.minY.xy[1]; y <= field.bounds.maxY.xy[1]; y++ {
		for x := field.bounds.minX.xy[0]; x <= field.bounds.maxX.xy[0]; x++ {
			if _, ok := field.grid[toKey(x, y)]; ok {
				if testMode {
					print("#")
				}
			} else {
				if testMode {
					print(".")
				}
				empty++
			}
		}

		if testMode {
			println()
		}
	}

	if testMode {
		println()
	}

	return empty
}
