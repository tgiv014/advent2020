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
	wpx    int
	wpy    int
}

func (s *ship) rotate_left(input int) {
	input /= 90
	for ; input > 0; input-- {
		s.wpx, s.wpy = -s.wpy, s.wpx
		// s.wpy = s.wpx
	}
}

func (s *ship) rotate_right(input int) {
	input /= 90
	for ; input > 0; input-- {
		s.wpx, s.wpy = s.wpy, -s.wpx
		// s.wpy = -s.wpx
	}
}

func (s *ship) move(input string) {
	command := input[0]
	value, err := strconv.Atoi(input[1:])
	if err != nil {
		log.Fatal(err)
	}
	switch command {
	case 'E':
		s.wpx += value
	case 'N':
		s.wpy += value
	case 'W':
		s.wpx -= value
	case 'S':
		s.wpy -= value
	case 'L':
		// s.facing = (s.facing + value/90) % 4
		s.rotate_left(value)
	case 'R':
		// s.facing = (s.facing - value/90) % 4
		s.rotate_right(value)
	case 'F':
		//CHARGE!
		for i := value; i > 0; i-- {
			s.xpos += s.wpx
			s.ypos += s.wpy
		}
	}
}

func (s *ship) printState() {
	fmt.Printf("Position: (%d,%d)\n", s.xpos, s.ypos)
	fmt.Printf("Manhattan dist: %f\n", math.Abs(float64(s.xpos))+math.Abs(float64(s.ypos)))
	fmt.Printf("WP Pos: (%d,%d)\n", s.wpx, s.wpy)
}

func main() {
	f, err := os.Open("./input")
	if err != nil {
		log.Fatalf("Uh oh %e", err)
	}
	defer f.Close()

	s := ship{}
	s.wpx = 10
	s.wpy = 1

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
