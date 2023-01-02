package main

import (
	"bufio"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"

	"github.com/ventsislav-georgiev/advent-of-code/golang/pkg/aoc"
)

func main() {
	aoc.Exec(task1, task2)
}

func task1(in io.Reader) {
	root := parse(in)

	total := 0
	for _, size := range root.Sizes() {
		if size <= 100000 {
			total += size
		}
	}

	fmt.Println(total)
}

func task2(in io.Reader) {
	root := parse(in)

	sizes := root.Sizes()
	required := 30000000 - (70000000 - sizes[0])
	sort.Ints(sizes)

	for _, size := range sizes {
		if size >= required {
			fmt.Println(size)
			break
		}
	}
}

type Dir struct {
	parent    *Dir
	dirs      map[string]*Dir
	filesSize int
	size      int
}

func (d *Dir) Size() int {
	if d.size != 0 {
		return d.size
	}
	d.size = d.filesSize
	for _, child := range d.dirs {
		d.size += child.Size()
	}
	return d.size
}

func (d *Dir) Sizes() []int {
	sizes := []int{}
	sizes = append(sizes, d.Size())
	for _, child := range d.dirs {
		sizes = append(sizes, child.Sizes()...)
	}
	return sizes
}

func parse(in io.Reader) *Dir {
	scanner := bufio.NewScanner(in)
	node := &Dir{
		dirs: map[string]*Dir{},
	}
	root := node

	for scanner.Scan() {
		cmd := scanner.Text()

		if cmd[0] == '$' {
			if cmd[2] == 'l' {
				continue
			}
			if cmd[5] == '/' {
				node = root
				continue
			}
			if cmd[5] == '.' {
				node = node.parent
				continue
			}

			name := cmd[5:]
			child, exists := node.dirs[name]
			if !exists {
				child = &Dir{
					parent: node,
					dirs:   map[string]*Dir{},
				}
				node.dirs[name] = child
			}

			node = child
			continue
		}

		if cmd[0] == 'd' {
			continue
		}

		fsize, _ := strconv.Atoi(strings.Fields(cmd)[0])
		node.filesSize += fsize
	}

	return root
}
