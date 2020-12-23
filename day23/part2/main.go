package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

type cup struct {
	id   int
	next *cup
}

func pickDestination(cupMap map[int]*cup, pickedCups []*cup, startId, min, max int) *cup {
	for {
		// Attempt to find the desired cup
		if c, ok := cupMap[startId]; ok {
			safe := true
			for _, pc := range pickedCups {
				if c == pc {
					safe = false
					break
				}
			}
			if safe {
				return c
			}
		}
		// If we got here, it's because we didn't find one
		startId--
		if startId < min {
			startId = max
		}
	}
}

func main() {
	f, err := os.Open("./input")
	if err != nil {
		log.Fatal(err)
	}
	cups := make(map[int]*cup)
	scanner := bufio.NewScanner(f)
	var firstCup *cup = nil
	var prevCup *cup = nil
	nCups := 0
	min := 5 // Sorry
	max := 5
	for scanner.Scan() {
		line := scanner.Text()
		for _, c := range line {
			nCups += 1
			n, err := strconv.Atoi(string(c))
			if err != nil {
				log.Fatal(err)
			}
			if n < min {
				min = n
			}
			if n > max {
				max = n
			}
			cups[n] = &cup{n, nil}
			if prevCup != nil {
				// Link'em
				prevCup.next = cups[n]
			}
			if firstCup == nil {
				firstCup = cups[n]
			}
			prevCup = cups[n]
		}

		// Generate remaining cups
		for n := max + 1; nCups < 1000000; n++ {
			cups[n] = &cup{n, nil}
			// Link'em
			prevCup.next = cups[n]
			prevCup = cups[n]
			nCups++
		}
		max = nCups

		// prevCup now holds the last cup
		prevCup.next = firstCup
	}
	if scanner.Err() != nil {
		log.Fatal(scanner.Err())
	}

	currentCup := firstCup
	for i := 0; i < 10000000; i++ {
		// "Pick up" 3 cups
		pickedCups := []*cup{currentCup.next, currentCup.next.next, currentCup.next.next.next}

		// Remove them from the circle
		currentCup.next = pickedCups[2].next // Current cup's next cup is our 3rd's next
		nCups -= 3

		// Pick a destination cup
		destinationCup := pickDestination(cups, pickedCups, currentCup.id-1, min, max)

		// Insert our slice of cups back into the circle
		pickedCups[2].next = destinationCup.next
		destinationCup.next = pickedCups[0]
		nCups += 3
		currentCup = currentCup.next
	}

	a := cups[1].next
	fmt.Println("Next to cup 1", a.id)
	b := a.next
	fmt.Println("And the next one!", b.id)
	fmt.Println("Product:", a.id*b.id)
}
