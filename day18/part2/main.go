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
	regA *regexp.Regexp
	regM *regexp.Regexp
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
		mA := regA.FindSubmatchIndex([]byte(currentE))
		if mA != nil {
			o1, err := strconv.Atoi(currentE[mA[2]:mA[3]])
			if err != nil {
				log.Fatal(err)
			}
			o2, err := strconv.Atoi(currentE[mA[4]:mA[5]])
			if err != nil {
				log.Fatal(err)
			}
			strBuild := currentE[:mA[0]]
			strBuild += fmt.Sprintf("%d", o1+o2)
			strBuild += currentE[mA[1]:]
			currentE = strBuild
			continue
		}
		mM := regM.FindSubmatchIndex([]byte(currentE))
		if mM != nil {
			o1, err := strconv.Atoi(currentE[mM[2]:mM[3]])
			if err != nil {
				log.Fatal(err)
			}
			o2, err := strconv.Atoi(currentE[mM[4]:mM[5]])
			if err != nil {
				log.Fatal(err)
			}
			strBuild := currentE[:mM[0]]
			strBuild += fmt.Sprintf("%d", o1*o2)
			strBuild += currentE[mM[1]:]
			currentE = strBuild
			continue
		}
		// If we could not reduce any further
		if start == -1 && mA == nil && mM == nil {
			return currentE
		}
	}
}

func init() {
	regA = regexp.MustCompile(`([0-9]+) \+ ([0-9]+)`)
	regM = regexp.MustCompile(`([0-9]+) \* ([0-9]+)`)
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
