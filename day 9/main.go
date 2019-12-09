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

	fmt.Printf("%v\n", output)
}
