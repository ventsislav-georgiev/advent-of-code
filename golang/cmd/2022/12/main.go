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
	matrix, src, dest := parse(in, false)
	fmt.Println(BFS(matrix, src, dest))
}

func task2(in io.Reader) {
	matrix, src, dest := parse(in, true)
	fmt.Println(BFS(matrix, src, dest))
}

type Pos struct {
	row int
	col int
}

type Cell struct {
	pos  Pos
	dist int
}

func parse(in io.Reader, checkAll bool) ([][]byte, Pos, []Pos) {
	var matrix [][]byte
	var src Pos
	var dest []Pos

	row := 0
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := scanner.Bytes()
		data := make([]byte, 0, len(line))

		for col, ch := range line {
			if ch == 'S' || (checkAll && ch == 'a') {
				dest = append(dest, Pos{row, col})
				ch = 'a'
			} else if ch == 'E' {
				src = Pos{row, col}
				ch = 'z'
			}

			data = append(data, ch)
		}

		matrix = append(matrix, data)
		row += 1
	}

	return matrix, src, dest
}

func BFS(mat [][]byte, src Pos, dest []Pos) int {
	visited := make([][]bool, 0, len(mat))
	for i := 0; i < len(mat); i++ {
		visited = append(visited, make([]bool, len(mat[0])))
	}

	q := []Cell{{src, 0}}
	visited[src.row][src.col] = true
	directions := []Pos{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}

	for len(q) > 0 {
		cur := q[0]
		q = q[1:]

		pos := cur.pos
		for _, d := range dest {
			if pos.row == d.row && pos.col == d.col {
				return cur.dist
			}
		}

		for _, d := range directions {
			row := pos.row + d.row
			col := pos.col + d.col
			inbounds := row >= 0 && row < len(mat) && col >= 0 && col < len(mat[0])
			if !inbounds {
				continue
			}

			canMove := mat[pos.row][pos.col]-1 <= mat[row][col]
			if canMove && !visited[row][col] {
				visited[row][col] = true
				adjcell := Cell{Pos{row, col}, cur.dist + 1}
				q = append(q, adjcell)
			}
		}
	}

	return -1
}
