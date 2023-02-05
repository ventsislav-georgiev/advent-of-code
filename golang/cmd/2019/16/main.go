package main

import (
	"bytes"
	"fmt"
	"io"

	"github.com/ventsislav-georgiev/advent-of-code/golang/pkg/aoc"
)

func main() {
	aoc.Exec(task1, task2)
}

var basePattern = []int{0, 1, 0, -1}

func task1(in io.Reader) {
	signal, _ := parse(in)
	signal = fft(signal, 100, false)
	printFirst8(signal)
}

func task2(in io.Reader) {
	inputSignal, offset := parse(in)

	size := len(inputSignal) * 10000
	signal := make([]int, size-offset)
	for i := range signal {
		signal[i] = inputSignal[(offset+i)%len(inputSignal)]
	}

	// the offset is in the second half of the signal, so we can optimize the FFT
	// the pattern for the second half of the signal is always 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, ...
	useOptimized := offset > size/2
	signal = fft(signal, 100, useOptimized)
	printFirst8(signal)
}

func fft(signal []int, phases int, useOptimized bool) []int {
	if useOptimized {
		for ; phases > 0; phases-- {
			sum := 0
			for i := len(signal) - 1; i >= 0; i-- {
				sum += signal[i]
				signal[i] = abs(sum) % 10
			}
		}

		return signal
	}

	for ; phases > 0; phases-- {
		phase := make([]int, len(signal))

		for row := range signal {
			sum := 0

			for col, digit := range signal {
				if row > col {
					continue
				}

				if row > len(signal)/2 {
					sum += digit
					continue
				}

				coef := basePattern[((col+1)/(row+1))%len(basePattern)]
				sum += digit * coef
			}

			phase[row] = abs(sum) % 10
		}

		signal = phase
	}

	return signal
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func printFirst8(fft []int) {
	for i := 0; i < 8; i++ {
		fmt.Print(fft[i])
	}
	fmt.Println()
}

func parse(in io.Reader) ([]int, int) {
	input, _ := io.ReadAll(in)
	input = bytes.TrimSpace(input)

	signal := make([]int, len(input))
	for i, b := range input {
		signal[i] = int(b - '0')
	}

	offset := aoc.StrToInt(string(input[:7]))
	return signal, offset
}
