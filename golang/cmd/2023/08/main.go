package main

import (
	"bufio"
	"fmt"
	"io"

	"github.com/ventsislav-georgiev/advent-of-code/golang/pkg/aoc"
)

func main() {
	aoc.Exec(task1, task2)
}

func task1(in io.Reader) {
	directions, tree := parse(in)
	curNode := tree["AAA"]
	endNode := tree["ZZZ"]
	var steps uint

	for {
		for _, dir := range directions {
			steps++

			switch dir {
			case 'L':
				curNode = curNode.Left
			case 'R':
				curNode = curNode.Right
			}

			if curNode == endNode {
				fmt.Println(steps)
				return
			}
		}
	}
}

func task2(in io.Reader) {
	directions, tree := parse(in)
	var totalSteps uint = 1

	for _, curNode := range tree {
		if !curNode.IsStart {
			continue
		}

		var steps uint

		for !curNode.IsExit {
			for _, dir := range directions {
				steps++

				switch dir {
				case 'L':
					curNode = curNode.Left
				case 'R':
					curNode = curNode.Right
				}

				if curNode.IsExit {
					totalSteps = aoc.LeastCommonDenominator(totalSteps, steps)
				}
			}
		}
	}

	fmt.Println(totalSteps)
}

type Node struct {
	IsStart bool
	IsExit  bool
	Left    *Node
	Right   *Node
}

func parse(in io.Reader) (directions string, tree map[string]*Node) {
	scanner := bufio.NewScanner(in)
	scanner.Scan()

	directions = string(scanner.Bytes())
	tree = map[string]*Node{}

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		cur := string(line[:3])
		left := string(line[7:10])
		right := string(line[12:15])

		var leftNode, rightNode, curNode *Node

		if node, ok := tree[left]; ok {
			leftNode = node
		} else {
			leftNode = &Node{
				IsStart: left[2] == 'A',
				IsExit:  left[2] == 'Z',
			}
			tree[left] = leftNode
		}

		if node, ok := tree[right]; ok {
			rightNode = node
		} else {
			rightNode = &Node{
				IsStart: right[2] == 'A',
				IsExit:  right[2] == 'Z',
			}
			tree[right] = rightNode
		}

		if node, ok := tree[cur]; ok {
			node.Left = leftNode
			node.Right = rightNode
		} else {
			curNode = &Node{
				IsStart: cur[2] == 'A',
				IsExit:  cur[2] == 'Z',
				Left:    leftNode,
				Right:   rightNode,
			}
			tree[cur] = curNode
		}
	}

	return
}
