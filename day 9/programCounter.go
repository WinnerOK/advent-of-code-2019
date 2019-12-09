package main

type ProgramCounter struct {
	running bool
	value   int
}

func (pc *ProgramCounter) add(offset int) {
	pc.value += offset
}

func (pc *ProgramCounter) jump(addr int) {
	pc.value = addr
}

func (pc *ProgramCounter) halt() {
	pc.running = false
}
