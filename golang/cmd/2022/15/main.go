package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"sort"

	"github.com/ventsislav-georgiev/advent-of-code/golang/pkg/aoc"
)

func main() {
	aoc.Exec(task1, task2)
}

func task1(in io.Reader) {
	scanner := bufio.NewScanner(in)
	row := int32(2000000)
	rowSignal := map[uint64]struct{}{}
	rowBeacons := []int32{}

	toKey := func(x, y int32) uint64 {
		return uint64(x)<<32 + uint64(y)
	}

	for scanner.Scan() {
		var sx, sy, bx, by int32
		fmt.Sscanf(scanner.Text(), "Sensor at x=%d, y=%d: closest beacon is at x=%d, y=%d", &sx, &sy, &bx, &by)

		bdist := int32(math.Abs(float64(sx-bx)) + math.Abs(float64(sy-by)))
		rowdist := int32(math.Abs(float64(sy - row)))
		if rowdist > bdist {
			continue
		}

		rminx := sx - bdist + rowdist
		rmaxx := rminx + 2*(bdist-rowdist) + 1
		for x := rminx; x < rmaxx; x++ {
			rowSignal[toKey(x, row)] = struct{}{}
		}

		if by == row {
			rowBeacons = append(rowBeacons, bx)
		}
	}

	for _, bx := range rowBeacons {
		delete(rowSignal, toKey(bx, row))
	}

	fmt.Println(len(rowSignal))
}

func task2(in io.Reader) {
	scanner := bufio.NewScanner(in)
	ranges := [4000000][][2]int{}

	updateRange := func(y, rminx, rmaxx int) bool {
		updated := false
		for i, r := range ranges[y] {
			if rminx >= r[0] && rminx <= r[1] && rmaxx >= r[1] {
				updated = false
				ranges[y][i][1] = rmaxx
			} else if rmaxx >= r[0] && rmaxx <= r[1] && rminx <= r[0] {
				updated = false
				ranges[y][i][0] = rminx
			} else if rminx <= r[0] && rmaxx >= r[1] {
				updated = false
				ranges[y][i][0] = rminx
				ranges[y][i][1] = rmaxx
			}
		}
		return updated
	}

	addRange := func(y, sx, sy, bdist int) {
		if y < 0 || y >= 4000000 {
			return
		}

		ydist := int(math.Abs(float64(sy - y)))
		rminx := sx - bdist + ydist
		rmaxx := rminx + 2*(bdist-ydist) + 1

		if len(ranges[y]) == 0 || !updateRange(y, rminx, rmaxx) {
			ranges[y] = append(ranges[y], [2]int{rminx, rmaxx})
		}
	}

	for scanner.Scan() {
		var sx, sy, bx, by int
		fmt.Sscanf(scanner.Text(), "Sensor at x=%d, y=%d: closest beacon is at x=%d, y=%d", &sx, &sy, &bx, &by)
		bdist := int(math.Abs(float64(sx-bx)) + math.Abs(float64(sy-by)))
		for y := sy; y <= sy+bdist; y++ {
			addRange(y, sx, sy, bdist)
		}
		for y := sy; y >= sy-bdist; y-- {
			addRange(y, sx, sy, bdist)
		}
	}

	mergeOverlapping := func(y int) {
		row := ranges[y]
		if len(row) == 0 {
			return
		}

		sort.Slice(row, func(i, j int) bool {
			return row[i][0] < row[j][0]
		})

		stack := [][2]int{row[0]}
		for _, i := range row[1:] {
			last := stack[len(stack)-1]
			if last[1] < i[0] {
				stack = append(stack, i)
			} else if last[1] < i[1] {
				stack[len(stack)-1][1] = i[1]
			}
		}

		ranges[y] = stack
	}

	for y := 0; y < 4000000; y++ {
		mergeOverlapping(y)
		if len(ranges[y]) > 1 && ranges[y][0][0] > 0 || ranges[y][0][1] < 4000000 {
			x := ranges[y][0][1]
			fmt.Println(x*4000000 + y)
			return
		}
	}
}
