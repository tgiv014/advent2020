package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type reference []int

type rule struct {
	n       int
	c       byte
	options []reference
	end     bool
}

func parse(line string) rule {
	substrings := strings.Split(line, ":")
	rulenumS := substrings[0]
	rulenum, err := strconv.Atoi(rulenumS)
	if err != nil {
		log.Fatal(err)
	}

	ruleString := substrings[1][1:]
	// Check if this is an end rule
	if ruleString[0] == '"' && ruleString[2] == '"' {
		return rule{rulenum, ruleString[1], nil, true}
	}

	// Build the list of references to other rules
	// Note: the first layer of the array is an either-or situation
	// Only one of the rules has to pass
	newOptions := []reference{}
	options := strings.Split(ruleString, "|")
	for _, option := range options {

		refs := []int{}
		splitrefs := strings.Split(option, " ")
		for _, refnum := range splitrefs {
			if len(refnum) == 0 {
				continue
			}
			num, err := strconv.Atoi(refnum)
			if err != nil {
				log.Fatal(err)
			}
			refs = append(refs, num)
		}
		newOptions = append(newOptions, refs)
	}
	return rule{rulenum, 0, newOptions, false}
}

func match(s string, rules map[int]*rule, rulenum int, p int) []int {
	// No matches if checking a char out of range
	if p >= len(s) {
		return []int{}
	}

	// If we know we are checking an end node, return p+1 if it matches
	// This ultimately drives the advance in position as we build matches
	if rules[rulenum].end {
		// return []int{p + 1}
		if s[p] == rules[rulenum].c {
			return []int{p + 1}
		} else {
			return []int{}
		}
	}

	matches := []int{}
	for _, option := range rules[rulenum].options {
		positions := []int{p}
		for _, ref := range option {
			newPositions := []int{}
			for _, pos := range positions {
				mcs := match(s, rules, ref, pos)
				for _, mc := range mcs {
					newPositions = append(newPositions, mc)
				}
			}
			positions = newPositions
		}

		for _, pos := range positions {
			matches = append(matches, pos)
		}
	}
	return matches
}

func main() {
	f, err := os.Open("./input")
	if err != nil {
		log.Fatal(err)
	}
	rules := make(map[int]*rule)
	strings := make([]string, 0)

	doneProcessingRules := false

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		if len(line) == 0 && !doneProcessingRules {
			doneProcessingRules = true
			continue
		}

		if !doneProcessingRules {
			r := parse(line)
			rules[r.n] = &r
		} else {
			strings = append(strings, line)
		}
	}
	if scanner.Err() != nil {
		log.Fatal(scanner.Err())
	}

	matches := 0
	for _, s := range strings {
		sMatches := match(s, rules, 0, 0)
		for _, sMatch := range sMatches {
			if len(s) == sMatch {
				matches++
			}
		}

	}
	fmt.Println("Total matches", matches)

	// PART 2! Replace two of the rules
	r1 := parse("8: 42 | 42 8")
	rules[r1.n] = &r1

	r2 := parse("11: 42 31 | 42 11 31")
	rules[r2.n] = &r2

	matches = 0
	for _, s := range strings {
		sMatches := match(s, rules, 0, 0)
		for _, sMatch := range sMatches {
			if len(s) == sMatch {
				matches++
			}
		}

	}
	fmt.Println("Total matches", matches)
}
