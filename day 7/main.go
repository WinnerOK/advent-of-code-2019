package main

import (
	"fmt"
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

func runAmplifier(source []int, order []int) int {
	nextInput := 0
	for _, amplifier := range order {
		output, _ := SimulateMachine(
			CreateState(source, []int{amplifier, nextInput}, false),
		)
		nextInput = output[len(output)-1]
	}
	return nextInput
}

func GreatestAmplifierSignal(source []int) (int, []int) {
	maxIDX := 0
	maxOut := 0
	possibleOrders := permutations([]int{0, 1, 2, 3, 4})
	for i, order := range possibleOrders {
		output := runAmplifier(source, order)
		if output > maxOut {
			maxIDX = i
			maxOut = output
		}
	}

	return maxOut, possibleOrders[maxIDX]
}

func runAmplifierFeedback(source []int, order []int) int {
	lastOutput := 0
	amplifierInput := []int{0}
	amplifiers := []ProgramState{}

	for _, phase := range order {
		amplifiers = append(amplifiers,
			CreateState(
				source,
				[]int{phase},
				true,
			),
		)
	}

	countRunningAmp := func() int {
		result := 0
		for _, amp := range amplifiers {
			if amp.PC.running {
				result += 1
			}
		}
		return result
	}

	condition := true
	for ok := true; ok; ok = condition { // do while emulation
		for i, amp := range amplifiers {
			if amp.PC.running {
				amp.input = append(amp.input, amplifierInput...)
				result, newState := SimulateMachine(amp)
				amplifiers[i] = newState

				if len(result) > 0 {
					lastOutput = result[len(result)-1]
					amplifierInput = result
				}
			}
		}

		condition = countRunningAmp() == len(amplifiers)
	}
	return lastOutput
}

func GreatestAmplifierFeedbackSignal(source []int) (int, []int) {
	maxIDX := 0
	maxOut := 0
	possibleOrders := permutations([]int{5, 6, 7, 8, 9})
	for i, order := range possibleOrders {
		output := runAmplifierFeedback(source, order)

		if output > maxOut {
			maxIDX = i
			maxOut = output
		}
	}

	return maxOut, possibleOrders[maxIDX]
}

func main() {
	fileInput := readInput("./in.txt")
	source := stringSliceToIntSlice(fileInput)

	maxOut1, maxOrder1 := GreatestAmplifierSignal(source)
	fmt.Printf("Part 1| MaxOut: %8d order: %v\n", maxOut1, maxOrder1)

	maxOut2, maxOrder2 := GreatestAmplifierFeedbackSignal(source)
	fmt.Printf("Part 2| MaxOut: %8d order: %v\n", maxOut2, maxOrder2)
}
