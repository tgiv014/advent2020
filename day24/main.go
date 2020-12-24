package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type coord struct {
	x, y, z int
}

type tile struct {
	pos   coord
	state bool // False = white
	next  bool
}

func (t *tile) ensureAllocated(m map[coord]*tile) {
	directions := []coord{coord{1, -1, 0}, coord{1, 0, -1}, coord{0, 1, -1},
		coord{-1, 1, 0}, coord{-1, 0, 1}, coord{0, -1, 1}}
	for _, d := range directions {
		newC := coord{t.pos.x + d.x, t.pos.y + d.y, t.pos.z + d.z}
		t := m[newC]
		if t == nil {
			m[newC] = &tile{newC, false, false}
		}
	}
}

func (t *tile) getNeighbors(m map[coord]*tile) int {
	n := 0
	directions := []coord{coord{1, -1, 0}, coord{1, 0, -1}, coord{0, 1, -1},
		coord{-1, 1, 0}, coord{-1, 0, 1}, coord{0, -1, 1}}
	for _, d := range directions {
		newC := coord{t.pos.x + d.x, t.pos.y + d.y, t.pos.z + d.z}
		t := m[newC]
		if t == nil {
			m[newC] = &tile{newC, false, false}
		}
		if m[newC].state {
			n++
		}
	}
	return n
}

func (t *tile) run(m map[coord]*tile) {
	nNeighbors := t.getNeighbors(m)
	t.next = t.state
	if t.state { // If black
		if nNeighbors == 0 || nNeighbors > 2 {
			t.next = false
		}
	} else { // White
		if nNeighbors == 2 {
			t.next = true
		}
	}
}

func (t *tile) apply() {
	t.state = t.next
}

func countBlackTiles(m map[coord]*tile) int {
	nBlack := 0
	for _, tile := range m {
		// fmt.Println(c, tile)
		if tile.state {
			nBlack++
		}
	}
	return nBlack
}

func coordFromLine(s string) coord {
	c := coord{}

	for i := 0; i < len(s); {
		// Be greedy and attempt to get the longer tokens first
		if i+2 <= len(s) {
			switch s[i : i+2] {
			case "ne":
				c.x += 1
				c.z -= 1
				i += 2
				continue
			case "se":
				i += 2
				c.y -= 1
				c.z += 1
				continue
			case "sw":
				i += 2
				c.x -= 1
				c.z += 1
				continue
			case "nw":
				i += 2
				c.y += 1
				c.z -= 1
				continue
			}
		}
		switch s[i : i+1] {
		case "e":
			i += 1
			c.x += 1
			c.y -= 1
		case "w":
			i += 1
			c.x -= 1
			c.y += 1
		default:
			fmt.Println("Heck!")
		}
	}
	return c
}

func main() {
	f, err := os.Open("./input")
	if err != nil {
		log.Fatal(err)
	}

	tiles := make(map[coord]*tile)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		c := coordFromLine(scanner.Text())
		if tiles[c] == nil {
			tiles[c] = &tile{c, false, false}
		}
		tiles[c].state = !tiles[c].state
	}
	if scanner.Err() != nil {
		log.Fatal(scanner.Err())
	}

	fmt.Println("Initial number of black tiles:", countBlackTiles(tiles))

	for i := 0; i < 100; i++ {
		// Make sure that tiles around other black tiles have been allocated
		// Only doing this around black tiles ensures that we don't overgrow
		for _, tile := range tiles {
			if tile.state {
				tile.ensureAllocated(tiles)
			}
		}
		for _, tile := range tiles {
			tile.run(tiles)
		}
		for _, tile := range tiles {
			tile.apply()
		}
	}
	fmt.Println("# Black tiles after 100 days:", countBlackTiles(tiles))
}
