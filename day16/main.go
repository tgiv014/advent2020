package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	reField *regexp.Regexp
)

type rule struct {
	lower int
	upper int
}

func (r *rule) validate(input int) bool {
	return (input >= r.lower) && (input <= r.upper)
}

type field struct {
	name  string
	rules []rule
	index int
}

func (f *field) validate(input int) bool {
	// Only one of the rules has to be valid
	for _, r := range f.rules {
		if r.validate(input) {
			return true
		}
	}
	return false
}

type ticket struct {
	values []int
	valid  bool
}

func newTicket(valuestring string) ticket {
	substrings := strings.Split(valuestring, ",")
	if len(substrings) == 0 {
		log.Fatal("Invalid ticket string")
	}
	values := make([]int, 0)
	for _, substring := range substrings {
		v, err := strconv.Atoi(substring)
		if err != nil {
			log.Fatal(err)
		}
		values = append(values, v)
	}
	return ticket{values, true}
}

func init() {
	// Gotta compile all the regex
	var err error
	reField, err = regexp.Compile("^([a-z ]+): ([0-9]+)-([0-9]+) or ([0-9]+)-([0-9]+)$")
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	f, err := os.Open("./input")
	if err != nil {
		log.Fatal(err)
	}
	fields := make([]field, 0)
	tickets := make([]ticket, 0)
	readPhase := 0

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		// line := scanner.Text()
		switch readPhase {
		case 0:
			// Reading in fields
			m := reField.FindSubmatch(scanner.Bytes())
			if m == nil {
				if scanner.Text() == "your ticket:" {
					readPhase = 1
				}
				continue
			}
			r1l, _ := strconv.Atoi(string(m[2]))
			r1u, _ := strconv.Atoi(string(m[3]))
			r2l, _ := strconv.Atoi(string(m[4]))
			r2u, _ := strconv.Atoi(string(m[5]))

			rules := []rule{rule{r1l, r1u}, rule{r2l, r2u}}

			fields = append(fields, field{string(m[1]), rules, -1})

		case 1:
			line := scanner.Text()
			if len(line) == 0 || line == "nearby tickets:" {
				continue
			}
			tickets = append(tickets, newTicket(line))
		}
	}
	if scanner.Err() != nil {
		log.Fatal(scanner.Err())
	}

	errRate := 0

	// Validate all the tickets!
	for i, ticket := range tickets {
		for _, v := range ticket.values {
			// Look for values that are invalid across all fields
			condition := true
			for _, f := range fields {
				// If valid for even one field, stop
				if f.validate(v) {
					condition = false
					break
				}
			}
			if condition {
				errRate += v
				tickets[i].valid = false
			}
		}
	}
	fmt.Println("Error rate:", errRate)

	// I apologize in advance for the horrors that follow this line
	// TODO: Repent
	//=============================================================

	// This will essentially work on process of elimination
	idxfound := make([]bool, len(tickets[0].values))
	finds := 0
	// Keep trying until we've matched them all
	for finds < len(fields) {
		for fn, f := range fields {
			if f.index != -1 {
				// Already found the index for this one
				continue
			}
			nMatches := 0
			lastFind := 0
			// Loop through all value indices
			for i := 0; i < len(tickets[0].values); i++ {
				if idxfound[i] {
					continue
				}
				condition := true
				for _, t := range tickets {
					if !t.valid {
						continue
					}
					if !f.validate(t.values[i]) {
						// Invalid, this column doesn't match up
						condition = false
						break
					}
				}
				if condition {
					nMatches++
					lastFind = i
				}
			}
			if nMatches == 1 {
				fmt.Println("Unique match found for", f.name, lastFind)
				fields[fn].index = lastFind
				finds++
				idxfound[lastFind] = true
			}
		}
	}

	// Get the product of all fields on my ticket that begin with departure
	var product uint64 = 1
	myTicket := tickets[0]
	for _, f := range fields {
		if strings.Split(f.name, " ")[0] == "departure" {
			product *= uint64(myTicket.values[f.index])
		}
	}

	fmt.Println("Product:", product)
}
