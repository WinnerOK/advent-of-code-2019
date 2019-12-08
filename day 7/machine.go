package main

import (
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

type ProgramState struct {
	memory         []int
	PC             ProgramCounter
	input          []int
	inputRead      int
	stopOnFirstOut bool
}

func CreateState(source []int, input []int, stopOnFirstOut bool) ProgramState {
	memory := make([]int, len(source))
	copy(memory, source)
	inp := make([]int, len(input))
	copy(inp, input)
	return ProgramState{
		memory: memory,
		PC: ProgramCounter{
			running: true,
			value:   0,
		},
		input:          input,
		inputRead:      0,
		stopOnFirstOut: stopOnFirstOut,
	}
}

func (s *ProgramState)addInput(input []int)  {
	s.input = append(s.input, input...)
}

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

func SimulateMachine(state ProgramState) ([]int, ProgramState) {
	var memory []int
	var PC ProgramCounter
	var input []int
	var inputRead int
	var stopOnFirstOut bool

	restoreState := func() {
		memory = state.memory
		PC = state.PC
		input = state.input
		inputRead = state.inputRead
		stopOnFirstOut = state.stopOnFirstOut
	}
	restoreState()

	saveState := func() ProgramState {
		newState := ProgramState{
			memory:         memory,
			PC:             PC,
			input:          input,
			inputRead:      inputRead,
			stopOnFirstOut: stopOnFirstOut,
		}
		return newState
	}

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
	return output, saveState()
}
