package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/ventsislav-georgiev/advent-of-code-22/golang/pkg/aoc"
)

func main() {
	aoc.Exec(task1, task2)
}

func task1(in io.Reader) {
	dsizes := parse(in)
	var total uint
	for _, size := range dsizes {
		if size <= 100000 {
			total += size
		}
	}
	fmt.Println(total)
}

func task2(in io.Reader) {
	dsizes := parse(in)
	required := 30000000 - (70000000 - dsizes["/"])
	min := ^uint(0)
	for _, size := range dsizes {
		if size >= required {
			if size < min {
				min = size
			}
		}
	}
	fmt.Println(min)
}

func parse(in io.Reader) map[string]uint {
	scanner := bufio.NewScanner(in)
	pwd := []string{"/"}
	dsizes := map[string]uint{}

	for scanner.Scan() {
		cmd := scanner.Text()

		if cmd[0] == '$' {
			if cmd[2] == 'l' {
				continue
			}
			if cmd[5] == '/' {
				pwd = []string{"/"}
				continue
			}
			if cmd[5] == '.' {
				pwd = pwd[:len(pwd)-1]
				continue
			}

			pwd = append(pwd, cmd[5:])
			continue
		}

		if cmd[0] == 'd' {
			continue
		}

		fsize, _ := strconv.Atoi(strings.Fields(cmd)[0])
		for i := range pwd {
			dsizes[strings.Join(pwd[:i+1], "/")] += uint(fsize)
		}
	}

	return dsizes
}
