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
	regE *regexp.Regexp
)

func findParen(e string) (int, int) {
	position := -1
	for i, c := range e {
		if c == '(' {
			position = i
		}
		if c == ')' && position != -1 {
			return position, i
		}
	}
	return -1, -1
}

func solve(e string) string {
	currentE := e
	for {
		// Attempt to find a parenthesis group
		start, end := findParen(currentE)
		if start != -1 {
			res := solve(currentE[start+1 : end])
			strBuild := currentE[:start]
			strBuild += res
			strBuild += currentE[end+1:]
			currentE = strBuild
			continue
		}
		m := regE.FindSubmatchIndex([]byte(currentE))
		if m != nil {
			var res int
			o1, err := strconv.Atoi(currentE[m[2]:m[3]])
			if err != nil {
				log.Fatal(err)
			}
			o2, err := strconv.Atoi(currentE[m[6]:m[7]])
			if err != nil {
				log.Fatal(err)
			}
			switch currentE[m[4]:m[5]] {
			case "+":
				res = o1 + o2
			case "*":
				res = o1 * o2
			}
			strBuild := currentE[:m[0]]
			strBuild += fmt.Sprintf("%d", res)
			strBuild += currentE[m[1]:]
			currentE = strBuild
			continue
		}
		if start == -1 && m == nil {
			return currentE
		}
	}
}

func init() {
	regE = regexp.MustCompile(`([0-9]+) ([+*]) ([0-9]+)`)
}

func main() {
	f, err := os.Open("./input")
	if err != nil {
		log.Fatal(err)
	}

	// Read in the seed
	scanner := bufio.NewScanner(f)
	sum := 0
	for scanner.Scan() {
		line := scanner.Text()
		result := solve(line)
		fmt.Println(result)
		r, err := strconv.Atoi(result)
		if err != nil {
			log.Fatal(err)
		}
		sum += r
	}
	if scanner.Err() != nil {
		log.Fatal(scanner.Err())
	}
	fmt.Println(sum)

}
