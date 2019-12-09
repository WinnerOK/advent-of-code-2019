package main

import (
	"fmt"
	"strconv"
)

const (
	DEBUG = true
)

var (
	step = 0
)

//Opcodes
const (
	ADDITION          = 1
	MULTIPLICATION    = 2
	INPUT             = 3
	OUTPUT            = 4
	JUMP_IF_TRUE      = 5
	JUMP_IF_FALSE     = 6
	LESS_THAN         = 7
	EQUALS            = 8
	RELATIVE_BASE_ADJ = 9
	HALT              = 99
)

//Parameter modes
const (
	POSITION  = 0
	IMMEDIATE = 1
	RELATIVE  = 2
)

const (
	DEFAULT_MEM = 0
)

type memoryType = map[int]int

func makeMemory(from []int) memoryType {
	result := memoryType{}
	for i, v := range from {
		result[i] = v
	}
	return result
}

type ProgramState struct {
	memory         memoryType
	PC             ProgramCounter
	input          []int
	inputRead      int
	stopOnFirstOut bool
	relativeBase   int
}

func CreateState(source []int, input []int, stopOnFirstOut bool) ProgramState {
	memory := makeMemory(source)
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
		relativeBase:   0,
	}
}

func (s *ProgramState) addInput(input []int) {
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
	var memory memoryType
	var PC ProgramCounter
	var input []int
	var inputRead int
	var stopOnFirstOut bool
	var relativeBase int

	restoreState := func() {
		memory = state.memory
		PC = state.PC
		input = state.input
		inputRead = state.inputRead
		stopOnFirstOut = state.stopOnFirstOut
		relativeBase = state.relativeBase
	}
	restoreState()

	saveState := func() ProgramState {
		return ProgramState{
			memory:         memory,
			PC:             PC,
			input:          input,
			inputRead:      inputRead,
			stopOnFirstOut: stopOnFirstOut,
			relativeBase:   relativeBase,
		}
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
	memSet := func(addr, value int) { memory[addr] = value }
	memGet := func(addr int) int {
		if value, ok := memory[addr]; ok {
			return value
		} else {
			return DEFAULT_MEM
		}
	}

	for PC.value < len(memory) && PC.running && outputCheck() {
		opCode, readingModes := interpretCode(memGet(PC.value))

		deRef := func(readingMode, address int) int {
			var result int
			switch readingMode {
			case POSITION:
				result = memGet(address)
			case IMMEDIATE:
				result = address
			case RELATIVE:
				result = memGet(relativeBase + address)
			default:
				panic("Unknown reading mode")
			}
			return result
		}

		param := func(index int) int {
			return deRef(getOrDefault(readingModes, index-1, POSITION), memGet(PC.value+index))
		}

		log := func(s string) {
			if DEBUG {
				fmt.Printf("step %d| ", step)
				print(s)
				step += 1
			}
		}

		addOp := func() {
			arg1 := param(1)
			arg2 := param(2)
			result := arg1 + arg2
			memSet(memGet(PC.value+3), result)
			log(fmt.Sprintf("[%d] opcode: %d (add), arg1: %d, arg2: %d, result: %d, memory: %v\n",
				PC.value, opCode, arg1, arg2, result, memory))
			PC.add(4)
		}

		multOp := func() {
			arg1 := param(1)
			arg2 := param(2)
			result := arg1 * arg2
			memSet(memGet(PC.value+3), result)
			log(fmt.Sprintf("[%d] opcode: %d (mul), arg1: %d, arg2: %d, result: %d, memory: %v\n",
				PC.value, opCode, arg1, arg2, result, memory))
			PC.add(4)
		}

		haltOp := func() {
			log(fmt.Sprintf("[%d] opcode: %d (halt), memory: %v\n",
				PC.value, opCode, memory))
			PC.halt()
		}

		inputOp := func() {
			memSet(memGet(PC.value+1), nextInput())
			log(fmt.Sprintf("[%d] opcode: %d (in), memory: %v\n",
				PC.value, opCode, memory))
			PC.add(2)
		}

		outputOp := func() {
			value := param(1)
			output = append(output, value)
			log(fmt.Sprintf("[%d] opcode: %d (out), memory: %v\n",
				PC.value, opCode, memory))
			PC.add(2)
		}

		jumpIfTrueOp := func() {
			arg1 := param(1)
			dest := param(2)
			log(fmt.Sprintf("[%d] opcode: %d (jump true), arg1: %d, dest: %d memory: %v\n",
				PC.value, opCode, arg1, dest, memory))
			if arg1 != 0 {
				PC.jump(dest)
			} else {
				PC.add(3)
			}
		}

		jumpIfFalseOp := func() {
			arg1 := param(1)
			dest := param(2)
			log(fmt.Sprintf("[%d] opcode: %d (jump false), arg1: %d, dest: %d memory: %v\n",
				PC.value, opCode, arg1, dest, memory))
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
			memSet(memGet(PC.value+3), result)
			log(fmt.Sprintf("[%d] opcode: %d (LE), arg1: %d, arg2: %d memory: %v\n",
				PC.value, opCode, arg1, arg2, memory))
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
			memSet(memGet(PC.value+3), result)
			log(fmt.Sprintf("[%d] opcode: %d (equals), arg1: %d, arg2: %d memory: %v\n",
				PC.value, opCode, arg1, arg2, memory))
			PC.add(4)
		}

		relativeBaseAdjOp := func() {
			arg1 := param(1)
			relativeBase += arg1
			log(fmt.Sprintf("[%d] opcode: %d (relativeBase adj), arg1: %d, memory: %v\n",
				PC.value, opCode, arg1, memory))
			PC.add(2)
		}

		switch opCode {
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
		case RELATIVE_BASE_ADJ:
			relativeBaseAdjOp()
		default:
			panic("Unexpected opCode")
		}

	}
	return output, saveState()
}
