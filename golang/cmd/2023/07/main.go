package main

import (
	"bufio"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/ventsislav-georgiev/advent-of-code/golang/pkg/aoc"
)

func main() {
	aoc.Exec(task1, task2)
}

func task1(in io.Reader) {
	hands := parse(in, false)

	total := uint(0)
	for rank, hand := range hands {
		total += hand.Bid * uint(rank+1)
	}

	fmt.Println(total)
}

func task2(in io.Reader) {
	hands := parse(in, true)

	total := uint(0)
	for rank, hand := range hands {
		total += hand.Bid * uint(rank+1)
	}

	fmt.Println(total)
}

type Hand struct {
	CardsWeight []uint
	CardsCount  map[byte]uint
	StrengthMap map[uint]uint
	Bid         uint
}

func (h *Hand) GetWeight() uint {
	switch true {
	case h.StrengthMap[5] == 1:
		return 10
	case h.StrengthMap[4] == 1:
		return 9
	case h.StrengthMap[3] == 1 && h.StrengthMap[2] == 1:
		return 8
	case h.StrengthMap[3] == 1:
		return 7
	case h.StrengthMap[2] == 2:
		return 6
	case h.StrengthMap[2] == 1:
		return 5
	default:
		return 0
	}
}

func (h *Hand) UpdateJokers() bool {
	if h.CardsCount['J'] == 0 {
		return false
	}

	h.CardsCount['J']--

	switch true {
	case h.StrengthMap[5] == 1:
		return true
	case h.StrengthMap[4] == 1:
		h.StrengthMap[4] = 0
		h.StrengthMap[5] = 1
		return true
	case h.StrengthMap[3] == 1 && h.StrengthMap[2] == 1:
		return true
	case h.StrengthMap[3] == 1:
		h.StrengthMap[3] = 0
		h.StrengthMap[4] = 1
		return true
	case h.StrengthMap[2] == 2:
		h.StrengthMap[2] = 1
		h.StrengthMap[3] = 1
		return true
	case h.StrengthMap[2] == 1:
		h.StrengthMap[2] = 0
		h.StrengthMap[3] = 1
		return true
	default:
		h.StrengthMap[2] = 1
		return true
	}
}

func (h *Hand) Compare(other *Hand) int {
	curStrength, otherStrength := h.GetWeight(), other.GetWeight()

	if curStrength == otherStrength {
		for i := 0; i < 5; i++ {
			if h.CardsWeight[i] == other.CardsWeight[i] {
				continue
			}

			if h.CardsWeight[i] > other.CardsWeight[i] {
				return 1
			}

			return -1
		}
	}

	if curStrength > otherStrength {
		return 1
	}

	return -1
}

func parse(in io.Reader, hasJoker bool) []Hand {
	scanner := bufio.NewScanner(in)
	hands := []Hand{}

	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), " ")
		handStr := parts[0]
		bidStr := parts[1]

		hand := Hand{
			CardsWeight: []uint{},
			CardsCount:  make(map[byte]uint, 0),
			StrengthMap: make(map[uint]uint, 0),
			Bid:         aoc.Atoui([]byte(bidStr)),
		}

		for i := 0; i < len(handStr); i++ {
			card := handStr[i]
			hand.CardsCount[card]++
			hand.CardsWeight = append(hand.CardsWeight, getWeight(card, hasJoker))
		}

		for card, count := range hand.CardsCount {
			if hasJoker && card == 'J' {
				continue
			}
			hand.StrengthMap[count]++
		}

		if hasJoker {
			for hand.UpdateJokers() {
			}
		}

		hands = append(hands, hand)
	}

	sort.SliceStable(hands, func(i, j int) bool {
		return hands[i].Compare(&hands[j]) == -1
	})

	return hands
}

func getWeight(card byte, hasJoker bool) uint {
	switch card {
	case 'T':
		return 10
	case 'Q':
		return 12
	case 'K':
		return 13
	case 'A':
		return 14
	case 'J':
		if hasJoker {
			return 0
		}

		return 11
	default:
		return uint(card - '0')
	}
}
