package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
)

var dirToCD = []byte{'E', 'N', 'W', 'S'}

type ship struct {
	xpos   int
	ypos   int // +x -> east, -x -> west, +y -> north, -y -> south
	facing int // 0 through 3 with 0 being east, proceeding ccw (like a radian circle)
}

func (s *ship) rotate(input int) {
	s.facing += input / 90
	if s.facing < 0 {
		s.facing += 4
	}
	s.facing %= 4
}

func (s *ship) move(input string) {
	command := input[0]
	value, err := strconv.Atoi(input[1:])
	if err != nil {
		log.Fatal(err)
	}
	switch command {
	case 'F':
		// Translate to NESW instr
		idx := s.facing % 4
		command = dirToCD[idx]
	case 'B':
		// Translate to NESW instr
		idx := (s.facing + 2) % 4
		command = dirToCD[idx]
	}

	switch command {
	case 'E':
		s.xpos += value
	case 'N':
		s.ypos += value
	case 'W':
		s.xpos -= value
	case 'S':
		s.ypos -= value
	case 'L':
		// s.facing = (s.facing + value/90) % 4
		s.rotate(value)
	case 'R':
		// s.facing = (s.facing - value/90) % 4
		s.rotate(-value)
	}
}

func (s *ship) printState() {
	fmt.Printf("Position: (%d,%d)\n", s.xpos, s.ypos)
	fmt.Printf("Manhattan dist: %f\n", math.Abs(float64(s.xpos))+math.Abs(float64(s.ypos)))
}

func main() {
	f, err := os.Open("./input")
	if err != nil {
		log.Fatalf("Uh oh %e", err)
	}
	defer f.Close()

	s := ship{}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		s.move(line)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	s.printState()

}
