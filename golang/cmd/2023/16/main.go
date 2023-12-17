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
	dirUp = iota + 1
	dirDown
	dirLeft
	dirRight
)

func task1(in io.Reader) {
	matrix, lightPath, paths := parse(in)

	startPos := image.Pt(-1, 0)
	startDir := dirRight
	energizedTiles := calcEnergizedTiles(startPos, startDir, paths, matrix, lightPath)

	fmt.Println(energizedTiles)
}

func task2(in io.Reader) {
	matrix, lightPath, paths := parse(in)

	var maxVal int
	calc := func(startPos image.Point, startDir int) {
		resetLightPath(lightPath)

		val := calcEnergizedTiles(startPos, startDir, paths, matrix, lightPath)

		if val > maxVal {
			maxVal = val
		}
	}

	for y := 0; y < len(matrix); y++ {
		startPos := image.Pt(-1, y)
		calc(startPos, dirRight)

		startPos = image.Pt(len(matrix[0]), y)
		calc(startPos, dirLeft)
	}

	for x := 0; x < len(matrix[0]); x++ {
		startPos := image.Pt(x, -1)
		calc(startPos, dirDown)

		startPos = image.Pt(x, len(matrix))
		calc(startPos, dirUp)
	}

	fmt.Println(maxVal)
}

type Path struct {
	Pos    image.Point
	Mirror string
	Edges  []*Edge
}

type Edge struct {
	From, To *Path
	Dir      int
	Length   int
}

type CurPos struct {
	Path *Path
	Dir  int
}

func parse(in io.Reader) (matrix [][]byte, lightPath [][]bool, paths map[image.Point]*Path) {
	matrix = aoc.ReadMatrix(in)
	paths = make(map[image.Point]*Path)

	for y, row := range matrix {
		for x, ch := range row {
			if ch == '.' {
				continue
			}

			pos := image.Pt(x, y)
			paths[pos] = &Path{Pos: pos, Mirror: string(ch)}
		}
	}

	updateEdges(matrix, paths)

	lightPath = make([][]bool, 0, len(matrix))
	for _, row := range matrix {
		lightPath = append(lightPath, make([]bool, len(row)))
	}

	return
}

func updateEdges(matrix [][]byte, paths map[image.Point]*Path) {
	xLimit := len(matrix[0])
	yLimit := len(matrix)

	for _, path := range paths {
		path.Edges = make([]*Edge, 0, 4)
		var edge *Edge

		switch path.Mirror {
		case "|", "/", "\\":
			edge = &Edge{From: path, Dir: dirUp}
			edge.To, edge.Length = findPath(path.Pos, aoc.DirUp, paths, xLimit, yLimit)
			path.Edges = append(path.Edges, edge)

			edge = &Edge{From: path, Dir: dirDown}
			edge.To, edge.Length = findPath(path.Pos, aoc.DirDown, paths, xLimit, yLimit)
			path.Edges = append(path.Edges, edge)
		}

		switch path.Mirror {
		case "-", "/", "\\":
			edge = &Edge{From: path, Dir: dirLeft}
			edge.To, edge.Length = findPath(path.Pos, aoc.DirLeft, paths, xLimit, yLimit)
			path.Edges = append(path.Edges, edge)

			edge = &Edge{From: path, Dir: dirRight}
			edge.To, edge.Length = findPath(path.Pos, aoc.DirRight, paths, xLimit, yLimit)
			path.Edges = append(path.Edges, edge)
		}
	}
}

func findPath(pos image.Point, dir image.Point, paths map[image.Point]*Path, xLimit, yLimit int) (*Path, int) {
	var step int

	for {
		pos = pos.Add(dir)

		if pos.X < 0 || pos.X > xLimit || pos.Y < 0 || pos.Y > yLimit {
			return nil, step
		}

		step++

		if path, ok := paths[pos]; ok {
			return path, step
		}
	}
}

func calcEnergizedTiles(startPos image.Point, startDir int, paths map[image.Point]*Path, matrix [][]byte, lightPath [][]bool) (energizedTiles int) {
	var moveVec image.Point
	switch startDir {
	case dirUp:
		moveVec = aoc.DirUp
	case dirDown:
		moveVec = aoc.DirDown
	case dirLeft:
		moveVec = aoc.DirLeft
	case dirRight:
		moveVec = aoc.DirRight
	}

	start, startLen := findPath(startPos, moveVec, paths, len(matrix[0]), len(matrix))

	if start == nil {
		return
	}

	updateLightPath(startPos, startLen, startDir, lightPath)

	queue := []CurPos{{Path: start, Dir: startDir}}
	visited := make(map[CurPos]struct{})

	for len(queue) > 0 {
		pos := queue[0]
		queue = queue[1:]

		lightPath[pos.Path.Pos.Y][pos.Path.Pos.X] = true

		var nextDir int
		var splitLight bool

		switch pos.Path.Mirror {
		case "/":
			switch pos.Dir {
			case dirUp:
				nextDir = dirRight
			case dirDown:
				nextDir = dirLeft
			case dirLeft:
				nextDir = dirDown
			case dirRight:
				nextDir = dirUp
			}
		case "\\":
			switch pos.Dir {
			case dirUp:
				nextDir = dirLeft
			case dirDown:
				nextDir = dirRight
			case dirLeft:
				nextDir = dirUp
			case dirRight:
				nextDir = dirDown
			}
		case "|", "-":
			splitLight = true
		}

		for _, edge := range pos.Path.Edges {
			nextPos := CurPos{Path: edge.To, Dir: edge.Dir}

			if _, ok := visited[nextPos]; ok {
				continue
			}

			if edge.Dir != nextDir && !splitLight {
				continue
			}

			updateLightPath(edge.From.Pos, edge.Length, edge.Dir, lightPath)

			if edge.To == nil {
				continue
			}

			queue = append(queue, nextPos)
			visited[nextPos] = struct{}{}
		}
	}

	for _, row := range lightPath {
		for _, active := range row {
			if active {
				energizedTiles++
			}
		}
	}

	return
}

func updateLightPath(pos image.Point, length int, dir int, matrixPath [][]bool) {
	if pos.X >= 0 && pos.Y >= 0 && pos.X < len(matrixPath[0]) && pos.Y < len(matrixPath) {
		matrixPath[pos.Y][pos.X] = true
	}

	for i := 1; i <= length; i++ {
		switch dir {
		case dirUp:
			if pos.Y-i < 0 {
				break
			}
			matrixPath[pos.Y-i][pos.X] = true
		case dirDown:
			if pos.Y+i >= len(matrixPath) {
				break
			}
			matrixPath[pos.Y+i][pos.X] = true
		case dirLeft:
			if pos.X-i < 0 {
				break
			}
			matrixPath[pos.Y][pos.X-i] = true
		case dirRight:
			if pos.X+i >= len(matrixPath[0]) {
				break
			}
			matrixPath[pos.Y][pos.X+i] = true
		}
	}
}

func resetLightPath(lightPath [][]bool) {
	for _, row := range lightPath {
		for i := range row {
			row[i] = false
		}
	}
}
