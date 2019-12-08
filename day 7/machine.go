package main

import (
	"strconv"
)

func interpretCode(code int) (int, []int) {
	CodeStr := strconv.FormatInt(int64(code), 10)
	var opcode int
	var readingModes []int
	if len(CodeStr) <= 2 {
		opcode = code
	} else {
		opcode, _ = strconv.Atoi(CodeStr[len(CodeStr)-2:])

		for i := len(CodeStr) - 3; i >= 0; i-- {
			readingMode := int(CodeStr[i]) - int('0')
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

func SimulateMachine(source []int, input []int, stopOnFirstOut bool) []int {
	memory := make([]int, len(source))
	copy(memory, source)
	PC := ProgramCounter{
		running: true,
		value:   0,
	}
	inputRead := 0
	output := []int{}

	outputCheck := func() bool {
		if stopOnFirstOut {
			return len(output) == 0
		} else {
			return true
		}
	}

	nextInput := func() int {
		if inputRead < len(input) {
			inputRead += 1
			return input[inputRead-1]
		} else {
			panic("Input ended!")
		}

	}
	memset := func(addr, value int) { memory[addr] = value }

	for PC.value < len(memory) && PC.running && outputCheck() {
		opcode, readingModes := interpretCode(memory[PC.value])

		deref := func(readingMode, address int) int {
			var result int
			switch readingMode {
			case POSITION:
				result = memory[address]
			case IMMEDIATE:
				result = address
			}
			return result
		}

		param := func(index int) int {
			return deref(getOrDefault(readingModes, index-1, POSITION), memory[PC.value+index])
		}

		addOp := func() {
			arg1 := param(1)
			arg2 := param(2)
			result := arg1 + arg2
			memset(memory[PC.value+3], result)
			PC.add(4)
		}

		multOp := func() {
			arg1 := param(1)
			arg2 := param(2)
			result := arg1 * arg2
			memset(memory[PC.value+3], result)
			PC.add(4)
		}

		haltOp := func() {
			PC.halt()
		}

		inputOp := func() {
			memset(memory[PC.value+1], nextInput())
			PC.add(2)
		}

		outputOp := func() {
			value := param(1)
			output = append(output, value)
			PC.add(2)
		}

		jumpIfTrueOp := func() {
			arg1 := param(1)
			dest := param(2)
			if arg1 != 0 {
				PC.jump(dest)
			} else {
				PC.add(3)
			}
		}

		jumpIfFalseOp := func() {
			arg1 := param(1)
			dest := param(2)
			if arg1 == 0 {
				PC.jump(dest)
			} else {
				PC.add(3)
			}
		}

		lessThanOp := func() {
			arg1 := param(1)
			arg2 := param(2)
			var result int
			if arg1 < arg2 {
				result = 1
			} else {
				result = 0
			}
			memset(memory[PC.value+3], result)
			PC.add(4)
		}

		equalsOp := func() {
			arg1 := param(1)
			arg2 := param(2)
			var result int
			if arg1 == arg2 {
				result = 1
			} else {
				result = 0
			}
			memset(memory[PC.value+3], result)
			PC.add(4)
		}

		switch opcode {
		case ADDITION:
			addOp()
		case MULTIPLICATION:
			multOp()
		case HALT:
			haltOp()
		case INPUT:
			inputOp()
		case OUTPUT:
			outputOp()
		case JUMP_IF_TRUE:
			jumpIfTrueOp()
		case JUMP_IF_FALSE:
			jumpIfFalseOp()
		case LESS_THAN:
			lessThanOp()
		case EQUALS:
			equalsOp()
		}

	}
	//fmt.Printf("Machine state before output: %v\n", memory)
	return output
}
