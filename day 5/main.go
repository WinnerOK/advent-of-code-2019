package main

import (
	"fmt"
	"strconv"
)

//Opcodes
const (
	ADDITION       = 1
	MULTIPLICATION = 2
	INPUT          = 3
	OUTPUT         = 4
	JUMP_IF_TRUE   = 5
	JUMP_IF_FALSE  = 6
	LESS_THAN      = 7
	EQUALS         = 8
	HALT           = 99
)

//Parameter modes
const (
	POSITION  = 0
	IMMEDIATE = 1
)

const INPUT_VALUE = 5

func interpretOpcode(code int) (int, []int) {
	opcodeStr := strconv.FormatInt(int64(code), 10)
	var opcode int
	var readingModes []int
	if len(opcodeStr) <= 2 {
		opcode = code
	} else {
		opcode, _ = strconv.Atoi(opcodeStr[len(opcodeStr)-2:])

		for i := len(opcodeStr) - 3; i >= 0; i-- {
			readingMode := int(opcodeStr[i]) - int('0')
			readingModes = append(readingModes, readingMode)
		}
	}

	return opcode, readingModes
}

func getOrDefault(container []int, idx int, defaultValue int) int {
	if idx >= 0 && idx < len(container) {
		return container[idx]
	} else {
		return defaultValue
	}
}

func performOP(memory []int, cursor int, input int) int {

	opcode, readingModes := interpretOpcode(memory[cursor])
	argumentsRead := 0
	getNextArgValue := func() (int, int) {
		arg := memory[cursor+1+argumentsRead]
		var returnValue int
		switch getOrDefault(readingModes, argumentsRead, POSITION) {
		case POSITION:
			returnValue = memory[arg]
		case IMMEDIATE:
			returnValue = arg
		default:
			panic("Unexpected mode!\n")
		}
		argumentsRead += 1
		return returnValue, arg
	}
	switch opcode {
	case ADDITION:
		arg1, _ := getNextArgValue()
		arg2, _ := getNextArgValue()
		_, dest := getNextArgValue()
		memory[dest] = arg1 + arg2
	case MULTIPLICATION:
		arg1, _ := getNextArgValue()
		arg2, _ := getNextArgValue()
		_, dest := getNextArgValue()
		memory[dest] = arg1 * arg2
	case INPUT:
		_, dest := getNextArgValue()
		memory[dest] = input
	case OUTPUT:
		arg1, _ := getNextArgValue()
		fmt.Printf("[Machine] %d\n", arg1)
	case JUMP_IF_TRUE:
		arg1, _ := getNextArgValue()
		dest, _ := getNextArgValue()
		if arg1 != 0 {
			return dest
		}
	case JUMP_IF_FALSE:
		arg1, _ := getNextArgValue()
		dest, _ := getNextArgValue()
		if arg1 == 0 {
			return dest
		}
	case LESS_THAN:
		arg1, _ := getNextArgValue()
		arg2, _ := getNextArgValue()
		_, dest := getNextArgValue()
		if arg1 < arg2 {
			memory[dest] = 1
		} else {
			memory[dest] = 0
		}
	case EQUALS:
		arg1, _ := getNextArgValue()
		arg2, _ := getNextArgValue()
		_, dest := getNextArgValue()
		if arg1 == arg2 {
			memory[dest] = 1
		} else {
			memory[dest] = 0
		}
	default:
		panic("Unexpected opcode!\n")
	}
	return cursor + 1 + argumentsRead
}

func SimulateMachine(source []int, input int) int {
	memory := make([]int, len(source))
	copy(memory, source)
	cursor := 0
	for memory[cursor] != HALT {
		cursor = performOP(memory, cursor, input)
	}

	return memory[0]
}

func main() {
	input := readInput("./in.txt")
	source := stringSliceToIntSlice(input)
	//fmt.Printf("%v\n", source)
	SimulateMachine(source, INPUT_VALUE)

}
