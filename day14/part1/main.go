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
	reMask *regexp.Regexp
	reMem  *regexp.Regexp
)

type input struct {
	mask_high uint64
	mask_low  uint64
	mem       map[uint64]uint64
}

func (in *input) build_mask(m []byte) {
	in.mask_low = uint64(0)
	in.mask_high = uint64(0)
	// Loops from start of string to end
	// We will switch bitshift around to make sense here
	for i, b := range m {
		n := len(m) - i - 1
		switch b {
		case '0':
			in.mask_low |= (1 << n)
		case '1':
			in.mask_high |= (1 << n)
		case 'X':
			continue
		default:
			fmt.Println("Invalid input")
		}
	}
}

func (in *input) apply(addr uint64, data uint64) {
	var r uint64 = data
	r |= in.mask_high
	r &= ^in.mask_low
	in.mem[addr] = r
}

func (in *input) sum() uint64 {
	var t uint64 = 0
	for _, val := range in.mem {
		t += val
	}
	return t
}

func init() {
	var err error
	reMask, err = regexp.Compile("^mask = ([01X]+)$")
	if err != nil {
		log.Fatal(err)
	}
	reMem, err = regexp.Compile("^mem\\[([0-9]+)\\] = ([0-9]+)$")
	if err != nil {
		log.Fatal(err)
	}
}
func main() {
	f, err := os.Open("./input")
	if err != nil {
		log.Fatalf("Uh oh %e", err)
	}
	defer f.Close()

	in := input{0, 0, make(map[uint64]uint64)}

	scanner := bufio.NewScanner(f)
	for i := 0; scanner.Scan(); i++ {
		line := scanner.Text()

		maskMatch := reMask.FindSubmatch([]byte(line))
		if maskMatch != nil {
			in.build_mask(maskMatch[1])
			continue
		}

		memMatch := reMem.FindSubmatch([]byte(line))
		if memMatch != nil {
			addrI, err := strconv.ParseUint(string(memMatch[1]), 10, 64)
			if err != nil {
				fmt.Println(err)
			}
			dataI, err := strconv.ParseUint(string(memMatch[2]), 10, 64)
			if err != nil {
				fmt.Println(err)
			}
			in.apply(addrI, dataI)
			continue
		}
		fmt.Println("Parser missed a line :(")
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Sum:", in.sum())
}
