package main

type ProgramCounter struct {
	running bool
	value   bigInt
}

func (pc *ProgramCounter) add(offset bigInt) {
	pc.value = makeBigInt(0).Add(pc.value, offset)
}

func (pc *ProgramCounter) jump(addr bigInt) {
	pc.value = addr
}

func (pc *ProgramCounter) halt() {
	pc.running = false
}
