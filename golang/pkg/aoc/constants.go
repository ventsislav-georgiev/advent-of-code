package aoc

import "image"

var (
	MoveLeft   = image.Pt(-1, 0)
	MoveRight  = image.Pt(1, 0)
	MoveUp     = image.Pt(0, -1)
	MoveDown   = image.Pt(0, 1)
	PosNone    = image.Pt(-1, -1)
	Directions = []image.Point{MoveLeft, MoveRight, MoveUp, MoveDown}
)

const (
	DirUp = iota + 1
	DirDown
	DirLeft
	DirRight
)
