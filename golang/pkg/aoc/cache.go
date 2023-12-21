package aoc

import (
	"crypto/md5"
	"fmt"
)

func Hash(str string) string {
	var hash [32]rune

	for i, ch := range str {
		hash[i%32] ^= ch
	}

	return fmt.Sprintf("%x", string(hash[:]))
}

func HashKeys[T any](m map[string]T) string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	return HashStrings(keys...)
}

func HashElements(list []string) string {
	return HashStrings(list...)
}

func HashStrings(strs ...string) string {
	return HashBytes([]byte(JoinStrings(strs...)))
}

func HashBytes(bytes []byte) string {
	return string(HashBytesRaw(bytes...))
}

func HashBytesRaw(bytes ...byte) []byte {
	return []byte(fmt.Sprintf("%x", md5.Sum(bytes)))
}
