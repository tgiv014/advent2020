package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

var (
	rePW *regexp.Regexp
)

func count(c rune, s string) int {
	t := 0
	for _, sc := range s {
		if sc == c {
			t++
		}
	}
	return t
}

func init() {
	rePW = regexp.MustCompile(`([0-9]+)-([0-9]+) ([a-z]): ([a-z]+)`)
}

func main() {
	f, err := os.Open("./input")
	if err != nil {
		log.Fatalf("Uh oh %e", err)
	}
	defer f.Close()

	nValid := 0

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		m := rePW.FindSubmatch(scanner.Bytes())
		min, err := strconv.Atoi(string(m[1]))
		if err != nil {
			log.Fatal(err)
		}
		max, err := strconv.Atoi(string(m[2]))
		if err != nil {
			log.Fatal(err)
		}
		c := rune(m[3][0])
		pw := string(m[4])

		n := count(c, pw)
		if n >= min && n <= max {
			nValid++
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println(nValid)
}
