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
	matrix, graph, _ := parse(in)
	maxSteps := 64
	count := getPositionsAfterNSteps(matrix, graph, maxSteps)
	fmt.Println(count)
}

func task2(in io.Reader) {
	matrix, graph, startPos := parse(in)

	// From input:
	// S is at the center of the grid in position (65,65)
	// The grid dimension is 131x131
	// The entire row and col of the start position is empty and the grid borders are empty
	// No need for the inBounds check for task1 as 64 steps are not enough to reach the border

	// Calculate the target position for the interpolation
	// Removing 65 from the target (26_501_365) leaves a multiple of 131 (The grid size)
	target := (26_501_365 - startPos.Y) / matrix.Bounds.Max.Y

	// The task is quadratic equation (ax^2 + bx + c)
	// Lagrange's Interpolation formula with x=[0,1,2] and y=[y0,y1,y2]
	// f(x) = (x^2-3x+2) * y0/2 - (x^2-2x)*y1 + (x^2-x) * y2/2

	// f(0) = y0 = 65 + 0 * 131
	// f(1) = y1 = 65 + 1 * 131
	// f(2) = y2 = 65 + 2 * 131
	y0 := getPositionsAfterNSteps(matrix, graph, startPos.Y)
	y1 := getPositionsAfterNSteps(matrix, graph, startPos.Y+matrix.Bounds.Max.Y)
	y2 := getPositionsAfterNSteps(matrix, graph, startPos.Y+matrix.Bounds.Max.Y*2)

	// Calculate the coefficients
	// a = y0/2 - y1 + y2/2
	// b = -3*y0/2 + 2*y1 - y2/2
	// c = y0
	a := int(y0/2 - y1 + y2/2)
	b := int(-3*(y0/2) + 2*y1 - y2/2)
	c := int(y0)

	// Evaluate f(target) = ax^2 + bx + c
	fmt.Println(a*target*target + b*target + c)
}

type State struct {
	pos   image.Point
	steps int
}

func (n State) Pos() image.Point {
	return n.pos
}

func getPositionsAfterNSteps(matrix *aoc.Matrix[byte], graph *aoc.Graph[State], maxSteps int) float64 {
	const repeatableMatrix = true // The matrix is repeatable in all directions

	var possiblePositions int
	visited := map[image.Point]struct{}{graph.StartNode.pos: {}}

	graph.GetNextNodes = func(state State) (next []State) {
		next = make([]State, 0, 4)
		remainingSteps := maxSteps - state.steps

		// If the remaining steps are even, we can always return to the same position so we can count it
		if (remainingSteps)%2 == 0 {
			possiblePositions++
		}

		if state.steps >= maxSteps {
			return
		}

		for _, dir := range aoc.Directions {
			nextPos := state.pos.Add(dir)

			if matrix.Get(nextPos, repeatableMatrix) == '#' {
				continue
			}

			if _, ok := visited[nextPos]; ok {
				continue
			}

			visited[nextPos] = struct{}{}

			nextState := State{
				pos:   nextPos,
				steps: state.steps + 1,
			}

			next = append(next, nextState)
		}

		return
	}

	graph.Traverse()

	return float64(possiblePositions)
}

func parse(in io.Reader) (matrix *aoc.Matrix[byte], graph *aoc.Graph[State], startPos image.Point) {
	matrix = aoc.ReadMatrixPositionsAs[byte](in, func(ch byte, x, y int) byte {
		if ch == 'S' {
			startPos = image.Pt(x, y)
			return '.'
		}

		return ch
	})

	graph = &aoc.Graph[State]{
		StartNode: State{
			pos:   startPos,
			steps: 0,
		},
	}

	return
}
