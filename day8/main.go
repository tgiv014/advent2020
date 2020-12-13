package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

var (
	reInst *regexp.Regexp
)

func init() {
	var err error
	reInst, err = regexp.Compile("([a-z]{3}) ([\\-\\+0-9]+)")
	if err != nil {
		log.Fatal(err)
	}
}

type instruction struct {
	op      string
	operand int
	runs    int // used for loop detect
}

type program struct {
	pc           int
	acc          int
	instructions []instruction
}

func (p program) clone() program {
	new_program := p // Clone!
	new_program.instructions = make([]instruction, len(p.instructions))
	copy(new_program.instructions, p.instructions)

	// Ensure pre-zeroed
	new_program.pc = 0
	new_program.acc = 0
	for i, _ := range new_program.instructions {
		new_program.instructions[i].runs = 0
	}
	return new_program
}

func (p *program) execute() (int, int) {
	for p.pc >= 0 && p.pc < len(p.instructions) {
		inst := &p.instructions[p.pc]
		if inst.runs > 0 {
			return p.acc, p.pc
		}
		switch inst.op {
		case "acc":
			p.acc += inst.operand
			p.pc += 1
		case "jmp":
			p.pc += inst.operand
		case "nop":
			p.pc += 1
		default:
			fmt.Println("Unsupported instruction")
		}
		inst.runs += 1
	}
	return p.acc, p.pc
}

func parse_line(line string) (instruction, error) {
	groups := reInst.FindSubmatch([]byte(line))
	if len(groups) == 0 {
		return instruction{}, errors.New("Yeah that didn't work")
	}
	op := string(groups[1])
	value, err := strconv.Atoi(string(groups[2]))
	if err != nil {
		log.Fatal(err)
	}
	inst := instruction{op, value, 0}
	return inst, nil
}

func main() {
	f, err := os.Open("./input")
	if err != nil {
		log.Fatalf("Uh oh %e", err)
	}
	defer f.Close()

	program := program{}
	mutable_indices := []int{}
	scanner := bufio.NewScanner(f)
	i := 0
	for scanner.Scan() {
		line := scanner.Text()
		inst, err := parse_line(line)
		if err != nil {
			log.Fatal(err)
		}
		program.instructions = append(program.instructions, inst)
		if inst.op == "nop" || inst.op == "jmp" {
			mutable_indices = append(mutable_indices, i)
		}
		i += 1
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	result, _ := program.execute()
	fmt.Println("Part 1:", result)

	for _, mutate := range mutable_indices {
		mutatedProgram := program.clone()
		instToMutate := &mutatedProgram.instructions[mutate]
		switch instToMutate.op {
		case "jmp":
			instToMutate.op = "nop"
		case "nop":
			instToMutate.op = "jmp"
		default:
			log.Fatal("You done messed up")
		}

		res_acc, res_pc := mutatedProgram.execute()
		if res_pc == len(program.instructions) {
			fmt.Println("Part 2:", res_acc)
			break
		}
	}
}
