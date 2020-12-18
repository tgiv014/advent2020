package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type universe map[pos]*cell

type pos struct {
	x, y, z int
}

type cell struct {
	u      universe
	p      pos
	active bool
	next   bool
}

func (c *cell) getNeighbor(dx, dy, dz int, canAlloc bool) *cell {
	p := c.p
	np := pos{p.x + dx, p.y + dy, p.z + dz}

	// Find and return neighbor
	if n, ok := c.u[np]; ok {
		return n
	}
	// No neighbor found... better make one
	newNeighbor := cell{c.u, np, false, false}
	if canAlloc {
		c.u[np] = &newNeighbor
		newNeighbor.run(false)
	}
	return &newNeighbor
}

func (c *cell) apply() {
	c.active = c.next
}

func (c *cell) run(canAlloc bool) {
	activeNeighbors := 0
	// Loop through all neighbors
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			for dz := -1; dz <= 1; dz++ {
				// Always skip self
				if dx == 0 && dy == 0 && dz == 0 {
					continue
				}
				n := c.getNeighbor(dx, dy, dz, canAlloc)
				if n.active {
					activeNeighbors++
				}
			}
		}
	}
	if c.active {
		c.next = (activeNeighbors == 2 || activeNeighbors == 3)
	} else {
		c.next = (activeNeighbors == 3)
	}
}

func main() {
	f, err := os.Open("./input")
	if err != nil {
		log.Fatal(err)
	}
	// Just for grins let's do this with a map
	u := make(universe)

	// Read in the seed
	scanner := bufio.NewScanner(f)
	for y := 0; scanner.Scan(); y++ {
		// line := scanner.Text()
		line := scanner.Text()
		for x, c := range line {
			p := pos{x, y, 0}
			newC := cell{u, p, c == '#', false}
			u[p] = &newC
		}
	}
	if scanner.Err() != nil {
		log.Fatal(scanner.Err())
	}

	for i := 0; i < 6; i++ {
		for _, c := range u {
			// Top level runs are allowed to allocate new cells
			c.run(true)
		}
		for _, c := range u {
			c.apply()
		}
	}

	activeCells := 0
	for _, c := range u {
		if c.active {
			activeCells++
		}
	}
	fmt.Println(activeCells)

}
