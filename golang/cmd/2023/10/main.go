package main

import (
	"bufio"
	"image"
	"io"

	"github.com/ventsislav-georgiev/advent-of-code/golang/pkg/aoc"
)

func main() {
	aoc.Exec(task1, task2)
}

func task1(in io.Reader) {
	startNode, _ := parse(in)

	var curNode, prevNode *Node
	curNode = startNode
	var dist uint

	for {
		if curNode == startNode && dist > 0 {
			break
		}

		dist++
		nextNode := curNode.MoveNext(prevNode)
		prevNode = curNode
		curNode = nextNode
	}

	println(dist / 2)
}

func task2(in io.Reader) {
	startNode, grid := parse(in)

	var curNode, prevNode *Node
	curNode = startNode
	loop := []image.Point{}
	loopCache := map[image.Point]struct{}{}

	for {
		if curNode == startNode && len(loop) > 0 {
			break
		}

		nextNode := curNode.MoveNext(prevNode)
		prevNode = curNode
		loop = append(loop, curNode.Point)
		loopCache[curNode.Point] = struct{}{}
		curNode = nextNode
	}

	enclosedTilesCount := 0
	for _, node := range grid {
		if _, ok := loopCache[node.Point]; !ok {
			if aoc.IsPointInsideLoop(node.Point, loop) {
				enclosedTilesCount++
			}
		}
	}

	println(enclosedTilesCount)
}

type Node struct {
	Pipe  byte
	Point image.Point
	West  *Node
	East  *Node
	North *Node
	South *Node
}

func (n *Node) MoveNext(prevNode *Node) *Node {
	var nextNode *Node

	if n.North != nil && n.North != prevNode {
		nextNode = n.North
	} else if n.South != nil && n.South != prevNode {
		nextNode = n.South
	} else if n.West != nil && n.West != prevNode {
		nextNode = n.West
	} else if n.East != nil && n.East != prevNode {
		nextNode = n.East
	}

	return nextNode
}

func parse(in io.Reader) (startNode *Node, grid map[uint64]*Node) {
	grid = make(map[uint64]*Node)

	scanner := bufio.NewScanner(in)
	var y, x int
	for scanner.Scan() {
		line := scanner.Bytes()

		x = 0
		for _, ch := range line {
			grid[aoc.ToKey(x, y)] = &Node{Pipe: ch, Point: image.Point{x, y}}
			x++
		}

		y++
	}

	rows := y
	cols := x
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			node := grid[aoc.ToKey(x, y)]

			switch node.Pipe {
			case '-':
				node.West = grid[aoc.ToKey(x-1, y)]
				node.East = grid[aoc.ToKey(x+1, y)]
			case '|':
				node.North = grid[aoc.ToKey(x, y-1)]
				node.South = grid[aoc.ToKey(x, y+1)]
			case 'L':
				node.North = grid[aoc.ToKey(x, y-1)]
				node.East = grid[aoc.ToKey(x+1, y)]
			case 'J':
				node.North = grid[aoc.ToKey(x, y-1)]
				node.West = grid[aoc.ToKey(x-1, y)]
			case 'F':
				node.South = grid[aoc.ToKey(x, y+1)]
				node.East = grid[aoc.ToKey(x+1, y)]
			case '7':
				node.South = grid[aoc.ToKey(x, y+1)]
				node.West = grid[aoc.ToKey(x-1, y)]
			case 'S':
				node.North = grid[aoc.ToKey(x, y-1)]
				node.South = grid[aoc.ToKey(x, y+1)]
				node.East = grid[aoc.ToKey(x+1, y)]
				node.West = grid[aoc.ToKey(x-1, y)]
				startNode = node
			}
		}
	}

	if startNode.North != nil && startNode.North.South == nil {
		startNode.North = nil
	}
	if startNode.South != nil && startNode.South.North == nil {
		startNode.South = nil
	}
	if startNode.East != nil && startNode.East.West == nil {
		startNode.East = nil
	}
	if startNode.West != nil && startNode.West.East == nil {
		startNode.West = nil
	}

	return
}
