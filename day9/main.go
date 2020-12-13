package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

type shiftreg struct {
	buff []int
	cap  int
	n    int
}

func (s *shiftreg) add(n int) {
	s.buff = append(s.buff[1:], n)
	if s.n < s.cap {
		s.n++
	}
}

func (s *shiftreg) full() bool {
	return s.n == s.cap
}

func (s *shiftreg) valid(n int) bool {
	for _, testnum := range s.buff {
		searchnum := n - testnum
		for _, findnum := range s.buff[1:] {
			if findnum == searchnum {
				return true
			}
		}
	}
	return false
}

func main() {
	f, err := os.Open("./input")
	if err != nil {
		log.Fatalf("Uh oh %e", err)
	}
	defer f.Close()
	invalidN := 0
	numbers := make([]int, 0)
	s := shiftreg{make([]int, 25), 25, 0}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		n, err := strconv.Atoi(line)
		if err != nil {
			log.Fatal(err)
		}
		if s.full() {
			valid := s.valid(n)
			if !valid {
				fmt.Println("Found invalid number", n)
				invalidN = n
			}
		}
		numbers = append(numbers, n)
		s.add(n)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	for start_idx, _ := range numbers {
		for end_idx := start_idx + 2; end_idx < len(numbers); end_idx++ {
			sumrange := numbers[start_idx:end_idx]
			sum := 0
			for _, n := range sumrange {
				sum += n
			}
			if sum == invalidN {
				fmt.Println("Found the range", start_idx, end_idx)
				rmin, rmax := sumrange[0], sumrange[0]
				for _, n := range sumrange {
					if n < rmin {
						rmin = n
					}
					if n > rmax {
						rmax = n
					}
				}
				fmt.Println("Weakness:", rmin+rmax)
			}
		}
	}
}
