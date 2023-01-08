package main

import (
	"bufio"
	"fmt"
	"io"
	"math"

	"github.com/ventsislav-georgiev/advent-of-code/golang/pkg/aoc"
)

func main() {
	aoc.Exec(task1, task2)
}

func task1(in io.Reader) {
	layers := parse(in, 25, 6)

	minZeros := math.MaxInt
	numOnes := 0
	numTwos := 0

	for _, layer := range layers {
		zeros := 0
		ones := 0
		twos := 0

		for _, row := range layer {
			for _, col := range row {
				switch col {
				case '0':
					zeros++
				case '1':
					ones++
				case '2':
					twos++
				}
			}
		}

		if zeros < minZeros {
			minZeros = zeros
			numOnes = ones
			numTwos = twos
		}
	}

	println(numOnes * numTwos)
}

func task2(in io.Reader) {
	width := 25
	height := 6
	layers := parse(in, width, height)

	image := make([][]byte, height)
	for row := 0; row < height; row++ {
		image[row] = make([]byte, width)

		for col := 0; col < width; col++ {
			for _, layer := range layers {
				if layer[row][col] != '2' {
					image[row][col] = layer[row][col]
					break
				}
			}
		}
	}

	for _, row := range image {
		for _, color := range row {
			if color == '1' {
				fmt.Print("█")
			} else {
				fmt.Print(" ")
			}
		}
		println()
	}
}

func parse(in io.Reader, width, height int) [][][]byte {
	scanner := bufio.NewScanner(in)
	scanner.Split(bufio.ScanRunes)

	layers := make([][][]byte, 0)
	for {
		layer := make([][]byte, height)

		for row := 0; row < height; row++ {
			layer[row] = make([]byte, width)

			for col := 0; col < width; col++ {
				if !scanner.Scan() {
					return layers
				}

				layer[row][col] = scanner.Bytes()[0]
			}
		}

		layers = append(layers, layer)
	}
}
