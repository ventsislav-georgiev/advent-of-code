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
	pipelines, parts := parse(in)
	rules := getAcceptRules(pipelines["in"], map[string]Range{})
	var sumOfProps int

	for _, part := range parts {
		for _, rule := range rules {
			accept := true

			for prop, val := range part.Values {
				if propRange, found := rule.Ranges[prop]; found {
					if val < propRange.Min || val > propRange.Max {
						accept = false
						break
					}
				}
			}

			if accept {
				for _, val := range part.Values {
					sumOfProps += val
				}
				break
			}
		}
	}

	println(sumOfProps)
}

func task2(in io.Reader) {
	pipelines, _ := parse(in)
	rules := getAcceptRules(pipelines["in"], map[string]Range{
		"x": {Min: 1, Max: 4000},
		"m": {Min: 1, Max: 4000},
		"a": {Min: 1, Max: 4000},
		"s": {Min: 1, Max: 4000},
	})

	var count uint
	for _, rule := range rules {
		product := uint(1)
		for _, propRange := range rule.Ranges {
			product *= uint(propRange.Max - propRange.Min + 1)
		}

		count += product
	}

	fmt.Println(count)
}

type Pipeline struct {
	Steps []*Step
}

type Step struct {
	Prop   string
	Op     byte
	Val    int
	NextID string
	Next   *Pipeline
}

type Part struct {
	Values map[string]int
}

type Rule struct {
	Ranges map[string]Range
}

type Range struct {
	Min, Max int
}

func getAcceptRules(pipeline *Pipeline, curRanges map[string]Range) (rules []Rule) {
	rules = make([]Rule, 0)

	for _, step := range pipeline.Steps {
		stepIntoRanges := aoc.CopyMap(curRanges)
		propRange := curRanges[step.Prop]
		if propRange.Min == 0 && propRange.Max == 0 {
			propRange.Min = 1
			propRange.Max = 4000
		}

		switch step.Op {
		case '>':
			stepIntoRanges[step.Prop] = Range{
				Min: aoc.Max(propRange.Min, step.Val+1),
				Max: propRange.Max,
			}
			curRanges[step.Prop] = Range{
				Min: propRange.Min,
				Max: aoc.Min(propRange.Max, step.Val),
			}
		case '<':
			stepIntoRanges[step.Prop] = Range{
				Min: propRange.Min,
				Max: aoc.Min(propRange.Max, step.Val-1),
			}
			curRanges[step.Prop] = Range{
				Min: aoc.Max(propRange.Min, step.Val),
				Max: propRange.Max,
			}
		}

		if step.NextID == "R" {
			continue
		}

		if step.Next != nil {
			rules = append(rules, getAcceptRules(step.Next, stepIntoRanges)...)
			continue
		}

		rules = append(rules, Rule{
			Ranges: stepIntoRanges,
		})
	}

	return
}

func parse(in io.Reader) (pipelines map[string]*Pipeline, parts []Part) {
	pipelines = make(map[string]*Pipeline)

	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}

		fields := strings.FieldsFunc(line, func(r rune) bool {
			return r == '{' || r == '}'
		})

		pipelineID := fields[0]
		pipeline := Pipeline{
			Steps: make([]*Step, 0),
		}
		pipelines[pipelineID] = &pipeline

		steps := strings.Split(fields[1], ",")
		for _, stepData := range steps {
			step := Step{}
			nextSepIdx := strings.Index(stepData, ":")

			if nextSepIdx == -1 {
				step.NextID = stepData
			} else {
				step.NextID = stepData[nextSepIdx+1:]
				step.Prop = stepData[0:1]
				step.Op = stepData[1:2][0]
				step.Val, _ = strconv.Atoi(stepData[2:nextSepIdx])
			}

			pipeline.Steps = append(pipeline.Steps, &step)
		}
	}

	for _, pipeline := range pipelines {
		for _, step := range pipeline.Steps {
			if step.NextID != "R" && step.NextID != "A" {
				step.Next = pipelines[step.NextID]
			}
		}
	}

	parts = []Part{}

	for scanner.Scan() {
		line := scanner.Text()
		part := Part{
			Values: make(map[string]int),
		}

		line = line[1 : len(line)-1] // remove the first and last curly braces
		fields := strings.Split(line, ",")
		for _, field := range fields {
			kv := strings.Split(field, "=")
			key := kv[0]
			val, _ := strconv.Atoi(kv[1])
			part.Values[key] = val
		}

		parts = append(parts, part)
	}

	return
}
