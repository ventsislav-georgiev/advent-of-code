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
	springs := parse(in, 1)
	sumOfCombinations := 0
	for _, spring := range springs {
		sumOfCombinations += countCombinations(spring.Pattern, spring.Groups, map[string]int{})
	}

	fmt.Println(sumOfCombinations)
}

func task2(in io.Reader) {
	springs := parse(in, 5)
	sumOfCombinations := 0
	for _, spring := range springs {
		sumOfCombinations += countCombinations(spring.Pattern, spring.Groups, map[string]int{})
	}

	fmt.Println(sumOfCombinations)
}

type Spring struct {
	Pattern string
	Groups  []int
}

func countCombinations(pattern string, groups []int, visited map[string]int) (count int) {
	pKey := pattern + fmt.Sprintf("%v", groups)
	if count, ok := visited[pKey]; ok {
		return count
	}

	if valid, end := isMatch(pattern, groups); end {
		if valid {
			visited[pKey] = 1
			return 1
		}
		return 0
	}

	ch := pattern[0]

	// Skip position if it's a '.'
	if ch == '.' {
		return countCombinations(pattern[1:], groups, visited)
	}

	if ch == '?' {
		count += countCombinations(pattern[1:], groups, visited)
	}

	// Consume the first group in order if possible
	consumedCount, ok := consumeGroup(pattern, groups[0])
	if !ok {
		visited[pKey] = count
		return count
	}

	// Next group
	count += countCombinations(pattern[consumedCount:], groups[1:], visited)

	visited[pKey] = count
	return count
}

func isMatch(pattern string, groupLengths []int) (valid bool, end bool) {
	numWildcards := strings.Count(pattern, "?")
	numFixed := strings.Count(pattern, "#")

	// Check if there are fixed positions remaining but no more groups
	if numFixed > 0 && len(groupLengths) == 0 {
		return false, true
	}

	pattern = strings.ReplaceAll(pattern, "?", ".")
	end = numWildcards == 0

	// Split the pattern into # groups
	groups := strings.FieldsFunc(pattern, func(r rune) bool {
		return r == '.'
	})

	if len(groups) != len(groupLengths) {
		return false, end
	}

	// Check if the length of each group matches the specified length
	for i, group := range groups {
		if len(group) != groupLengths[i] {
			return false, end
		}
	}

	return true, true
}

func consumeGroup(pattern string, count int) (consumed int, ok bool) {
	if len(pattern) < count {
		return 0, false
	}

	// Fail if the pattern contains a '.' in the specified count
	for i := 0; i < count; i++ {
		if pattern[i] == '.' {
			return 0, false
		}
	}

	// Success if the pattern left is exactly the same length as the count
	if len(pattern) == count {
		return count, true
	}

	nextCh := pattern[count]

	// Consume next character also if it's a '.' or '?'
	if nextCh == '.' || nextCh == '?' {
		return count + 1, true
	}

	return 0, false
}

func parse(in io.Reader, copies int) (springs []Spring) {
	springs = make([]Spring, 0)

	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, " ")
		pattern := parts[0]

		if copies > 1 {
			pattern += "?"
		}

		pattern = strings.Repeat(pattern, copies)

		if copies > 1 {
			pattern = pattern[:len(pattern)-1]
		}

		groups := strings.Split(parts[1], ",")
		spring := Spring{
			Pattern: pattern,
			Groups:  make([]int, 0, len(groups)*copies),
		}

		for c := 0; c < copies; c++ {
			for _, group := range groups {
				spring.Groups = append(spring.Groups, aoc.StrToInt(group))
			}
		}

		springs = append(springs, spring)
	}

	return
}
