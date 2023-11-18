package main

import (
	"strings"
	"testing"
)

var inputsAndAnswers = map[string]int{
	`########################
#f.D.E.e.C.b.A.@.a.B.c.#
######################.#
#d.....................#
########################`: 86,
	`########################
#...............b.C.D.f#
#.######################
#.....@.a.B.c.d.A.e.F.g#
########################`: 132,
	`#################
#i.G..c...e..H.p#
########.########
#j.A..b...f..D.o#
########@########
#k.E..a...g..B.n#
########.########
#l.F..d...h..C.m#
#################`: 136,
	`########################
#@..............ac.GI.b#
###d#e#f################
###A#B#C################
###g#h#i################
########################`: 81,
}

func TestTask1(t *testing.T) {
	for in, ans := range inputsAndAnswers {
		grid := parse(strings.NewReader(in))
		response := leastStepsToAllKeys(grid)
		if response != ans {
			t.Errorf("%d != %d", response, ans)
		}
	}
}
