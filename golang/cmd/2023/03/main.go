package main

import (
	"bufio"
	"fmt"
	"image"
	"io"

	"github.com/ventsislav-georgiev/advent-of-code/golang/pkg/aoc"
)

func main() {
	aoc.Exec(task1, task2)
}

var directions = []image.Point{
	{-1, -1},
	{0, -1},
	{1, -1},
	{1, 0},
	{1, 1},
	{0, 1},
	{-1, 1},
	{-1, 0},
}

func task1(in io.Reader) {
	symbols, numbers := parse(in)
	engineParts := map[*uint]struct{}{}

	for xy := range symbols {
		for _, dir := range directions {
			xy2 := xy.Add(dir)
			if num, ok := numbers[xy2]; ok {
				engineParts[num] = struct{}{}
			}
		}
	}

	var enginePartsSum uint
	for num := range engineParts {
		enginePartsSum += *num
	}

	fmt.Println(enginePartsSum)
}

func task2(in io.Reader) {
	symbols, numbers := parse(in)
	var gearsRatioSum uint

	for xy, ch := range symbols {
		if ch != '*' {
			continue
		}

		engineParts := map[*uint]struct{}{}

		for _, dir := range directions {
			xy2 := xy.Add(dir)
			if num, ok := numbers[xy2]; ok {
				engineParts[num] = struct{}{}
			}
		}

		if len(engineParts) == 2 {
			gearRatio := uint(1)
			for num := range engineParts {
				gearRatio *= *num
			}

			gearsRatioSum += gearRatio
		}
	}

	fmt.Println(gearsRatioSum)
}

func parse(in io.Reader) (symbols map[image.Point]byte, numbers map[image.Point]*uint) {
	scanner := bufio.NewScanner(in)

	xy := image.Point{-1, -1}
	symbols = map[image.Point]byte{}
	numbers = map[image.Point]*uint{}

	for scanner.Scan() {
		line := scanner.Bytes()
		xy.Y++

		numParts := []byte{}
		num := new(uint)
		for _, ch := range line {
			xy.X++

			if ch == '.' {
				goto add_number
			}

			if ch >= '0' && ch <= '9' {
				numParts = append(numParts, ch)
				numbers[xy] = num
				continue
			}

			symbols[xy] = ch

		add_number:
			if len(numParts) > 0 {
				*num = aoc.Atoui(numParts)
				num = new(uint)
				numParts = numParts[:0]
			}
		}

		// last number on the line
		if len(numParts) > 0 {
			*num = aoc.Atoui(numParts)
			num = new(uint)
		}

		xy.X = -1
	}

	return
}
