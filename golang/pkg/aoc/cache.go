package aoc

import "fmt"

func Hash(str string) string {
	var hash [32]rune

	for i, ch := range str {
		hash[i%32] ^= ch
	}

	return fmt.Sprintf("%x", string(hash[:]))
}
