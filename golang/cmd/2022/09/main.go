package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"

	"github.com/ventsislav-georgiev/advent-of-code/golang/pkg/aoc"
)

func main() {
	aoc.Exec(task1, task2)
}

func task1(in io.Reader) {
	fmt.Println(solve(in, 2))
}

func task2(in io.Reader) {
	fmt.Println(solve(in, 10))
}

func solve(in io.Reader, length int) int {
	path := map[string]struct{}{}
	snake := make([][]int, 0, length)
	for i := 0; i < length; i++ {
		snake = append(snake, []int{0, 0})
	}

	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		cmd := strings.Fields(scanner.Text())
		n, _ := strconv.Atoi(cmd[1])
		for i := 0; i < n; i++ {
			switch cmd[0] {
			case "L":
				snake[0][0] -= 1
			case "R":
				snake[0][0] += 1
			case "U":
				snake[0][1] += 1
			case "D":
				snake[0][1] -= 1
			}

			for i := 0; i < len(snake)-1; i++ {
				prev := snake[i]
				next := snake[i+1]
				if distance(prev, next) <= 1 {
					break
				}
				if prev[0] > next[0] {
					next[0] += 1
				} else if prev[0] < next[0] {
					next[0] -= 1
				}
				if prev[1] > next[1] {
					next[1] += 1
				} else if prev[1] < next[1] {
					next[1] -= 1
				}
			}

			path[fmt.Sprintf("%d,%d", snake[length-1][0], snake[length-1][1])] = struct{}{}
		}
	}

	return len(path)
}

func distance(vec1, vec2 []int) int {
	xdiff := vec1[0] - vec2[0]
	ydiff := vec1[1] - vec2[1]
	return int(math.Sqrt(float64(xdiff*xdiff + ydiff*ydiff)))
}
