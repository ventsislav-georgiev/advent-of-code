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
	minForward := 0
	maxForward := 3

	graph, dest := parse(in, minForward, maxForward)
	bestPath := graph.Dijkstra(func(a, b *aoc.Item[aoc.Path[State]]) bool {
		return b.Val.Node.pos.XY == dest && (a == nil || b.Val.Dist < a.Val.Dist)
	})

	fmt.Println(bestPath.Val.Dist)
}

func task2(in io.Reader) {
	minForward := 4
	maxForward := 10

	graph, dest := parse(in, minForward, maxForward)
	bestPath := graph.Dijkstra(func(a, b *aoc.Item[aoc.Path[State]]) bool {
		return b.Val.Node.pos.XY == dest && (a == nil || b.Val.Dist < a.Val.Dist && b.Val.Node.forwardSteps >= minForward)
	})

	fmt.Println(bestPath.Val.Dist)
}

type State struct {
	pos          aoc.Pos
	forwardSteps int
}

func (n State) Pos() image.Point {
	return n.pos.XY
}

func parse(in io.Reader, minForward, maxForward int) (graph *aoc.Graph[State], dest image.Point) {
	matrix := aoc.ReadMatrixPositionsAs[int](in, func(ch byte) int {
		return int(ch - '0')
	})

	graph = &aoc.Graph[State]{
		StartNode: State{
			pos:          aoc.Pos{XY: image.Pt(0, 0), Dir: aoc.DirDown},
			forwardSteps: -1,
		},
		GetNextNodes: func(state State) (next []State) {
			next = make([]State, 0, 3)

			moveForwardAllowed := state.forwardSteps < maxForward
			if moveForwardAllowed {
				next = append(next, State{
					pos:          state.pos.MoveForward(),
					forwardSteps: state.forwardSteps + 1,
				})
			}

			turnAllowed := state.forwardSteps >= minForward || state.forwardSteps == -1
			if turnAllowed {
				next = append(next, State{
					pos:          state.pos.MoveLR(aoc.DirRight),
					forwardSteps: 1,
				})
				next = append(next, State{
					pos:          state.pos.MoveLR(aoc.DirLeft),
					forwardSteps: 1,
				})
			}

			return
		},
		GetWeight: func(node State) int {
			return matrix.Map[node.pos.XY]
		},
		FilterNode: func(node State) bool {
			return node.pos.XY.In(matrix.Bounds)
		},
	}

	dest = matrix.Bounds.Max.Sub(image.Pt(1, 1))

	return
}
