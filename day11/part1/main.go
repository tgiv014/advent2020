package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

// An array of directions to look for seats in
type vec struct {
	x int
	y int
}

var vecs = []vec{
	{1, 0},
	{0, 1},
	{-1, 0},
	{0, -1},
	{1, 1},
	{-1, 1},
	{-1, -1},
	{1, -1},
}

// Container for simulation state
type sim struct {
	u      [][]byte // Current state
	f      [][]byte // Future state
	deltas int      // # Changes in this cycle
}

// Get adjacent chairs in all 8 directions
func (s *sim) getNeighbors(x int, y int) []byte {
	neighbors := make([]byte, 0)
	for _, dir := range vecs {
		nx := x + dir.x
		ny := y + dir.y
		if nx < 0 || ny < 0 ||
			nx >= len(s.u[0]) || ny >= len(s.u) {
			continue
		}
		neighbors = append(neighbors, s.u[ny][nx])
	}
	return neighbors
}

// Get the number of neighbors matching b
func (s *sim) countNeighbors(x int, y int, b byte) int {
	total := 0
	neighbors := s.getNeighbors(x, y)
	for _, neighbor := range neighbors {
		if neighbor == b {
			total++
		}
	}
	return total
}

// Count total number of taken seats in the sim
func (s *sim) nTaken() int {
	total := 0
	for _, r := range s.u {
		for _, c := range r {
			if c == '#' {
				total++
			}
		}
	}
	return total
}

func (s *sim) iterate() {
	s.deltas = 0
	for y, row := range s.u {
		for x, cell := range row {
			switch cell {
			case '.':
				s.f[y][x] = '.'
			case 'L':
				if s.countNeighbors(x, y, '#') == 0 {
					s.f[y][x] = '#'
					s.deltas++
				}
			case '#':
				if s.countNeighbors(x, y, '#') >= 4 {
					s.f[y][x] = 'L'
					s.deltas++
				}
			}
		}
	}
	for i := range s.f {
		copy(s.u[i], s.f[i])
	}
}

func main() {
	f, err := os.Open("./input")
	if err != nil {
		log.Fatalf("Uh oh %e", err)
	}
	defer f.Close()

	s := sim{}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.Trim(line, "\n\r ")
		s.u = append(s.u, []byte(line))
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// Allocate future array
	// Same size as current state
	s.f = make([][]byte, len(s.u))
	for i := range s.f {
		s.f[i] = make([]byte, len(s.u[i]))
	}

	for {
		s.iterate()
		if s.deltas == 0 {
			break
		}
	}

	fmt.Println("Taken seats", s.nTaken())
}
