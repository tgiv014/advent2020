package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type rule struct {
	n    int
	s    string
	refs [][]int
	end  bool
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
		return rule{rulenum, string(ruleString[1]), nil, true}
	}

	// Build the list of references to other rules
	// Note: the first layer of the array is an either-or situation
	// Only one of the rules has to pass
	refs := make([][]int, 0)
	ruleOrs := strings.Split(ruleString, "|")
	for _, ruleOr := range ruleOrs {
		// The other rules that must be true in order in the OR'd rules
		subruleRefs := make([]int, 0)
		subnums := strings.Split(ruleOr, " ")
		for _, subnum := range subnums {
			if len(subnum) == 0 {
				continue
			}
			num, err := strconv.Atoi(subnum)
			if err != nil {
				log.Fatal(err)
			}
			subruleRefs = append(subruleRefs, num)
		}
		refs = append(refs, subruleRefs)
	}
	return rule{rulenum, "", refs, false}
}

func match(s string, rules map[int]*rule, rulenum int, p int) []int {
	// If we know we are checking an end node, return p+1 if it matches
	if rules[rulenum].end && string(s[p]) == rules[rulenum].s {
		return []int{p + 1}
	}
	// No matches if checking a char out of range
	if p >= len(s) {
		return []int{}
	}
	matches := []int{}
	// We will check both sub-rules and append matches from both
	for _, ruleOr := range rules[rulenum].refs {
		// Always start with the given position
		positions := []int{p}

		// Check the required subrules
		for _, subrule := range ruleOr {
			newPositions := []int{}
			for _, pos := range positions {
				// Check if this subrule matches, evaluate any resulting matches
				mcs := match(s, rules, subrule, pos)
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
	doneProcessingRules := false
	scanner := bufio.NewScanner(f)
	strings := make([]string, 0)
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
