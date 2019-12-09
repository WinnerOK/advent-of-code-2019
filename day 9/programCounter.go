package main

type ProgramCounter struct {
	running bool
	value   int64
}

func (pc *ProgramCounter) add(offset int64) {
	pc.value += offset
}

func (pc *ProgramCounter) jump(addr int64) {
	pc.value = addr
}

func (pc *ProgramCounter) halt() {
	pc.running = false
}
