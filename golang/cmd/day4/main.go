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

func task1(in io.ReadCloser) {
	scanner := bufio.NewScanner(in)

	result := 0
	for scanner.Scan() {
		pairs := strings.Split(scanner.Text(), ",")
		a := strings.Split(pairs[0], "-")
		b := strings.Split(pairs[1], "-")
		a_x, _ := strconv.Atoi(a[0])
		a_y, _ := strconv.Atoi(a[1])
		b_x, _ := strconv.Atoi(b[0])
		b_y, _ := strconv.Atoi(b[1])
		if (a_x <= b_x && a_y >= b_y) || (b_x <= a_x && b_y >= a_y) {
			result++
		}
	}

	fmt.Println(result)
}

func task2(in io.ReadCloser) {
	scanner := bufio.NewScanner(in)

	result := 0
	for scanner.Scan() {
		pairs := strings.Split(scanner.Text(), ",")
		a := strings.Split(pairs[0], "-")
		b := strings.Split(pairs[1], "-")
		a_x, _ := strconv.Atoi(a[0])
		a_y, _ := strconv.Atoi(a[1])
		b_x, _ := strconv.Atoi(b[0])
		b_y, _ := strconv.Atoi(b[1])
		if (a_x <= b_x && a_y >= b_x) || (b_x <= a_x && b_y >= a_x) {
			result++
		}
	}

	fmt.Println(result)
}
