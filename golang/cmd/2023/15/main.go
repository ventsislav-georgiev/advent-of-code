package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/ventsislav-georgiev/advent-of-code/golang/pkg/aoc"
)

func main() {
	aoc.Exec(task1, task2)
}

func task1(in io.Reader) {
	var initSum int

	parse(in, func(seq string) {
		initSum += hashSeq([]byte(seq))
	})

	fmt.Println(initSum)
}

func task2(in io.Reader) {
	boxes := Boxes{}
	parse(in, boxes.Update)
	fmt.Println(boxes.CalcFocusingPower())
}

func hashSeq(seq []byte) int {
	var hash int
	salt := 17
	upperBound := 256

	for _, ch := range seq {
		hash += int(ch)
		hash *= salt
		hash %= upperBound
	}

	return hash
}

func parse(in io.Reader, action func(seq string)) {
	scanner := bufio.NewScanner(in)
	scanner.Split(bufio.ScanBytes)

	seq := make([]byte, 0, 20)
	for scanner.Scan() {
		ch := scanner.Bytes()[0]
		if ch == '\n' {
			continue
		}

		if ch == ',' {
			action(string(seq))
			seq = seq[:0]
			continue
		}

		seq = append(seq, ch)
	}

	action(string(seq))
}

func parseStep(step string) (label string, valStr string, remove bool) {
	remove = true

	sepIdx := strings.Index(step, "-")
	if sepIdx != -1 {
		label = step[:sepIdx]
		return
	}

	stepParts := strings.Split(step, "=")
	label = stepParts[0]
	valStr = stepParts[1]
	remove = false

	return
}

type Boxes [256]*Box

func (b *Boxes) Update(step string) {
	if step == "" {
		return
	}

	label, valStr, remove := parseStep(step)
	boxIndex := hashSeq([]byte(label))

	if b[boxIndex] == nil {
		b[boxIndex] = &Box{
			LensesMap: map[string]*Lense{},
		}
	}

	b[boxIndex].Update(label, valStr, remove)
}

func (b *Boxes) CalcFocusingPower() int {
	var focusingPower int
	for i, box := range b {
		if box == nil {
			continue
		}

		focusingPower += box.CalcFocusingPower(i)
	}

	return focusingPower
}

type Box struct {
	FirstLense *Lense
	LastLense  *Lense
	LensesMap  map[string]*Lense
}

func (b *Box) Update(label string, valStr string, remove bool) {
	if remove {
		b.RemoveLense(label)
		return
	}

	val, _ := strconv.Atoi(valStr)
	b.AddLense(label, val)
}

func (b *Box) CalcFocusingPower(boxIdx int) int {
	var focusingPower int
	var i int
	for lense := b.FirstLense; lense != nil; lense = lense.NextLense {
		focusingPower += (boxIdx + 1) * (i + 1) * lense.FocalLength
		i++
	}

	return focusingPower
}

func (b *Box) AddLense(label string, val int) {
	if l, found := b.LensesMap[label]; found {
		l.FocalLength = val
		return
	}

	lense := &Lense{
		Label:       label,
		FocalLength: val,
	}

	b.LensesMap[label] = lense

	if b.FirstLense == nil {
		b.FirstLense = lense
		b.LastLense = lense
	} else {
		b.LastLense.NextLense = lense
		lense.PrevLense = b.LastLense
		b.LastLense = lense
	}
}

func (b *Box) RemoveLense(label string) {
	lense, found := b.LensesMap[label]
	if !found {
		return
	}

	lense.RemoveFromChain()
	delete(b.LensesMap, label)

	if b.FirstLense == lense {
		b.FirstLense = lense.NextLense
	}

	if b.LastLense == lense {
		b.LastLense = lense.PrevLense
	}
}

type Lense struct {
	Label       string
	FocalLength int
	NextLense   *Lense
	PrevLense   *Lense
}

func (l *Lense) RemoveFromChain() {
	if l.PrevLense != nil {
		l.PrevLense.NextLense = l.NextLense
	}

	if l.NextLense != nil {
		l.NextLense.PrevLense = l.PrevLense
	}
}
