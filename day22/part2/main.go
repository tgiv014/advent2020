package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type history struct {
	p1Cards hand
	p2Cards hand
}

type hand []int

func draw(h hand) int {
	c := h[0]
	h = h[1:]
	return c
}

func addToBottom(h hand, c int) {
	h = append(h, c)
}

func score(h hand) int {
	s := 0
	for i, c := range h {
		s += (len(h) - i) * c
	}
	return s
}

func main() {
	f, err := os.Open("./input")
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(f)
	hands := []*hand{}
	var h *hand
	for scanner.Scan() {
		if len(scanner.Text()) == 0 {
			continue
		}
		if strings.HasPrefix(scanner.Text(), "Player") {
			if h != nil {
				hands = append(hands, h)
			}
			// New hand
			h = &hand{}
		} else if h != nil {
			c, err := strconv.Atoi(scanner.Text())
			if err != nil {
				log.Fatal(err)
			}
			*h = append(*h, c)
		}
	}
	if scanner.Err() != nil {
		log.Fatal(scanner.Err())
	}

	hands = append(hands, h)
	winner, p1, p2 := playGame(*hands[0], *hands[1], 1)
	fmt.Println("Player", winner, "wins overall game")
	fmt.Println(p1)
	fmt.Println(p2)
	fmt.Println(score(p1))
	fmt.Println(score(p2))
}

func playGame(p1 hand, p2 hand, gameNum int) (int, hand, hand) {
	hist := []*history{}

	for round := 0; len(p1) > 0 && len(p2) > 0; round++ {

		// Confirm this round has not happened before
		for _, h := range hist {
			matched := true

			if len(p1) != len(h.p1Cards) {
				continue
			}

			if len(p2) != len(h.p2Cards) {
				continue
			}

			for i := range p1 {
				if p1[i] != h.p1Cards[i] {
					matched = false
					continue
				}
			}
			for i := range p2 {
				if p2[i] != h.p2Cards[i] {
					matched = false
					continue
				}
			}
			if !matched {
				continue
			}
			return 1, p1, p2
		}

		newHist := &history{make(hand, len(p1)), make(hand, len(p2))}
		for i, c := range p1 {
			newHist.p1Cards[i] = c
		}
		for i, c := range p2 {
			newHist.p2Cards[i] = c
		}
		hist = append(hist, newHist)

		// Ordinary start
		var roundWinner int
		// Draw cards
		c1 := p1[0]
		p1 = p1[1:]
		c2 := p2[0]
		p2 = p2[1:]
		if len(p1) >= c1 && len(p2) >= c2 {
			// Clone hands
			rp1 := make(hand, c1)
			rp2 := make(hand, c2)
			for i := range rp1 {
				rp1[i] = p1[i]
			}
			for i := range rp2 {
				rp2[i] = p2[i]
			}
			roundWinner, _, _ = playGame(rp1, rp2, gameNum+1)
		} else if c1 > c2 {
			roundWinner = 1
		} else if c2 > c1 {
			roundWinner = 2
		}
		switch roundWinner {
		case 1:
			p1 = append(p1, c1)
			p1 = append(p1, c2)
		case 2:
			p2 = append(p2, c2)
			p2 = append(p2, c1)
		}
	}

	if len(p1) == 0 {
		return 2, p1, p2
	} else if len(p2) == 0 {
		return 1, p1, p2
	} else {
		return 0, p1, p2
	}
}
