package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type mem struct {
	a uint64
	b uint64
}

type game struct {
	m map[uint64]*mem
	i uint64
	p uint64
}

func (g *game) insert(n uint64) {
	if m, ok := g.m[n]; ok {
		// m is a pointer so we can modify its contents
		// Shift in i
		m.a = m.b
		m.b = g.i
	} else {
		// Allocate a new memory object
		g.m[n] = &mem{g.i, g.i}
	}
	// Increment and store most recent number
	g.i++
	g.p = n
}

func (g *game) play() {
	m := g.m[g.p]
	d := m.b - m.a
	g.insert(d)
}

func main() {
	f, err := os.Open("./input")
	if err != nil {
		log.Fatal(err)
	}

	g := game{make(map[uint64]*mem), 0, 0}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		substrings := strings.Split(line, ",")
		for _, substring := range substrings {
			n, err := strconv.ParseUint(substring, 10, 64)
			if err != nil {
				log.Fatal(err)
			}
			g.insert(n)
		}
	}
	if scanner.Err() != nil {
		log.Fatal(scanner.Err())
	}

	for g.i < 2020 {
		g.play()
	}
	fmt.Println("Value after 2020 rounds:", g.p)

	for g.i < 30000000 {
		g.play()
	}
	fmt.Println("Value after 30000000 rounds:", g.p)
}
