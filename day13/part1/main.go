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
	earliest int
	busses   []int
}

func main() {
	f, err := os.Open("./input")
	if err != nil {
		log.Fatalf("Uh oh %e", err)
	}
	defer f.Close()

	in := input{0, make([]int, 0)}

	scanner := bufio.NewScanner(f)
	for i := 0; scanner.Scan(); i++ {
		line := scanner.Text()
		fmt.Println(i, line)
		if i == 0 {
			var err error
			in.earliest, err = strconv.Atoi(line)
			if err != nil {
				log.Fatal(err)
			}
		}
		if i == 1 {
			substrings := strings.Split(line, ",")
			for _, substring := range substrings {
				if substring == "x" {
					continue
				}
				bus, err := strconv.Atoi(substring)
				if err != nil {
					log.Fatal(err)
				}
				in.busses = append(in.busses, bus)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// Data is read in
	fmt.Println(in)

	candidates := make(map[int]int)
	for _, bus := range in.busses {
		time := 0
		for time < in.earliest {
			time += bus
		}
		candidates[bus] = time
	}
	fmt.Println(candidates)

	lowest := 0
	lowestBus := 0
	// Preload minimum
	for lowestBus, lowest = range candidates {
		break
	}
	for bus, time := range candidates {
		if time < lowest {
			lowest = time
			lowestBus = bus
		}
	}

	fmt.Println("First available bus", lowestBus, "departing at", lowest)
	proof := lowestBus * (lowest - in.earliest)
	fmt.Println("Problem 13 proof", proof)
}
