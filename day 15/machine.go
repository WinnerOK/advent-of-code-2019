package main

import (
	"fmt"
	"strconv"
)

const (
	Debug = false
)

var (
	step = 0
)

//Opcodes
const (
	Addition        = 1
	Multiplication  = 2
	Input           = 3
	Output          = 4
	JumpIfTrue      = 5
	JumpIfFalse     = 6
	LessThan        = 7
	Equals          = 8
	RelativeBaseAdj = 9
	Halt            = 99
)

//Parameter modes
const (
	Position  = 0
	Immediate = 1
	Relative  = 2
)

var (
	DefaultMem = makeBigInt(0)
)

type memoryType map[string]bigInt

func makeMemory(from []int) memoryType {
	result := memoryType{}
	for i, v := range from {
		result[makeBigInt(i).String()] = makeBigInt(v)
	}
	return result
}

type ProgramState struct {
	memory       memoryType
	PC           ProgramCounter
	input        []bigInt
	inputRead    int
	waitForNOut  int
	relativeBase bigInt
}

func CreateState(source []int, input []int, waitForNOut int) ProgramState {
	// pass waitForNOut = 0 to wait for halt
	memory := makeMemory(source)
	var bigInput []bigInt
	for _, v := range input {
		bigInput = append(bigInput, makeBigInt(v))
	}
	return ProgramState{
		memory: memory,
		PC: ProgramCounter{
			running: true,
			value:   makeBigInt(0),
		},
		input:        bigInput,
		inputRead:    0,
		waitForNOut:  waitForNOut,
		relativeBase: makeBigInt(0),
	}
}

func CopyState(state ProgramState) ProgramState {
	newMem := memoryType{}
	for k, v := range state.memory {
		newMem[k] = copyBigInt(v)
	}
	var newBigInput []bigInt
	for _, v := range state.input {
		newBigInput = append(newBigInput, copyBigInt(v))
	}

	return ProgramState{
		memory: newMem,
		PC: ProgramCounter{
			running: state.PC.running,
			value:   copyBigInt(state.PC.value),
		},
		input:        newBigInput,
		inputRead:    state.inputRead,
		waitForNOut:  state.waitForNOut,
		relativeBase: copyBigInt(state.relativeBase),
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
	var waitForNOut int
	var relativeBase bigInt

	restoreState := func() {
		memory = state.memory
		PC = state.PC
		input = state.input
		inputRead = state.inputRead
		waitForNOut = state.waitForNOut
		relativeBase = state.relativeBase
	}
	restoreState()

	saveState := func() ProgramState {
		return ProgramState{
			memory:       memory,
			PC:           PC,
			input:        input,
			inputRead:    inputRead,
			waitForNOut:  waitForNOut,
			relativeBase: relativeBase,
		}
	}

	var output []bigInt

	outputCheck := func() bool {
		if waitForNOut != 0 {
			return len(output) < waitForNOut
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
			return DefaultMem
		}
	}

	printMem := func() {
		print("map[")
		for k, v := range memory {
			fmt.Printf(" %s:%s", k, v.String())
		}
		print("]")
	}

	log := func(s string) {
		if Debug {
			fmt.Printf("step %d| ", step)
			print(s)
			printMem()
			println()
			step += 1
		}
	}

	for leBig(PC.value, makeBigInt(len(memory))) && PC.running && outputCheck() {
		opCode, readingModes := interpretCode(int(memGet(PC.value).Int64()))

		deRef := func(readingMode int, address bigInt) bigInt {
			var result bigInt
			switch readingMode {
			case Position:
				result = memGet(address)
			case Immediate:
				result = address
			case Relative:
				result = memGet(addBig(relativeBase, address))
			default:
				panic("Unknown reading mode")
			}
			return result
		}

		param := func(index int) bigInt {
			return deRef(getOrDefault(readingModes, index-1, Position),
				memGet(addBig(PC.value, makeBigInt(index))))
		}

		address := func(index int) bigInt {
			readingMode := getOrDefault(readingModes, index-1, Position)
			addr := memGet(addBig(PC.value, makeBigInt(index)))
			var result bigInt
			switch readingMode {
			case Relative:
				result = addBig(relativeBase, addr)
			case Position:
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
		case Addition:
			addOp()
		case Multiplication:
			multOp()
		case Halt:
			haltOp()
		case Input:
			inputOp()
		case Output:
			outputOp()
		case JumpIfTrue:
			jumpIfTrueOp()
		case JumpIfFalse:
			jumpIfFalseOp()
		case LessThan:
			lessThanOp()
		case Equals:
			equalsOp()
		case RelativeBaseAdj:
			relativeBaseAdjOp()
		default:
			panic("Unexpected opCode")
		}

	}
	return output, saveState()
}
