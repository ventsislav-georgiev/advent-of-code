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
	m := parse(in)

	result := len(m)*2 + len(m[0])*2 - 4
	for row := 1; row < len(m)-1; row++ {
		for col := 1; col < len(m[row])-1; col++ {
			v := [4]bool{true, true, true, true}
			for i := 0; i < col; i++ {
				if m[row][i] >= m[row][col] {
					v[0] = false
					break
				} else if i+1 == col {
					goto next
				}
			}
			for i := col + 1; i < len(m[row]); i++ {
				if m[row][i] >= m[row][col] {
					v[1] = false
					break
				} else if i+1 == len(m[row]) {
					goto next
				}
			}
			for i := 0; i < row; i++ {
				if m[i][col] >= m[row][col] {
					v[2] = false
					break
				} else if i+1 == row {
					goto next
				}
			}
			for i := row + 1; i < len(m); i++ {
				if m[i][col] >= m[row][col] {
					v[3] = false
					break
				} else if i+1 == len(m) {
					goto next
				}
			}
		next:
			if v[0] || v[1] || v[2] || v[3] {
				result += 1
			}
		}
	}

	fmt.Println(result)
}

func task2(in io.Reader) {
	m := parse(in)

	result := 0
	for row := 1; row < len(m)-1; row++ {
		for col := 1; col < len(m[row])-1; col++ {
			s := [4]int{0, 0, 0, 0}
			for i := col - 1; i >= 0; i-- {
				if m[row][i] >= m[row][col] {
					s[0] = col - i
					break
				} else if i == 0 {
					s[0] = col
				}
			}
			for i := col + 1; i < len(m[row]); i++ {
				if m[row][i] >= m[row][col] {
					s[1] = i - col
					break
				} else if i+1 == len(m[row]) {
					s[1] = len(m[row]) - col - 1
				}
			}
			for i := row - 1; i >= 0; i-- {
				if m[i][col] >= m[row][col] {
					s[2] = row - i
					break
				} else if i == 0 {
					s[2] = row
				}
			}
			for i := row + 1; i < len(m); i++ {
				if m[i][col] >= m[row][col] {
					s[3] = i - row
					break
				} else if i+1 == len(m) {
					s[3] = len(m) - row - 1
				}
			}

			score := s[0] * s[1] * s[2] * s[3]
			if result < score {
				result = score
			}
		}
	}

	fmt.Println(result)
}

func parse(in io.Reader) [][]byte {
	scanner := bufio.NewScanner(in)
	var result [][]byte
	for scanner.Scan() {
		result = append(result, scanner.Bytes())
	}
	return result
}
