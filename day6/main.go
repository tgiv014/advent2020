package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type Group struct {
	answers  [26]int
	n_people int
}

func (g *Group) parse(s string) {
	g.n_people += 1
	for _, c := range s {
		c = c - 'a'
		g.answers[c] += 1
	}
}

func (g *Group) count() int {
	total := 0
	for _, n := range g.answers {
		if n == g.n_people {
			total += 1
		}
	}
	return total
}

func main() {
	f, err := os.Open("./input")
	if err != nil {
		log.Fatalf("Uh oh %e", err)
	}
	defer f.Close()

	current_group := Group{}
	total := 0
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		// If the line is empty, we've encountered a group boundary
		if len(line) == 0 {
			total += current_group.count()
			current_group = Group{}
			continue
		}
		current_group.parse(line)
	}
	total += current_group.count()
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Total: ", total)
}
