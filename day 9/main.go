package main

import (
	"fmt"
	"runtime"
	"runtime/debug"
)

func main() {

	debug.SetGCPercent(-1)

	fileInput := readInput("./in.txt")
	source := stringSliceToIntSlice(fileInput)

	state := CreateState(source, []int{1}, false)

	output, _ := SimulateMachine(
		state,
	)

	fmt.Printf("%v\n", output)

	runtime.GC()
}
