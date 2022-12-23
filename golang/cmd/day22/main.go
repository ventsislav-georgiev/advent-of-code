package main

import (
	"bytes"
	"io"
	"strings"

	"github.com/ventsislav-georgiev/advent-of-code-22/golang/pkg/aoc"
)

var boardString = `    ........
    ........
    ........
    ........
    ....    
    ....    
    ....    
    ....    
........    
........    
........    
........    
....        
....        
....        
....        

1L2L1L3R5L1L6R3L1L4R5R7L3`

func main() {
	aoc.Exec(task1, task2)
}

var testMode = false

func task1(in io.Reader) {
	if testMode {
		in = strings.NewReader(boardString)
	}

	board, cmds, pos := parse(in)
	for _, cmd := range cmds {
		pos, _ = walk(board, pos, cmd, false)
	}

	if testMode {
		printBoard(board, cmds)
	}

	println(calcAnswer(pos))
}

func task2(in io.Reader) {
	if testMode {
		in = strings.NewReader(boardString)
	}

	board, cmds, pos := parse(in)
	initCubeSides(board)
	for _, cmd := range cmds {
		pos, _ = walk(board, pos, cmd, true)
	}

	if testMode {
		printBoard(board, cmds)
	}

	println(calcAnswer(pos))
}

func calcAnswer(pos Vector) int {
	facing := (pos.dirIndex - 2 + 4) % 4
	println(pos.cords.y+1, pos.cords.x+1, facing)
	return 1000*(pos.cords.y+1) + 4*(pos.cords.x+1) + facing
}

func walk(board [][]byte, pos Vector, cmd Command, cubeLike bool) (Vector, bool) {
	switch cmd.typ {
	case move:
		for i := 0; i < cmd.steps; i++ {
			warp := warpBounds
			if cubeLike {
				warp = warpCube
			}

			tmpx, tmpy, rotate, ok := warp(board, pos)
			if !ok {
				return pos, false
			}

			if rotate != nil {
				pos = rotate(pos)
			}

			switch pos.dirIndex {
			case 0:
				board[tmpy][tmpx] = '<'
			case 1:
				board[tmpy][tmpx] = '^'
			case 2:
				board[tmpy][tmpx] = '>'
			case 3:
				board[tmpy][tmpx] = 'v'
			}

			pos.cords.x = tmpx
			pos.cords.y = tmpy
		}
	case rotate:
		pos.dirIndex = (pos.dirIndex + cmd.rotate + 4) % 4
		pos.dir = &directions[pos.dirIndex]
		switch pos.dirIndex {
		case 0:
			board[pos.cords.y][pos.cords.x] = '<'
		case 1:
			board[pos.cords.y][pos.cords.x] = '^'
		case 2:
			board[pos.cords.y][pos.cords.x] = '>'
		case 3:
			board[pos.cords.y][pos.cords.x] = 'v'
		}
	}

	return pos, true
}

type Vector struct {
	cords    *Point
	dir      *Point
	dirIndex int
}

type Point struct {
	x, y int
}

type CommandType uint8

const (
	move CommandType = iota
	rotate
)

type Command struct {
	typ    CommandType
	steps  int
	rotate int
}

var (
	dirLeft  = Point{-1, 0}
	dirUp    = Point{0, -1}
	dirRight = Point{1, 0}
	dirDown  = Point{0, 1}
)

var directions = [4]Point{
	dirLeft,
	dirUp,
	dirRight,
	dirDown,
}

func parse(in io.Reader) ([][]byte, []Command, Vector) {
	var board [][]byte
	pos := Vector{
		cords:    &Point{},
		dir:      &directions[2],
		dirIndex: 2,
	}

	input, _ := io.ReadAll(in)
	lines := bytes.Split(input, []byte{'\n'})
	var path []byte
	longestLine := 0

	for i, line := range lines {
		if len(line) > longestLine {
			longestLine = len(line)
		}

		if len(line) == 0 {
			path = lines[i+1]
			break
		}

		board = append(board, line)
		if len(board) == 1 {
			pos.cords.x = bytes.IndexByte(line, '.')
		}
	}

	for i := range board {
		board[i] = append(board[i], bytes.Repeat([]byte{' '}, longestLine-len(board[i]))...)
	}

	cmds := make([]Command, 0, len(path))

	n := 0
	for _, ch := range path {
		switch ch {
		case 'L':
			cmds = append(cmds, Command{typ: move, steps: n})
			n = 0
			cmds = append(cmds, Command{typ: rotate, rotate: -1})
		case 'R':
			cmds = append(cmds, Command{typ: move, steps: n})
			n = 0
			cmds = append(cmds, Command{typ: rotate, rotate: 1})
		default:
			n = n*10 + int(ch-'0')
		}
	}

	if n > 0 {
		cmds = append(cmds, Command{typ: move, steps: n})
	}

	return board, cmds, pos
}

func printBoard(board [][]byte, cmds []Command) {
	for _, line := range board {
		println(string(line))
	}

	println()

	if len(cmds) == 0 {
		return
	}

	for _, cmd := range cmds {
		switch cmd.typ {
		case move:
			print(cmd.steps)
		case rotate:
			if cmd.rotate < 0 {
				print("L")
			} else {
				print("R")
			}
		}
	}

	println()
}

// ################################################################ PART1

func warpBounds(board [][]byte, pos Vector) (int, int, func(pos Vector) Vector, bool) {
	leftBounds, upBounds, rightBounds, downBounds := getPositionBounds(board, pos)
	tmpx := pos.cords.x + pos.dir.x
	tmpy := pos.cords.y + pos.dir.y

	if tmpx < leftBounds {
		tmpx = rightBounds
	} else if tmpx > rightBounds {
		tmpx = leftBounds
	} else if tmpy < upBounds {
		tmpy = downBounds
	} else if tmpy > downBounds {
		tmpy = upBounds
	}

	if board[tmpy][tmpx] == '#' {
		return 0, 0, nil, false
	}

	return tmpx, tmpy, nil, true
}

var boundsCache = map[uint64]*[4]int{}

func getPositionBounds(board [][]byte, pos Vector) (int, int, int, int) {
	var leftBounds, upBounds, rightBounds, downBounds int
	posKey := uint64(pos.cords.y)<<32 + uint64(pos.cords.x)
	if v, ok := boundsCache[posKey]; ok {
		return v[0], v[1], v[2], v[3]
	}

	bounds := [4]int{0, 0, 0, 0}
	boundsCache[posKey] = &bounds

	for i := 0; i < len(board[0]); i++ {
		if board[pos.cords.y][i] != ' ' {
			leftBounds = i
			bounds[0] = leftBounds
			break
		}
	}
	for i := 0; i < len(board); i++ {
		if board[i][pos.cords.x] != ' ' {
			upBounds = i
			bounds[1] = upBounds
			break
		}
	}
	for i := len(board[0]) - 1; i > 0; i-- {
		if board[pos.cords.y][i] != ' ' {
			rightBounds = i
			bounds[2] = rightBounds
			break
		}
	}
	for i := len(board) - 1; i > 0; i-- {
		if board[i][pos.cords.x] != ' ' {
			downBounds = i
			bounds[3] = downBounds
			break
		}
	}

	return leftBounds, upBounds, rightBounds, downBounds
}

// ################################################################ PART2

func warpCube(board [][]byte, pos Vector) (int, int, func(pos Vector) Vector, bool) {
	var curSide *CubeSide
	for _, side := range boardSides {
		if side.bounds.minY <= pos.cords.y && pos.cords.y <= side.bounds.maxY && side.bounds.minX <= pos.cords.x && pos.cords.x <= side.bounds.maxX {
			curSide = side
			break
		}
	}

	leftBounds, upBounds, rightBounds, downBounds := curSide.bounds.minX, curSide.bounds.minY, curSide.bounds.maxX, curSide.bounds.maxY

	tmpx := pos.cords.x + pos.dir.x
	tmpy := pos.cords.y + pos.dir.y
	dir := *pos.dir

	if tmpx < leftBounds {
		tmpx, tmpy, dir = curSide.getPos(pos, dirLeft, tmpx, tmpy)
	} else if tmpx > rightBounds {
		tmpx, tmpy, dir = curSide.getPos(pos, dirRight, tmpx, tmpy)
	} else if tmpy < upBounds {
		tmpx, tmpy, dir = curSide.getPos(pos, dirUp, tmpx, tmpy)
	} else if tmpy > downBounds {
		tmpx, tmpy, dir = curSide.getPos(pos, dirDown, tmpx, tmpy)
	}

	if board[tmpy][tmpx] == '#' {
		return 0, 0, nil, false
	}

	rotate := func(pos Vector) Vector {
		for dirIndex, d := range directions {
			if d == dir {
				pos.dirIndex = dirIndex
				pos.dir = &directions[dirIndex]
				break
			}
		}
		return pos
	}

	return tmpx, tmpy, rotate, true
}

type CubeSide struct {
	bounds *Bounds
	left   *CubeSideFace
	right  *CubeSideFace
	up     *CubeSideFace
	down   *CubeSideFace
}

func (s CubeSide) getPos(pos Vector, dir Point, tmpx, tmpy int) (int, int, Point) {
	var newDir Point

	switch dir {
	case dirLeft:
		newDir = s.left.dir
		shift := pos.cords.y - s.bounds.minY
		switch newDir {
		case dirLeft:
			tmpx = s.left.side.bounds.maxX
			tmpy = s.left.side.bounds.minY + shift
		case dirRight:
			tmpx = s.left.side.bounds.minX
			tmpy = s.left.side.bounds.maxY - shift
		case dirUp:
			tmpy = s.left.side.bounds.maxY
			tmpx = s.left.side.bounds.minX + shift
		case dirDown:
			tmpy = s.left.side.bounds.minY
			tmpx = s.left.side.bounds.minX + shift
		}
	case dirRight:
		newDir = s.right.dir
		shift := pos.cords.y - s.bounds.minY
		switch newDir {
		case dirRight:
			tmpx = s.right.side.bounds.minX
			tmpy = s.right.side.bounds.minY + shift
		case dirLeft:
			tmpx = s.right.side.bounds.maxX
			tmpy = s.right.side.bounds.maxY - shift
		case dirUp:
			tmpy = s.right.side.bounds.maxY
			tmpx = s.right.side.bounds.minX + shift
		case dirDown:
			tmpy = s.right.side.bounds.minY
			tmpx = s.right.side.bounds.minX + shift
		}
	case dirUp:
		newDir = s.up.dir
		shift := pos.cords.x - s.bounds.minX
		switch newDir {
		case dirUp:
			tmpy = s.up.side.bounds.maxY
			tmpx = s.up.side.bounds.minX + shift
		case dirDown:
			tmpy = s.up.side.bounds.minY
			tmpx = s.up.side.bounds.minX + shift
		case dirLeft:
			tmpx = s.up.side.bounds.maxX
			tmpy = s.up.side.bounds.minY + shift
		case dirRight:
			tmpx = s.up.side.bounds.minX
			tmpy = s.up.side.bounds.minY + shift
		}
	case dirDown:
		newDir = s.down.dir
		shift := pos.cords.x - s.bounds.minX
		switch newDir {
		case dirDown:
			tmpy = s.down.side.bounds.minY
			tmpx = s.down.side.bounds.minX + shift
		case dirUp:
			tmpy = s.down.side.bounds.maxY
			tmpx = s.down.side.bounds.minX + shift
		case dirLeft:
			tmpx = s.down.side.bounds.maxX
			tmpy = s.down.side.bounds.minY + shift
		case dirRight:
			tmpx = s.down.side.bounds.minX
			tmpy = s.down.side.bounds.minY + shift
		}
	}

	return tmpx, tmpy, newDir
}

type CubeSideFace struct {
	side *CubeSide
	dir  Point
}

type Bounds struct {
	minX, maxX, minY, maxY int
}

var (
	base  = &CubeSide{}
	front = &CubeSide{}
	back  = &CubeSide{}
	left  = &CubeSide{}
	right = &CubeSide{}
	top   = &CubeSide{}
)

var boardSides []*CubeSide
var cubeSize int

func initCubeSides(board [][]byte) {
	boardSides = []*CubeSide{front, right, base, left, back, top}
	if testMode {
		cubeSize = 4
	} else {
		cubeSize = 50
	}

	front.left = &CubeSideFace{left, dirRight}
	front.up = &CubeSideFace{top, dirRight}
	front.right = &CubeSideFace{right, dirRight}
	front.down = &CubeSideFace{base, dirDown}

	top.left = &CubeSideFace{front, dirDown}
	top.up = &CubeSideFace{left, dirUp}
	top.right = &CubeSideFace{back, dirUp}
	top.down = &CubeSideFace{right, dirDown}

	left.left = &CubeSideFace{front, dirRight}
	left.up = &CubeSideFace{base, dirRight}
	left.right = &CubeSideFace{back, dirRight}
	left.down = &CubeSideFace{top, dirDown}

	base.left = &CubeSideFace{left, dirDown}
	base.up = &CubeSideFace{front, dirUp}
	base.right = &CubeSideFace{right, dirUp}
	base.down = &CubeSideFace{back, dirDown}

	back.left = &CubeSideFace{left, dirLeft}
	back.up = &CubeSideFace{base, dirUp}
	back.right = &CubeSideFace{right, dirLeft}
	back.down = &CubeSideFace{top, dirLeft}

	right.left = &CubeSideFace{front, dirLeft}
	right.up = &CubeSideFace{top, dirUp}
	right.right = &CubeSideFace{back, dirLeft}
	right.down = &CubeSideFace{base, dirLeft}

	sideIdx := 0
	for y := 0; y < len(board); y += cubeSize {
		row := board[y]
		for x := 0; x < len(row); x++ {
			cell := row[x]
			if cell == ' ' {
				x += cubeSize - 1
				continue
			}
			boardSides[sideIdx].bounds = &Bounds{
				minX: x,
				maxX: x + cubeSize - 1,
				minY: y,
				maxY: y + cubeSize - 1,
			}
			sideIdx++
			x += cubeSize - 1
		}
	}
}
