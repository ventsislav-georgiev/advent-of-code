package aoc

func LastIdx[T any](arr []T) int {
	return len(arr) - 1
}

func LastElement[T any](arr []T) T {
	return arr[len(arr)-1]
}

func Reverse[T any](s []T) []T {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}

	return s
}
