package main

import (
	"fmt"
	"strconv"
)

const (
	DEBUG = false
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

var (
	DEFAULT_MEM = makeBigInt(0)
)

type memoryType = map[string]bigInt

func makeMemory(from []int) memoryType {
	result := memoryType{}
	for i, v := range from {
		result[makeBigInt(i).String()] = makeBigInt(v)
	}
	return result
}

type ProgramState struct {
	memory         memoryType
	PC             ProgramCounter
	input          []bigInt
	inputRead      int
	stopOnFirstOut bool
	relativeBase   bigInt
}

func CreateState(source []int, input []int, stopOnFirstOut bool) ProgramState {
	memory := makeMemory(source)
	bigInput := []bigInt{}
	for _, v := range input {
		bigInput = append(bigInput, makeBigInt(v))
	}
	return ProgramState{
		memory: memory,
		PC: ProgramCounter{
			running: true,
			value:   makeBigInt(0),
		},
		input:          bigInput,
		inputRead:      0,
		stopOnFirstOut: stopOnFirstOut,
		relativeBase:   makeBigInt(0),
	}
}

func (s *ProgramState) addInput(input []int) {
	for _, v := range input {
		s.input = append(s.input, makeBigInt(v))
	}
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

func SimulateMachine(state ProgramState) ([]bigInt, ProgramState) {
	var memory memoryType
	var PC ProgramCounter
	var input []bigInt
	var inputRead int
	var stopOnFirstOut bool
	var relativeBase bigInt

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

	output := []bigInt{}

	outputCheck := func() bool {
		if stopOnFirstOut {
			return len(output) == 0
		} else {
			return true
		}
	}

	nextInput := func() bigInt {
		if inputRead < len(input) {
			inputRead += 1
			return input[inputRead-1]
		} else {
			panic("Input ended!")
		}

	}
	memSet := func(addr, value bigInt) {
		memory[addr.String()] = value
	}
	memGet := func(addr bigInt) bigInt {
		if value, ok := memory[addr.String()]; ok {
			return value
		} else {
			return DEFAULT_MEM
		}
	}

	printMem := func() {
		print("map[")
		for k, v := range memory{
			fmt.Printf(" %s:%s", k, v.String())
		}
		print("]")
	}

	log := func(s string) {
		if DEBUG {
			fmt.Printf("step %d| ", step)
			print(s)
			printMem()
			println()
			step += 1
		}
	}

	for leBig(PC.value, makeBigInt(len(memory))) && PC.running && outputCheck() && step <= 80 {
		opCode, readingModes := interpretCode(int(memGet(PC.value).Int64()))

		deRef := func(readingMode int, address bigInt) bigInt {
			var result bigInt
			switch readingMode {
			case POSITION:
				result = memGet(address)
			case IMMEDIATE:
				result = address
			case RELATIVE:
				result = memGet(addBig(relativeBase, address))
			default:
				panic("Unknown reading mode")
			}
			return result
		}

		param := func(index int) bigInt {
			return deRef(getOrDefault(readingModes, index-1, POSITION),
				memGet(addBig(PC.value, makeBigInt(index))))
		}

		address := func(index int) bigInt {
			readingMode := getOrDefault(readingModes, index-1, POSITION)
			addr := addBig(PC.value, makeBigInt(index))
			var result bigInt
			switch readingMode {
			case RELATIVE:
				result = addBig(relativeBase, addr)
			case POSITION:
				result = addr
			default:
				panic("Unknown reading mode for addr")
			}
			return result
		}

		addOp := func() {
			arg1 := param(1)
			arg2 := param(2)
			result := addBig(arg1, arg2)
			memSet(address(3), result)
			log(fmt.Sprintf("[%s] opcode: %d (add), arg1: %s, arg2: %s, result: %s, memory: ",
				PC.value.String(), opCode, arg1.String(), arg2.String(), result.String(),
				))
			PC.add(makeBigInt(4))
		}

		multOp := func() {
			arg1 := param(1)
			arg2 := param(2)
			result := makeBigInt(0).Mul(arg1, arg2)
			memSet(address(3), result)
			log(fmt.Sprintf("[%s] opcode: %d (mul), arg1: %s, arg2: %s, result: %s, memory: ",
				PC.value.String(), opCode, arg1.String(), arg2.String(), result.String(),
			))
			PC.add(makeBigInt(4))
		}

		haltOp := func() {
			log(fmt.Sprintf("[%s] opcode: %d (halt),  memory: ",
				PC.value.String(), opCode,
			))
			PC.halt()
		}

		inputOp := func() {

			memSet(address(1), nextInput())
			log(fmt.Sprintf("[%s] opcode: %d (in),  memory: ",
				PC.value.String(), opCode,
			))
			PC.add(makeBigInt(2))
		}

		outputOp := func() {
			value := param(1)
			output = append(output, value)
			log(fmt.Sprintf("[%s] opcode: %d (out),  memory: ",
				PC.value.String(), opCode,
			))
			PC.add(makeBigInt(2))
		}

		jumpIfTrueOp := func() {
			arg1 := param(1)
			dest := param(2)
			log(fmt.Sprintf("[%s] opcode: %d (jump true), arg1: %s, dest: %s memory: ",
				PC.value.String(), opCode, arg1.String(), dest.String(),
			))
			if !equalsBig(arg1, makeBigInt(0)) {
				PC.jump(dest)
			} else {
				PC.add(makeBigInt(3))
			}
		}

		jumpIfFalseOp := func() {
			arg1 := param(1)
			dest := param(2)
			log(fmt.Sprintf("[%s] opcode: %d (jump false), arg1: %s, dest: %s memory: ",
				PC.value.String(), opCode, arg1.String(), dest.String(),
			))
			if equalsBig(arg1, makeBigInt(0)) {
				PC.jump(dest)
			} else {
				PC.add(makeBigInt(3))
			}
		}

		lessThanOp := func() {
			arg1 := param(1)
			arg2 := param(2)
			var result bigInt
			if leBig(arg1, arg2) {
				result = makeBigInt(1)
			} else {
				result = makeBigInt(0)
			}
			memSet(address(3), result)
			log(fmt.Sprintf("[%s] opcode: %d (LE), arg1: %s, arg2: %s memory: ",
				PC.value.String(), opCode, arg1.String(), arg2.String(),
			))
			PC.add(makeBigInt(4))
		}

		equalsOp := func() {
			arg1 := param(1)
			arg2 := param(2)
			var result bigInt
			if equalsBig(arg1, arg2) {
				result = makeBigInt(1)
			} else {
				result = makeBigInt(0)
			}
			memSet(address(3), result)
			log(fmt.Sprintf("[%s] opcode: %d (equals), arg1: %s, arg2: %s memory: ",
				PC.value.String(), opCode, arg1.String(), arg2.String(),
			))
			PC.add(makeBigInt(4))
		}

		relativeBaseAdjOp := func() {
			arg1 := param(1)
			relativeBase = addBig(relativeBase, arg1)
			log(fmt.Sprintf("[%s] opcode: %d (relativeBase adj), arg1: %s, memory: ",
				PC.value.String(), opCode, arg1.String(),
			))
			PC.add(makeBigInt(2))
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
