package aoc

import (
	"bufio"
	"fmt"
	"image"
	"io"
)

type Matrix[T any] struct {
	Rows   [][]T
	Map    map[image.Point]T
	Bounds image.Rectangle
}

type MatrixType struct {
	Rows bool
	Map  bool
}

func ReadMatrix[T any](in io.Reader, typ MatrixType, parse func(ch byte, x, y int) T) *Matrix[T] {
	scanner := bufio.NewScanner(in)

	var rows [][]T
	var posMap map[image.Point]T

	if typ.Rows {
		rows = make([][]T, 0)
	}

	if typ.Map {
		posMap = make(map[image.Point]T)
	}

	var y, maxX int
	var row []T

	if typ.Rows {
		row = make([]T, 0)
	}

	for scanner.Scan() {
		line := scanner.Bytes()

		if len(line) == 0 {
			continue
		}

		if maxX == 0 {
			maxX = len(line)
		}

		for x, ch := range line {
			val := parse(ch, x, y)

			if typ.Rows {
				row = append(row, val)
			}

			if typ.Map {
				pos := image.Pt(x, y)
				posMap[pos] = val
			}
		}

		y++

		if typ.Rows {
			rows = append(rows, row)
			row = make([]T, 0)
		}
	}

	bounds := image.Rectangle{
		Min: image.Pt(0, 0),
		Max: image.Pt(maxX, y),
	}

	return &Matrix[T]{Rows: rows, Bounds: bounds, Map: posMap}
}

var (
	readAsByte = func(ch byte, x, y int) byte {
		return ch
	}
	readAsRune = func(ch byte, x, y int) rune {
		return rune(ch)
	}
)

func ReadMatrixAs[T any](in io.Reader, parse func(ch byte, x, y int) T) *Matrix[T] {
	return ReadMatrix[T](in, MatrixType{Rows: true}, parse)
}

func ReadMatrixPositionsAs[T any](in io.Reader, parse func(ch byte, x, y int) T) *Matrix[T] {
	return ReadMatrix[T](in, MatrixType{Map: true}, parse)
}

func ReadMatrixAsBytes(in io.Reader) *Matrix[byte] {
	return ReadMatrix[byte](in, MatrixType{Rows: true}, readAsByte)
}

func ReadMatrixAsRunes(in io.Reader) *Matrix[rune] {
	return ReadMatrix[rune](in, MatrixType{Rows: true}, readAsRune)
}

func ReadMatrixPositionsAsBytes(in io.Reader) *Matrix[byte] {
	return ReadMatrix[byte](in, MatrixType{Map: true}, readAsByte)
}

func (m *Matrix[T]) Print() {
	if m.Rows == nil {
		m.printWithMap()
		return
	}

	for _, row := range m.Rows {
		for _, val := range row {
			fmt.Print(ToStr(val))
		}
		fmt.Println()
	}
}

func (m *Matrix[T]) printWithMap() {
	if m.Map == nil {
		return
	}

	for y := 0; y < m.Bounds.Max.Y; y++ {
		for x := 0; x < m.Bounds.Max.X; x++ {
			pos := image.Pt(x, y)
			val := m.Map[pos]
			fmt.Print(ToStr(val))
		}
		println()
	}
}

func (m *Matrix[T]) Get(pos image.Point, repeatable bool) T {
	var none T

	if !pos.In(m.Bounds) {
		if !repeatable {
			return none
		} else {
			pos = pos.Mod(m.Bounds)
		}
	}

	if m.Map != nil {
		return m.Map[pos]
	}

	return m.Rows[pos.Y][pos.X]
}
