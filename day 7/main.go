package main

import (
	"fmt"
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

// ref https://play.golang.org/p/Ulyo1H2Bii
func permutations(arr []int) [][]int {
	var helper func([]int, int)
	res := [][]int{}

	helper = func(arr []int, n int) {
		if n == 1 {
			tmp := make([]int, len(arr))
			copy(tmp, arr)
			res = append(res, tmp)
		} else {
			for i := 0; i < n; i++ {
				helper(arr, n-1)
				if n%2 == 1 {
					tmp := arr[i]
					arr[i] = arr[n-1]
					arr[n-1] = tmp
				} else {
					tmp := arr[0]
					arr[0] = arr[n-1]
					arr[n-1] = tmp
				}
			}
		}
	}
	helper(arr, len(arr))
	return res
}

func GreatestAmplifierSignal(source []int) (int, []int) {
	maxIDX := 0
	maxOut := 0
	possibleOrders := permutations([]int{0, 1, 2, 3, 4})
	for i, order := range possibleOrders {
		previousOut := 0
		for _, amplifier := range order {
			output := SimulateMachine(source, []int{amplifier, previousOut}, false)
			previousOut = output[len(output)-1]
		}

		if previousOut > maxOut {
			maxIDX = i
			maxOut = previousOut
		}
	}

	return maxOut, possibleOrders[maxIDX]
}

func main() {
	fileInput := readInput("./in.txt")
	source := stringSliceToIntSlice(fileInput)

	maxOut, maxOrder := GreatestAmplifierSignal(source)
	fmt.Printf("Part 1| MaxOut: %d, order: %v\n", maxOut, maxOrder)
}
