package main

import (
	"fmt"
	"os"
)

/*
OpCodes:
	- 1 add
	- 2 mult
	- 99 halt
 */

var source []int
var desiredOutput = 19690720

func performOP(memory []int, opcode int, src1 int, src2 int, dest int) {
	switch opcode {
	case 1:
		memory[dest] = memory[src1] + memory[src2]
	case 2:
		memory[dest] = memory[src1] * memory[src2]
	}
}

func simulateMachine(noun int, verb int) int {
	memory:= make([]int, len(source))
	copy(memory, source)
	cursor := 0
	memory[1] = noun
	memory[2] = verb
	for memory[cursor] != 99 {
		performOP(memory, memory[cursor], memory[cursor+1], memory[cursor+2], memory[cursor+3])
		cursor += 4
	}
	return memory[0]
}

func main() {
	input := readInput("./in.txt")
	source = stringSliceToIntSlice(input)

	fmt.Printf("Answer for part 1: %d\n", simulateMachine(12, 2))

	for noun := 0; noun <= 99; noun++ {
		for verb := 0; verb <= 99; verb++ {
			if simulateMachine(noun, verb) == desiredOutput {
				fmt.Printf("Answer for part 2: %d", 100*noun+verb)
				os.Exit(0)
			}
		}
	}
}
