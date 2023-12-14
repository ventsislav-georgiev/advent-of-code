package aoc

import (
	"image"

	"golang.org/x/exp/constraints"
)

func Max[T constraints.Ordered](values ...T) T {
	max := values[0]
	for _, v := range values[1:] {
		if v > max {
			max = v
		}
	}
	return max
}

func Min[T constraints.Ordered](values ...T) T {
	min := values[0]
	for _, v := range values[1:] {
		if v < min {
			min = v
		}
	}
	return min
}

func LeastCommonDenominator(a, b uint) uint {
	return a * b / GreatestCommonDivisor(a, b)
}

func GreatestCommonDivisor(a, b uint) uint {
	if b == 0 {
		return a
	}

	return GreatestCommonDivisor(b, a%b)
}

func ManhattanDistance(a, b image.Point) int {
	return Abs(a.X-b.X) + Abs(a.Y-b.Y)
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}

	return x
}
