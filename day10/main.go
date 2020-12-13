package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
)

type graphNode struct {
	value int
	peers []*graphNode
}

var cache = make(map[int]uint64)

func (g *graphNode) count() uint64 {
	if val, ok := cache[g.value]; ok {
		return val
	}
	if len(g.peers) == 0 {
		return uint64(1)
	}
	t := uint64(0)
	for _, p := range g.peers {
		t += p.count()
	}
	cache[g.value] = t
	return t
}

func main() {
	f, err := os.Open("./input")
	if err != nil {
		log.Fatalf("Uh oh %e", err)
	}
	defer f.Close()

	adapters := make([]int, 0)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		val, err := strconv.Atoi(line)
		if err != nil {
			log.Fatal(err)
		}
		adapters = append(adapters, val)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// Add src and dest adapters and sort
	adapters = append(adapters, 0)
	sort.Ints(adapters)
	adapters = append(adapters, adapters[len(adapters)-1]+3)

	adapterNodes := make([]graphNode, 0)
	deltas := make(map[int]int)
	prev := 0
	for i, val := range adapters {
		adapterNodes = append(adapterNodes, graphNode{val, make([]*graphNode, 0)})
		if i == 0 {
			continue
		}
		delta := val - prev
		prev = val
		deltas[delta]++
	}
	// fmt.Println(deltas)
	fmt.Println("Part 1 result:", deltas[1]*deltas[3])

	// Build peer tables
	for i, _ := range adapterNodes {
		for j := i + 1; j < i+4; j++ {
			if j >= len(adapterNodes) {
				break
			}
			d := adapterNodes[j].value - adapterNodes[i].value
			if d <= 3 {
				// node.peers = append(node.peers, &bond)
				adapterNodes[i].peers = append(adapterNodes[i].peers, &(adapterNodes[j]))
			} else {
				break
			}
		}
	}
	// fmt.Println(adapterNodes)
	combinations := adapterNodes[0].count()
	fmt.Println("Part 2 result: ", combinations)
}
