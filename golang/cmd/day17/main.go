package main

import (
	"io"
	"math"

	"github.com/ventsislav-georgiev/advent-of-code-22/golang/pkg/aoc"
)

func main() {
	aoc.Exec(task1, task2)
}

func task1(in io.Reader) {
	tower, _ := solve(in, 2022)
	println(len(tower))
}

func task2(in io.Reader) {
	target := 1000000000000
	_, heights := solve(in, 10000)
	matches := 0
	repeatIdx := 1000

test:
	for ; repeatIdx < len(heights); repeatIdx++ {
		for beginIdx := 0; beginIdx <= repeatIdx; beginIdx++ {
			if heights[beginIdx] != heights[repeatIdx+beginIdx] {
				matches = 0
				continue
			}

			matches += 1
			if matches == 1000 {
				break test
			}
		}
	}

	initialHeight := sum(heights[:repeatIdx])
	repetitionHeight := sum(heights[repeatIdx : repeatIdx*2])
	includedRepeats := int(math.Floor(float64(target) / float64(repeatIdx)))
	remaining := (target - repeatIdx*includedRepeats) + 1
	println(initialHeight + repetitionHeight*(includedRepeats-1) + sum(heights[repeatIdx:repeatIdx+remaining]))
}

const (
	left  uint8 = 1
	right uint8 = 2
)

type Shape struct {
	bits   []uint8
	spaceR int
	cache  [7][]uint8
}

func (s *Shape) SpaceRight(space int) []uint8 {
	if s.cache[space] != nil {
		return s.cache[space]
	}

	shift := space - s.spaceR
	if shift == 0 {
		s.cache[space] = s.bits
		return s.bits
	}

	bits := make([]uint8, len(s.bits))
	copy(bits, s.bits)

	if shift < 0 {
		for t := 0; t < -shift; t++ {
			for i := 0; i < len(bits); i++ {
				bits[i] >>= 1
			}
		}
	} else {
		for t := 0; t < shift; t++ {
			for i := 0; i < len(bits); i++ {
				bits[i] <<= 1
			}
		}
	}
	s.cache[space] = bits
	return bits
}

func solve(in io.Reader, maxShapes int) ([]uint8, []int) {
	directions, _ := io.ReadAll(in)
	directions = directions[:len(directions)-1]

	tower := make([]uint8, 0, maxShapes*2)
	heights := make([]int, 0, maxShapes)
	minus := Shape{[]uint8{
		0b00011110,
	}, 1, [7][]uint8{}}
	plus := Shape{[]uint8{
		0b00001000,
		0b00011100,
		0b00001000,
	}, 2, [7][]uint8{}}
	l := Shape{[]uint8{
		0b00000100,
		0b00000100,
		0b00011100,
	}, 2, [7][]uint8{}}
	i := Shape{[]uint8{
		0b00010000,
		0b00010000,
		0b00010000,
		0b00010000,
	}, 4, [7][]uint8{}}
	cube := Shape{[]uint8{
		0b00011000,
		0b00011000,
	}, 3, [7][]uint8{}}
	shapes := [5]Shape{
		minus,
		plus,
		l,
		i,
		cube,
	}

	var dir uint8
	var dirIdx int
	var counter int
	var prevHeight int

loop:
	for {
		shapeData := shapes[counter%5]
		spaceL := 2
		spaceR := shapeData.spaceR
		counter++

		// remove empty rows
		for len(tower) > 0 && tower[len(tower)-1] == 0 {
			tower = tower[:len(tower)-1]
		}

		height := len(tower)
		heights = append(heights, height-prevHeight)
		prevHeight = height

		// add 3 empty rows
		tower = append(tower, make([]uint8, 3)...)

		hasShape := false
		shapeTop := 0
		shape := shapeData.bits
		movedDown, canShiftLeft, canShiftRight := false, true, true
		for {
			if directions[dirIdx] == '>' {
				dir = right
			} else {
				dir = left
			}
			dirIdx++

			if len(directions) == dirIdx {
				dirIdx = 0
			}

			if !hasShape {
				shapeTop = len(tower) + len(shape) - 1
				hasShape = true
				tower = append(tower, make([]uint8, len(shape)-1)...)
			}

			// shift shape
			if canShiftLeft && dir == left && spaceL > 0 {
				spaceL--
				spaceR++
				shape = shapeData.SpaceRight(spaceR)
			} else if canShiftRight && dir == right && spaceR > 0 {
				spaceR--
				spaceL++
				shape = shapeData.SpaceRight(spaceR)
			}

			movedDown, canShiftLeft, canShiftRight = moveDown(tower, shapeTop, shape)
			if !movedDown {
				hasShape = false
				if counter == maxShapes {
					break loop
				}

				break
			}
			shapeTop -= 1
		}
	}

	// remove empty rows
	for len(tower) > 0 && tower[len(tower)-1] == 0 {
		tower = tower[:len(tower)-1]
	}

	return tower, heights
}

func moveDown(tower []uint8, shapeTop int, shape []uint8) (movedDown bool, canShiftLeft bool, canShiftRight bool) {
	canShiftLeft = true
	canShiftRight = true

	addShape := func() {
		// add shape to tower
		for i := 0; i < len(shape); i++ {
			tower[shapeTop-i] |= shape[i]
		}
	}

	for i := 0; i < len(shape); i++ {
		tidx := shapeTop - i - 1
		if tidx < 0 {
			addShape()
			return false, false, false
		}

		if tower[tidx]&shape[i] != 0 {
			addShape()
			return false, false, false
		}

		if canShiftLeft {
			canShiftLeft = tower[tidx]&(shape[i]<<1) == 0
		}

		if canShiftRight {
			canShiftRight = tower[tidx]&(shape[i]>>1) == 0
		}
	}

	return true, canShiftLeft, canShiftRight
}

func sum(a []int) int {
	var s int
	for _, v := range a {
		s += v
	}
	return s
}
