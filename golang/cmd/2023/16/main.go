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
	_, lightPath, paths := parse(in)

	startPos := image.Pt(-1, 0)
	startDir := aoc.DirRight
	energizedTiles := calcEnergizedTiles(startPos, startDir, paths, lightPath)

	fmt.Println(energizedTiles)
}

func task2(in io.Reader) {
	matrix, lightPath, paths := parse(in)

	var maxVal int
	calc := func(startPos image.Point, startDir int) {
		resetLightPath(lightPath)

		val := calcEnergizedTiles(startPos, startDir, paths, lightPath)

		if val > maxVal {
			maxVal = val
		}
	}

	for y := 0; y < matrix.Bounds.Max.Y; y++ {
		startPos := image.Pt(-1, y)
		calc(startPos, aoc.DirRight)

		startPos = image.Pt(matrix.Bounds.Max.X, y)
		calc(startPos, aoc.DirLeft)
	}

	for x := 0; x < matrix.Bounds.Max.X; x++ {
		startPos := image.Pt(x, -1)
		calc(startPos, aoc.DirDown)

		startPos = image.Pt(x, matrix.Bounds.Max.Y)
		calc(startPos, aoc.DirUp)
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

func parse(in io.Reader) (matrix *aoc.Matrix[byte], lightPath [][]bool, paths map[image.Point]*Path) {
	matrix = aoc.ReadMatrixAsBytes(in)
	paths = make(map[image.Point]*Path)

	for y, row := range matrix.Rows {
		for x, ch := range row {
			if ch == '.' {
				continue
			}

			pos := image.Pt(x, y)
			paths[pos] = &Path{Pos: pos, Mirror: string(ch)}
		}
	}

	updateEdges(matrix, paths)

	lightPath = make([][]bool, 0, matrix.Bounds.Max.X)
	for _, row := range matrix.Rows {
		lightPath = append(lightPath, make([]bool, len(row)))
	}

	return
}

func updateEdges(matrix *aoc.Matrix[byte], paths map[image.Point]*Path) {
	xMax := matrix.Bounds.Max.X
	yMax := matrix.Bounds.Max.Y

	for _, path := range paths {
		path.Edges = make([]*Edge, 0, 4)
		var edge *Edge

		switch path.Mirror {
		case "|", "/", "\\":
			edge = &Edge{From: path, Dir: aoc.DirUp}
			edge.To, edge.Length = findPath(path.Pos, aoc.MoveUp, paths, xMax, yMax)
			path.Edges = append(path.Edges, edge)

			edge = &Edge{From: path, Dir: aoc.DirDown}
			edge.To, edge.Length = findPath(path.Pos, aoc.MoveDown, paths, xMax, yMax)
			path.Edges = append(path.Edges, edge)
		}

		switch path.Mirror {
		case "-", "/", "\\":
			edge = &Edge{From: path, Dir: aoc.DirLeft}
			edge.To, edge.Length = findPath(path.Pos, aoc.MoveLeft, paths, xMax, yMax)
			path.Edges = append(path.Edges, edge)

			edge = &Edge{From: path, Dir: aoc.DirRight}
			edge.To, edge.Length = findPath(path.Pos, aoc.MoveRight, paths, xMax, yMax)
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

func calcEnergizedTiles(startPos image.Point, startDir int, paths map[image.Point]*Path, lightPath [][]bool) (energizedTiles int) {
	var moveVec image.Point
	switch startDir {
	case aoc.DirUp:
		moveVec = aoc.MoveUp
	case aoc.DirDown:
		moveVec = aoc.MoveDown
	case aoc.DirLeft:
		moveVec = aoc.MoveLeft
	case aoc.DirRight:
		moveVec = aoc.MoveRight
	}

	start, startLen := findPath(startPos, moveVec, paths, len(lightPath[0]), len(lightPath))

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
			case aoc.DirUp:
				nextDir = aoc.DirRight
			case aoc.DirDown:
				nextDir = aoc.DirLeft
			case aoc.DirLeft:
				nextDir = aoc.DirDown
			case aoc.DirRight:
				nextDir = aoc.DirUp
			}
		case "\\":
			switch pos.Dir {
			case aoc.DirUp:
				nextDir = aoc.DirLeft
			case aoc.DirDown:
				nextDir = aoc.DirRight
			case aoc.DirLeft:
				nextDir = aoc.DirUp
			case aoc.DirRight:
				nextDir = aoc.DirDown
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
		case aoc.DirUp:
			if pos.Y-i < 0 {
				break
			}
			matrixPath[pos.Y-i][pos.X] = true
		case aoc.DirDown:
			if pos.Y+i >= len(matrixPath) {
				break
			}
			matrixPath[pos.Y+i][pos.X] = true
		case aoc.DirLeft:
			if pos.X-i < 0 {
				break
			}
			matrixPath[pos.Y][pos.X-i] = true
		case aoc.DirRight:
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
