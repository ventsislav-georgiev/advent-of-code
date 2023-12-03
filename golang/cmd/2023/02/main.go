package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/ventsislav-georgiev/advent-of-code/golang/pkg/aoc"
)

func main() {
	aoc.Exec(task1, task2)
}

func task1(in io.Reader) {
	gameIdsSum := 0
	maxAmount := map[string]int{
		"red":   12,
		"green": 13,
		"blue":  14,
	}

	scanner := bufio.NewScanner(in)

	var amount int
	var sets, color string

	gameId := 0
	for scanner.Scan() {
		line := scanner.Text()
		gameId++

		sepIdx := 7
		if gameId > 9 {
			sepIdx++
		}
		if gameId > 99 {
			sepIdx++
		}

		sets = line[sepIdx:]

		for _, set := range strings.Split(sets, ";") {
			for _, cube := range strings.Split(set, ",") {
				fmt.Sscanf(cube, "%d %s", &amount, &color)
				if amount > maxAmount[color] {
					goto next
				}
			}
		}

		gameIdsSum += gameId

	next:
	}

	fmt.Println(gameIdsSum)
}

func task2(in io.Reader) {
	powersSum := 0

	scanner := bufio.NewScanner(in)

	var amount int
	var sets, color string

	gameId := 0
	for scanner.Scan() {
		line := scanner.Text()
		gameId++

		sepIdx := 7
		if gameId > 9 {
			sepIdx++
		}
		if gameId > 99 {
			sepIdx++
		}

		sets = line[sepIdx:]
		minAmount := map[string]int{
			"red":   0,
			"green": 0,
			"blue":  0,
		}

		for _, set := range strings.Split(sets, ";") {
			for _, cube := range strings.Split(set, ",") {
				fmt.Sscanf(cube, "%d %s", &amount, &color)
				if amount > minAmount[color] {
					minAmount[color] = amount
				}
			}
		}

		powersSum += minAmount["red"] * minAmount["green"] * minAmount["blue"]
	}

	fmt.Println(powersSum)
}
