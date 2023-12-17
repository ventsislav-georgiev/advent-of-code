package aoc

import "image"

var (
	DirLeft    = image.Pt(-1, 0)
	DirRight   = image.Pt(1, 0)
	DirUp      = image.Pt(0, -1)
	DirDown    = image.Pt(0, 1)
	PosNone    = image.Pt(-1, -1)
	Directions = []image.Point{DirLeft, DirRight, DirUp, DirDown}
)
