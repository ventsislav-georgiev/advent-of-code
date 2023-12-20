package aoc

func CopyMap[T any](m map[string]T) map[string]T {
	c := make(map[string]T, len(m))
	for k, v := range m {
		c[k] = v
	}

	return c
}
