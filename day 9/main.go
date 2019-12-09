package main

import (
	"fmt"
)

func main() {
	fileInput := readInput("./in.txt")
	source := stringSliceToIntSlice(fileInput)

	output, _ := SimulateMachine(
		CreateState(source, []int{1}, false),
	)

	fmt.Printf("Answer part 1: %v\n", output[0])

	output, _ = SimulateMachine(
		CreateState(source, []int{2}, false),
	)

	fmt.Printf("Answer part 2: %v\n", output[0])
}
