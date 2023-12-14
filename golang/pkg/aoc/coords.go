package aoc

func ToKeyXY(xy [2]int) uint64 {
	return ToKey(xy[0], xy[1])
}

func ToKey(x, y int) uint64 {
	return uint64(x)<<16 + uint64(y)
}
