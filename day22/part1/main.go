package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type hand struct {
	player int
	cards  []int
}

func (h *hand) draw() int {
	c := h.cards[0]
	h.cards = h.cards[1:]
	return c
}

func (h *hand) addToBottom(c int) {
	h.cards = append(h.cards, c)
}

func (h *hand) score() int {
	s := 0
	for i, c := range h.cards {
		s += (len(h.cards) - i) * c
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
			ps := strings.Split(strings.Trim(scanner.Text(), ":"), " ")[1]
			p, err := strconv.Atoi(ps)
			if err != nil {
				log.Fatal(err)
			}
			h = &hand{p, make([]int, 0)}
		} else if h != nil {
			c, err := strconv.Atoi(scanner.Text())
			if err != nil {
				log.Fatal(err)
			}
			h.cards = append(h.cards, c)
		}
	}
	if scanner.Err() != nil {
		log.Fatal(scanner.Err())
	}

	hands = append(hands, h)
	fmt.Println(hands)

	p1 := hands[0]
	p2 := hands[1]

	for len(p1.cards) > 0 && len(p2.cards) > 0 {
		c1 := p1.draw()
		c2 := p2.draw()
		if c1 > c2 {
			p1.addToBottom(c1)
			p1.addToBottom(c2)
		} else if c2 > c1 {
			p2.addToBottom(c2)
			p2.addToBottom(c1)
		}
	}
	fmt.Println("Player 1's deck:", p1.cards, "Score:", p1.score())
	fmt.Println("Player 2's deck:", p2.cards, "Score:", p2.score())
}
