package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	reHgt *regexp.Regexp
	reHcl *regexp.Regexp
	rePID *regexp.Regexp
)

type pass struct {
	fields map[string]string
}

func (p *pass) hasAllFields() bool {
	requiredFields := []string{"byr", "iyr", "eyr", "hgt", "hcl", "ecl", "pid"}

	for _, requiredField := range requiredFields {
		if _, ok := p.fields[requiredField]; !ok {
			return false
		}
	}
	return true
}

func (p *pass) validate() bool {
	if !p.hasAllFields() {
		return false
	}

	for field, val := range p.fields {
		switch field {
		case "byr":
			if !validateYear(val, 1920, 2002) {
				return false
			}
		case "iyr":
			if !validateYear(val, 2010, 2020) {
				return false
			}
		case "eyr":
			if !validateYear(val, 2020, 2030) {
				return false
			}
		case "hgt":
			if !validateHgt(val) {
				return false
			}
		case "hcl":
			if !validateHcl(val) {
				return false
			}
		case "ecl":
			if !validateEcl(val) {
				return false
			}
		case "pid":
			if !validatePID(val) {
				return false
			}
		}
	}
	return true
}

func validateYear(yr string, min, max int) bool {
	yrN, err := strconv.Atoi(yr)
	if err != nil {
		log.Fatal(err)
	}
	return yrN >= min && yrN <= max
}

func validateHgt(hgt string) bool {
	m := reHgt.FindSubmatch([]byte(hgt))
	if m == nil {
		return false
	}
	v, err := strconv.Atoi(string(m[1]))
	if err != nil {
		log.Fatal(err)
	}
	switch string(m[2]) {
	case "cm":
		return v >= 150 && v <= 193
	case "in":
		return v >= 59 && v <= 76
	}
	return true
}

func validateHcl(hcl string) bool {
	m := reHcl.Find([]byte(hcl))
	if m == nil {
		return false
	}
	return true
}

func validateEcl(ecl string) bool {
	validEcls := []string{"amb", "blu", "brn", "gry", "grn", "hzl", "oth"}
	for _, t := range validEcls {
		if ecl == t {
			return true
		}
	}
	return false
}

func validatePID(pid string) bool {
	m := rePID.Find([]byte(pid))
	if m == nil {
		return false
	}
	return true
}

func init() {
	reHgt = regexp.MustCompile(`^([0-9]+)([a-z]+)$`)
	reHcl = regexp.MustCompile(`^#[0-9a-f]+$`)
	rePID = regexp.MustCompile(`^[0-9]{9}$`)
}

func main() {
	f, err := os.Open("./input")
	if err != nil {
		log.Fatal(err)
	}

	passes := make([]pass, 0)

	// Read in the seed
	currentPass := pass{make(map[string]string)}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			passes = append(passes, currentPass)
			currentPass = pass{make(map[string]string)}
			continue
		}
		subs := strings.Split(line, " ")
		for _, sub := range subs {
			sep := strings.Split(sub, ":")
			currentPass.fields[sep[0]] = sep[1]
		}
	}
	if scanner.Err() != nil {
		log.Fatal(scanner.Err())
	}
	hasFields := 0
	validFields := 0
	for _, pass := range passes {
		if pass.hasAllFields() {
			hasFields++
		}
		if pass.validate() {
			validFields++
		}
	}
	fmt.Println("# Passports with required fields", hasFields)
	fmt.Println("# Passports with valid fields", validFields)
}
