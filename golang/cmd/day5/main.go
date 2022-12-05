package main

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"

	"github.com/ventsislav-georgiev/advent-of-code-22/golang/pkg/aoc"
)

func main() {
	aoc.Exec(task1, task2)
}

var regex *regexp.Regexp = regexp.MustCompile(`move (\d+) from (\d+) to (\d+)`)

func task1(in io.ReadCloser) {
	scanner := bufio.NewScanner(in)
	stacks := parseStacks(scanner)
	for scanner.Scan() {
		from, to, count := parseCommand(scanner)
		rearrange(stacks, from, to, count, false)
	}
	print(stacks)
}

func task2(in io.ReadCloser) {
	scanner := bufio.NewScanner(in)
	stacks := parseStacks(scanner)
	for scanner.Scan() {
		from, to, count := parseCommand(scanner)
		rearrange(stacks, from, to, count, true)
	}
	print(stacks)
}

func print(stacks [][]byte) {
	for _, stack := range stacks {
		fmt.Print(string(stack[len(stack)-1]))
	}
}

func rearrange(stacks [][]byte, from, to, count int, multi bool) {
	if multi {
		topindex := len(stacks[from])
		stacks[to] = append(stacks[to], stacks[from][topindex-count:]...)
		stacks[from] = stacks[from][:topindex-count]
		return
	}

	for i := 0; i < count; i++ {
		topindex := len(stacks[from]) - 1
		stacks[to] = append(stacks[to], stacks[from][topindex])
		stacks[from] = stacks[from][:topindex]
	}
}

func parseCommand(scanner *bufio.Scanner) (from, to, count int) {
	cmds := scanner.Text()
	matches := regex.FindStringSubmatch(cmds)
	count, _ = strconv.Atoi(matches[1])
	from, _ = strconv.Atoi(matches[2])
	to, _ = strconv.Atoi(matches[3])
	return from - 1, to - 1, count
}

func parseStacks(scanner *bufio.Scanner) [][]byte {
	stacks := make([][]byte, 0)

	rows := make([]string, 0)
	for scanner.Scan() {
		line := scanner.Text()
		if line[1] == '1' {
			break
		}
		rows = append(rows, line)
	}

	// skip empty line after stacks
	scanner.Scan()

	for c := 1; c < len(rows[0]); c += 4 {
		col := make([]byte, 0)
		for r := len(rows) - 1; r >= 0; r-- {
			if rows[r][c] != ' ' {
				col = append(col, rows[r][c])
			}
		}
		stacks = append(stacks, col)
	}

	return stacks
}
