package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func find(numbers []int) int {
	for _, n := range numbers {
		t := 2020 - n
		for _, n2 := range numbers {
			if n2 == t {
				return n * n2
			}
		}
	}
	log.Fatal("Failed to find an answer")
	return -1
}

func main() {
	f, err := os.Open("./input")
	if err != nil {
		log.Fatalf("Uh oh %e", err)
	}
	defer f.Close()

	numbers := make([]int, 0)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		n, err := strconv.Atoi(line)
		if err != nil {
			log.Fatal(err)
		}
		numbers = append(numbers, n)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Println(find(numbers))
}
