package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type Plane struct {
	seats [128][8]bool
}

func (pl *Plane) take_seat(pt Point) {
	pl.seats[pt.r][pt.c] = true
}

func (pl *Plane) seat_taken(pt Point) bool {
	return pl.seats[pt.r][pt.c]
}

type Point struct {
	r, c int
}

func (pt Point) get_sid() int {
	return pt.r*8 + pt.c
}

func parse_bsp(bsp string) Point {
	lp_r := 0
	hp_r := 127
	lp_c := 0
	hp_c := 7

	for _, c := range bsp {
		height := hp_r - lp_r + 1
		width := hp_c - lp_c + 1
		switch c {
		case 'F':
			hp_r -= height / 2
		case 'B':
			lp_r += height / 2
		case 'L':
			hp_c -= width / 2
		case 'R':
			lp_c += width / 2
		}
	}
	return Point{lp_r, lp_c}
}

func pt_from_sid(sid int) Point {
	r := sid / 8
	c := sid % 8
	return Point{r, c}
}

func main() {
	f, err := os.Open("./input")
	if err != nil {
		log.Fatalf("Uh oh %e", err)
	}
	defer f.Close()

	highest := 0
	plane := Plane{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		pt := parse_bsp(scanner.Text())
		sid := pt.get_sid()
		if sid > highest {
			highest = sid
		}
		plane.take_seat(pt)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	for y := 0; y < 128; y++ {
		for x := 0; x < 8; x++ {
			if y == 0 || y == 127 {
				continue
			}
			my_pt := Point{y, x}
			my_sid := my_pt.get_sid()
			pt_left := pt_from_sid(my_sid - 1)
			pt_right := pt_from_sid(my_sid + 1)

			if !plane.seat_taken(my_pt) &&
				plane.seat_taken(pt_left) &&
				plane.seat_taken(pt_right) {
				fmt.Printf("My ticker: Row %d Col %d sid %d\n", y, x, my_sid)
			}
		}
	}
	fmt.Println("Highest sid encountered: ", highest)
}
