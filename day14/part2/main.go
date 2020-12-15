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
	mask_x    []int
	mem       map[uint64]uint64
}

func (in *input) build_mask(m []byte) {
	in.mask_high = uint64(0)
	in.mask_x = make([]int, 0) // Reset X mask
	// Loops from start of string to end
	// We will switch bitshift around to make sense here
	for i, b := range m {
		n := len(m) - i - 1 // Actual bit index
		switch b {
		case '0':
			continue
		case '1':
			in.mask_high |= (1 << n)
		case 'X':
			in.mask_x = append(in.mask_x, n)
		default:
			fmt.Println("Invalid input")
		}
	}
}

func (in *input) apply(addr uint64, data uint64) {
	var r uint64 = addr
	r |= in.mask_high // Still apply highs

	// Consider I a "len(in.mask_x)"-bit wide counter
	for i := 0; i < (1 << len(in.mask_x)); i++ {
		for j, v := range in.mask_x {
			if i&(1<<j) != 0 {
				r |= (1 << v)
			} else {
				r &= ^(1 << v)
			}
		}
		in.mem[r] = data
	}
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

	in := input{0, make([]int, 0), make(map[uint64]uint64)}

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
