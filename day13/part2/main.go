package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type input struct {
	earliest int64
	busses   []int64
	offsets  []int64
}

func main() {
	f, err := os.Open("./input")
	if err != nil {
		log.Fatalf("Uh oh %e", err)
	}
	defer f.Close()

	in := input{0, make([]int64, 0), make([]int64, 0)}

	scanner := bufio.NewScanner(f)
	for i := 0; scanner.Scan(); i++ {
		line := scanner.Text()
		if i == 0 {
			earliest, err := strconv.Atoi(line)
			if err != nil {
				log.Fatal(err)
			}
			in.earliest = int64(earliest)
		}
		if i == 1 {
			substrings := strings.Split(line, ",")
			for j, substring := range substrings {
				if substring == "x" {
					continue
				}
				bus, err := strconv.Atoi(substring)
				if err != nil {
					log.Fatal(err)
				}
				in.busses = append(in.busses, int64(bus))
				in.offsets = append(in.offsets, int64(j))
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// Data is read in
	fmt.Println(in)

	// Consider:
	// We're looking for the solution t to the system of equations
	// 0 = (t + offsets[i]) % busses[i]
	// for i = [0:len(busses)]
	//
	// This solution works by finding the earliest occurence and
	// period of solutions for 2 busses, then 3 busses, and so on...
	//
	// This works because we know that the solution for n busses
	// has to occur at one of the times that works for n-1 busses as well.
	//
	// This dramatically lowers the size of the search space

	repeat := int64(in.busses[0])
	t := int64(0)
	for busidx := 1; busidx < len(in.busses); busidx++ {
		finds := 0
		for time := t; finds < 2; time += repeat {
			bus := in.busses[busidx]
			if (time+in.offsets[busidx])%bus == 0 {
				// Success!
				if finds == 0 {
					t = time
				} else if finds == 1 {
					repeat = time - t
				}
				finds++
			}
		}
	}
	fmt.Println(t)
}
