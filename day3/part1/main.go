package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

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

	nTrees := 0
	dx := 3
	dy := 1
	x := 0
	for y := 0; y < len(u); y += dy {
		if u[y][x] == '#' {
			nTrees += 1
		}
		x = (x + dx) % len(u[0])
	}

	fmt.Println("Trees:", nTrees)
}
