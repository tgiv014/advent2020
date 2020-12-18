package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type slope struct {
	dx, dy int
}

type universe []string

func main() {
	f, err := os.Open("./input")
	if err != nil {
		log.Fatalf("Uh oh %e", err)
	}
	defer f.Close()

	u := make(universe, 0)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		u = append(u, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	slopes := []slope{
		{1, 1},
		{3, 1},
		{5, 1},
		{7, 1},
		{1, 2},
	}
	product := 1
	for _, s := range slopes {
		nTrees := 0
		x := 0
		for y := 0; y < len(u); y += s.dy {
			if u[y][x] == '#' {
				nTrees += 1
			}
			x = (x + s.dx) % len(u[0])
		}
		product *= nTrees
	}

	fmt.Println("Product:", product)
}
